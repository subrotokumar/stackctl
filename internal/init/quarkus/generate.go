package quarkus

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
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL     = "https://code.quarkus.io"
	clientName  = "stackctl"
	maxZipBytes = 100 << 20 // 100 MiB safety limit

	DefaultGroup    = "org.acme"
	DefaultArtifact = "code-with-quarkus"
	DefaultVersion  = "1.0.0-SNAPSHOT"
)

const (
	BuildMaven        = "MAVEN"
	BuildGradle       = "GRADLE"
	BuildGradleKotlin = "GRADLE_KOTLIN_DSL"
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

// GENERATE

var (
	groupRe    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$`)
	artifactRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

func getJSON(path string, out any) error {
	resp, err := httpClient.Get(BaseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: unexpected status %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// javaVersions returns supported Java versions, recommended one first.
func javaVersions() []string {
	fallback := []string{"21", "17"}

	var streams []stream
	if err := getJSON("/api/streams", &streams); err != nil || len(streams) == 0 {
		return fallback
	}
	chosen := streams[0]
	for _, s := range streams {
		if s.Recommended {
			chosen = s
			break
		}
	}
	vs := append([]int(nil), chosen.JavaCompatibility.Versions...)
	if len(vs) == 0 {
		return fallback
	}
	sort.Sort(sort.Reverse(sort.IntSlice(vs)))

	rec := chosen.JavaCompatibility.Recommended
	out := make([]string, 0, len(vs))
	if rec != 0 {
		out = append(out, strconv.Itoa(rec))
	}
	for _, v := range vs {
		if v != rec {
			out = append(out, strconv.Itoa(v))
		}
	}
	return out
}

// FetchPresets returns the starter presets offered by code.quarkus.io.
func FetchPresets() ([]Preset, error) {
	var presets []Preset
	if err := getJSON("/api/presets", &presets); err != nil {
		return nil, err
	}
	return presets, nil
}

// Run fetches everything needed to drive the interactive prompts.
func Run() (*Starter, error) {
	var exts []Extension
	if err := getJSON("/api/extensions", &exts); err != nil {
		return nil, fmt.Errorf("could not fetch Quarkus extensions (check your internet connection): %w", err)
	}
	if len(exts) == 0 {
		return nil, fmt.Errorf("code.quarkus.io returned no extensions")
	}
	sort.SliceStable(exts, func(i, j int) bool {
		if exts[i].Order != exts[j].Order {
			return exts[i].Order < exts[j].Order
		}
		return strings.ToLower(exts[i].Name) < strings.ToLower(exts[j].Name)
	})

	return &Starter{
		Group:       DefaultGroup,
		Artifact:    DefaultArtifact,
		Version:     DefaultVersion,
		BuildTool:   []string{BuildMaven, BuildGradle, BuildGradleKotlin},
		JavaVersion: javaVersions(),
		Extensions:  exts,
	}, nil
}

func (p ProjectInitializr) validate() error {
	switch {
	case !groupRe.MatchString(p.Group):
		return fmt.Errorf("invalid group %q (expected something like org.acme)", p.Group)
	case !artifactRe.MatchString(p.Artifact):
		return fmt.Errorf("invalid artifact id %q (letters, digits, '.', '_' and '-' only)", p.Artifact)
	case strings.TrimSpace(p.Version) == "":
		return fmt.Errorf("version must not be empty")
	}
	switch p.BuildTool {
	case BuildMaven, BuildGradle, BuildGradleKotlin:
	default:
		return fmt.Errorf("unsupported build tool %q", p.BuildTool)
	}
	if _, err := strconv.Atoi(p.JavaVersion); err != nil {
		return fmt.Errorf("invalid java version %q", p.JavaVersion)
	}
	return nil
}

// DownloadURL returns the code.quarkus.io URL that produces the project zip.
func (p ProjectInitializr) DownloadURL() string {
	q := url.Values{}
	q.Set("g", p.Group)
	q.Set("a", p.Artifact)
	q.Set("v", p.Version)
	q.Set("b", p.BuildTool)
	q.Set("j", p.JavaVersion)
	q.Set("cn", clientName)
	if !p.StarterCode {
		q.Set("nc", "true") // no starter code
	}
	for _, e := range p.Extension {
		q.Add("e", e)
	}
	return BaseURL + "/d?" + q.Encode()
}

// EXTRACT

func (p ProjectInitializr) download() ([]byte, error) {
	resp, err := httpClient.Get(p.DownloadURL())
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("code.quarkus.io responded with %s: %s",
			resp.Status, strings.TrimSpace(string(msg)))
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

func (p ProjectInitializr) Generate() error {
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
