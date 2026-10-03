package micronaut

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	BaseURL     = "https://launch.micronaut.io"
	maxZipBytes = 100 << 20 // 100 MiB safety limit
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

var (
	groupRe    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$`)
	artifactRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

func FetchFeatures(appType string) ([]Feature, error) {
	path := "/application-types/" + url.PathEscape(appType) + "/features"
	resp, err := httpClient.Get(BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("could not fetch Micronaut features (check your internet connection): %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: unexpected status %s", path, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var features []Feature
	var wrapped featuresResponse
	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Features) > 0 {
		features = wrapped.Features
	} else {
		var bare []Feature
		if err := json.Unmarshal(body, &bare); err == nil {
			features = bare
		}
	}
	if len(features) == 0 {
		snippet := string(body)
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return nil, fmt.Errorf("no features found in response from %s: %s", path, snippet)
	}

	sort.SliceStable(features, func(i, j int) bool {
		a, b := features[i], features[j]
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		return strings.ToLower(a.Title) < strings.ToLower(b.Title)
	})
	return features, nil
}

func (p Project) validate() error {
	switch {
	case !groupRe.MatchString(p.Group):
		return fmt.Errorf("invalid group %q (expected something like com.example)", p.Group)
	case !artifactRe.MatchString(p.Artifact):
		return fmt.Errorf("invalid artifact id %q (letters, digits, '.', '_' and '-' only)", p.Artifact)
	}
	return nil
}

// DownloadURL returns the launch.micronaut.io URL that produces the project zip.
// The name is sent as "<group>.<artifact>" so the generated package is the group.
func (p Project) DownloadURL() string {
	q := url.Values{}
	q.Set("lang", p.Lang)
	q.Set("build", p.Build)
	q.Set("test", p.Test)
	q.Set("javaVersion", p.JavaVersion)
	for _, f := range p.Features {
		q.Add("features", f)
	}
	name := p.Group + "." + p.Artifact
	return BaseURL + "/create/" + url.PathEscape(p.Type) + "/" + url.PathEscape(name) + "?" + q.Encode()
}

func (p Project) download() ([]byte, error) {
	resp, err := httpClient.Get(p.DownloadURL())
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("launch.micronaut.io responded with %s: %s",
			resp.Status, apiError(msg))
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxZipBytes+1))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if len(data) > maxZipBytes {
		return nil, fmt.Errorf("response too large")
	}
	return data, nil

}

func (p Project) Generate() error {
	if err := p.validate(); err != nil {
		return err
	}
	dest, err := os.Getwd()
	if err != nil {
		return err
	}
	target := filepath.Join(dest, p.Artifact)
	if entries, err := os.ReadDir(target); err == nil && len(entries) > 0 {
		return fmt.Errorf("directory %q already exists and is not empty", p.Artifact)
	}

	data, err := p.download()
	if err != nil {
		return err
	}
	return unzip(data, dest)
}

func unzip(data []byte, dest string) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("invalid zip from server: %w", err)
	}
	dest = filepath.Clean(dest)

	for _, f := range zr.File {
		path := filepath.Join(dest, filepath.FromSlash(f.Name))
		if path != dest && !strings.HasPrefix(path, dest+string(os.PathSeparator)) {
			return fmt.Errorf("illegal path in archive: %q", f.Name)
		}
		if f.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := extractFile(f, path); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, path string) error {
	mode := f.Mode().Perm()
	if mode == 0 {
		mode = 0o644
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// apiError extracts a readable message from a Micronaut Launch error body.
func apiError(body []byte) string {
	var e struct {
		Message  string `json:"message"`
		Embedded struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"_embedded"`
	}
	if err := json.Unmarshal(body, &e); err == nil {
		var msgs []string
		for _, x := range e.Embedded.Errors {
			if x.Message != "" {
				msgs = append(msgs, x.Message)
			}
		}
		if len(msgs) > 0 {
			return strings.Join(msgs, "; ")
		}
		if e.Message != "" {
			return e.Message
		}
	}
	s := strings.TrimSpace(string(body))
	if len(s) > 300 {
		s = s[:300] + "..."
	}
	return s
}
