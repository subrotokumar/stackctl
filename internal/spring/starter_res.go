package spring

type SpringInitializrResponse struct {
	Links                   map[string]Link    `json:"_links"`
	Dependencies            DependenciesConfig `json:"dependencies"`
	Type                    SimpleConfigOption `json:"type"`
	Packaging               SimpleConfigOption `json:"packaging"`
	JavaVersion             SimpleConfigOption `json:"javaVersion"`
	Language                SimpleConfigOption `json:"language"`
	BootVersion             SimpleConfigOption `json:"bootVersion"`
	GroupID                 TextConfig         `json:"groupId"`
	ArtifactID              TextConfig         `json:"artifactId"`
	Version                 TextConfig         `json:"version"`
	Name                    TextConfig         `json:"name"`
	Description             TextConfig         `json:"description"`
	PackageName             TextConfig         `json:"packageName"`
	ConfigurationFileFormat SimpleConfigOption `json:"configurationFileFormat"`
}

// Link represents a single hypermedia link
type Link struct {
	Href      *string `json:"href,omitempty"`
	Title     *string `json:"title,omitempty"`
	Templated *bool   `json:"templated,omitempty"`
}

// DependenciesConfig represents the dependencies configuration
type DependenciesConfig struct {
	Type   string            `json:"type"`
	Values []DependencyGroup `json:"values"`
}

// DependencyGroup represents a group of related dependencies
type DependencyGroup struct {
	Name   string             `json:"name"`
	Values []DependencyDetail `json:"values"`
}

// DependencyDetail represents detailed information about a dependency
type DependencyDetail struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	VersionRange *string         `json:"versionRange,omitempty"`
	Links        *map[string]any `json:"_links,omitempty"`
	Tag          string          `json:"-"`
	Selected     bool            `json:"-"`
}

// ConfigOption represents a configuration option with multiple values
type ConfigOption struct {
	Type    string        `json:"type"`
	Default string        `json:"default"`
	Values  []ConfigValue `json:"values"`
}

// ConfigValue represents a single configuration value
type ConfigValue struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Action      string            `json:"action,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// SimpleConfigOption represents a simple configuration option
type SimpleConfigOption struct {
	Type    string              `json:"type"`
	Default string              `json:"default"`
	Values  []SimpleConfigValue `json:"values"`
}

// SimpleConfigValue represents a simple configuration value
type SimpleConfigValue struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TextConfig represents a text input configuration
type TextConfig struct {
	Type    string        `json:"type"`
	Default string        `json:"default"`
	Values  *[]IdNamePair `json:"values,omitempty"`
}

type IdNamePair struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DependencyReference struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated,omitempty"`
}

type DependencyGuide struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

func (r *SpringInitializrResponse) GetDependenciesByGroup() map[string][]DependencyDetail {
	result := make(map[string][]DependencyDetail)
	for _, group := range r.Dependencies.Values {
		result[group.Name] = group.Values
	}
	return result
}

func (r *SpringInitializrResponse) FindDependencyByID(id string) *DependencyDetail {
	for _, group := range r.Dependencies.Values {
		for _, dep := range group.Values {
			if dep.ID == id {
				return &dep
			}
		}
	}
	return nil
}

func (r *SpringInitializrResponse) GetAllDependencies() []DependencyDetail {
	var deps []DependencyDetail
	for _, group := range r.Dependencies.Values {
		deps = append(deps, group.Values...)
	}
	return deps
}

func (r *SpringInitializrResponse) GetDependencyIDs() []string {
	var ids []string
	for _, group := range r.Dependencies.Values {
		for _, dep := range group.Values {
			ids = append(ids, dep.ID)
		}
	}
	return ids
}

func (r *SpringInitializrResponse) GetJavaVersions() []string {
	var versions []string
	for _, v := range r.JavaVersion.Values {
		versions = append(versions, v.ID)
	}
	return versions
}

func (r *SpringInitializrResponse) GetBootVersions() []string {
	var versions []string
	for _, v := range r.BootVersion.Values {
		versions = append(versions, v.Name)
	}
	return versions
}

func (r *SpringInitializrResponse) GetLanguages() []string {
	var languages []string
	for _, l := range r.Language.Values {
		languages = append(languages, l.ID)
	}
	return languages
}

func (r *SpringInitializrResponse) GetBuildTypes() []string {
	var types []string
	for _, p := range r.Packaging.Values {
		types = append(types, p.ID)
	}
	return types
}

func (r *SpringInitializrResponse) GetProjectTypes() []string {
	var types []string
	for _, p := range r.Type.Values {
		types = append(types, p.Name)
	}
	return types
}

func (r *SpringInitializrResponse) GetPackagingTypes() []string {
	var types []string
	for _, p := range r.Packaging.Values {
		types = append(types, p.ID)
	}
	return types
}

func (r *SpringInitializrResponse) GetConfigurationFileFormat() []string {
	types := []string{"Properties", "YAML"}
	return types
}
