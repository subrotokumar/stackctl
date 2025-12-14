package quarkus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	QuarkusInitializrURL = "https://code.quarkus.io"
)

func Run() (QuarkusStarterResponse, error) {
	extensions, presets, err := fetch(context.Background())
	return QuarkusStarterResponse{
		Group:       "com.acme",
		Artifact:    "code-with-quarkus",
		Version:     "1.0.0-SNAPSHOT",
		JavaVersion: []string{"21", "17"},
		BuildTool:   []string{"Maven", "Grade", "Grade With Kotlin DSL"},
		StarterCode: true,
		Extensions:  extensions,
		Presets:     presets,
	}, err
}

func fetch(ctx context.Context) ([]Extension, []Preset, error) {
	var (
		extensions []Extension
		presets    []Preset
	)

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		extensions, err = GetExtensions()
		return err
	})

	g.Go(func() error {
		var err error
		presets, err = GetPresets()
		return err
	})

	if err := g.Wait(); err != nil {
		return []Extension{}, []Preset{}, err
	}

	return extensions, presets, nil
}

func GetExtensions() ([]Extension, error) {
	res, err := http.Get(QuarkusInitializrURL + "/api/extensions?platformOnly=false")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching metadata: %v\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "GetExtensions | Unexpected status code: %d\n", res.StatusCode)
		os.Exit(1)
	}

	var response []Extension
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []Extension{}, fmt.Errorf("failed to decode extension api reponse: %w", err)
	}

	return response, nil
}

func GetPresets() ([]Preset, error) {
	res, err := http.Get(QuarkusInitializrURL + "/api/presets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching metadata: %v\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "GetPresets | Unexpected status code: %d\n", res.StatusCode)
		os.Exit(1)
	}

	var response []Preset
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []Preset{}, fmt.Errorf("failed to decode extension api reponse: %w", err)
	}

	return response, nil
}
