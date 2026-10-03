package quarkus

// Extension is a single entry of https://code.quarkus.io/api/extensions.
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Extension struct {
	ID          string   `json:"id"` // e.g. io.quarkus:quarkus-rest
	Name        string   `json:"name"`
	ShortName   string   `json:"shortName"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
	Tags        []string `json:"tags"`
	Platform    bool     `json:"platform"`
	Order       int      `json:"order"`

	// Selected is UI state only.
	Selected bool `json:"-"`
}

// Preset is a single entry of https://code.quarkus.io/api/presets.
type Preset struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	Extensions []string `json:"extensions"` // extension IDs (groupId:artifactId)
}

// Starter holds the defaults and the choices shown to the user.
type Starter struct {
	Group       string
	Artifact    string
	Version     string
	BuildTool   []string // first entry is the default
	JavaVersion []string // first entry is the default
	Extensions  []Extension
}

// ProjectInitializr is the user's final configuration.
type ProjectInitializr struct {
	Group       string
	Artifact    string
	BuildTool   string
	Version     string
	JavaVersion string
	StarterCode bool
	Extension   []string // extension IDs (groupId:artifactId)
}

type stream struct {
	Recommended       bool `json:"recommended"`
	JavaCompatibility struct {
		Versions    []int `json:"versions"`
		Recommended int   `json:"recommended"`
	} `json:"javaCompatibility"`
}
