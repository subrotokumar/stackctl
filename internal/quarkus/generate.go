package quarkus

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/cmd/core"
)

const (
	ParamGroup       = "g"
	ParamArtifact    = "a"
	ParamBuildTool   = "b"
	ParamExtension   = "e"
	ParamVersion     = "v"
	ParamJavaVersion = "j"
	ParamCN          = "cn"
)

var (
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF5F"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FD7FF"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD75F"))
)

func (pi ProjectInitializr) Generate() error {
	zipFile := pi.Artifact + ".zip"

	err := pi.downloadStarterZip()
	if err != nil {
		return err
	}
	fmt.Println("⌀ Downloaded:", zipFile)

	err = unzip(zipFile, "./")
	if err != nil {
		return err
	}
	fmt.Println("⌀ Extracted successfully")
	fmt.Println()

	err = os.Remove(zipFile)
	if err != nil {
		return err
	}

	fmt.Println(successStyle.Render("⌀ Project created successfully 🚀!"))

	return nil
}

func (pi ProjectInitializr) URL() string {

	base := "https://code.quarkus.io/d"

	params := url.Values{}
	params.Add(ParamGroup, pi.Group)
	params.Add(ParamArtifact, pi.Artifact)
	params.Add(ParamBuildTool, strings.ToUpper(strings.ReplaceAll(pi.BuildTool, " ", "_")))
	params.Add(ParamJavaVersion, pi.JavaVersion)
	params.Add(ParamVersion, pi.Version)

	for _, ext := range pi.Extension {
		params.Add(ParamExtension, strings.ReplaceAll(ext, "io.quarkus:quarkus-", ""))
	}

	url := base + "?" + params.Encode()

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#20ff93ff"))

	fmt.Println()
	fmt.Println(titleStyle.Render("🔗 Quarkus Project Generate URL :"))
	fmt.Printf("%s\n\n", core.GreyStyle.Render(url))

	fmt.Println(titleStyle.Render(base + "?"))
	for key, values := range params {
		for _, val := range values {
			fmt.Printf("  %s=%s\n",
				keyStyle.Render(key),
				valueStyle.Render(val),
			)
		}
	}
	fmt.Println()

	return url
}

func (pi ProjectInitializr) downloadStarterZip() error {
	zipFile := pi.Artifact + ".zip"

	req, err := http.NewRequest("GET", pi.URL(), strings.NewReader(url.Values{}.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("⌀ Status:", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read response body as string
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("non-2xx response: %d; failed to read body: %v", resp.StatusCode, readErr)
		}
		bodyStr := string(bodyBytes)
		return fmt.Errorf("non-2xx response: %d; body: %s", resp.StatusCode, bodyStr)
	}

	out, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(destFile, rc)
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
