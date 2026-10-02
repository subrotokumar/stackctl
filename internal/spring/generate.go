package spring

import (
	"archive/zip"
	"encoding/json"
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

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF5F"))

const (
	SpringInitializrMetadataURL = "https://start.spring.io/metadata/client"
)

func Run() (SpringInitializrResponse, error) {
	res, err := http.Get(SpringInitializrMetadataURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching metadata: %v\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected status code: %d\n", res.StatusCode)
		os.Exit(1)
	}

	var response SpringInitializrResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return SpringInitializrResponse{}, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return response, nil
}

func (pi ProjectInitializr) Generate() error {
	zipFile := pi.ProjectMetadata.Name + ".zip"

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

	base := "https://start.spring.io/starter.zip"
	pi.SpringBootVersion = strings.ReplaceAll(pi.SpringBootVersion, " (SNAPSHOT)", "-SNAPSHOT")

	params := url.Values{}
	params.Add("type", pi.Project)
	params.Add("language", pi.Language)
	params.Add("bootVersion", pi.SpringBootVersion)
	params.Add("groupId", pi.ProjectMetadata.GroupID)
	params.Add("artifactId", pi.ProjectMetadata.ArtifactID)
	params.Add("baseDir", pi.ProjectMetadata.Name)
	params.Add("name", pi.ProjectMetadata.Name)
	params.Add("description", pi.ProjectMetadata.Description)
	params.Add("packageName", pi.ProjectMetadata.PackageName)
	params.Add("packaging", pi.ProjectMetadata.Packaging)
	params.Add("javaVersion", pi.ProjectMetadata.JavaVersion)
	params.Add("configurationFileFormat", strings.ToLower(pi.ProjectMetadata.Configuration))
	params.Set("dependencies", strings.Join(pi.Dependencies, ","))

	url := base + "?" + params.Encode()

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#20ff93ff"))

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5FD7FF"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD75F"))

	fmt.Println()
	fmt.Println(titleStyle.Render("🔗 Spring Project Generate URL :"))
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
	zipFile := pi.ProjectMetadata.Name + ".zip"

	req, err := http.NewRequest("POST", pi.URL(), strings.NewReader(url.Values{}.Encode()))
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
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 response: %d", resp.StatusCode)
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
