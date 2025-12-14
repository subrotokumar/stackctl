package quarkus

type Preset struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	Icon       string   `json:"icon"`
	Extensions []string `json:"extensions"`
}

type Extension struct {
	ID                   string   `json:"id"`
	ShortID              string   `json:"shortId"`
	Version              string   `json:"version"`
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	ShortName            string   `json:"shortName"`
	Category             string   `json:"category"`
	TransitiveExtensions []string `json:"transitiveExtensions"`
	Tags                 []string `json:"tags"`
	Keywords             []string `json:"keywords"`
	ProvidesExampleCode  bool     `json:"providesExampleCode"`
	ProvidesCode         bool     `json:"providesCode"`
	Guide                string   `json:"guide"`
	Order                int      `json:"order"`
	Platform             bool     `json:"platform"`
	BOM                  string   `json:"bom"`
	Selected             bool     `json:"-"`
}

type QuarkusStarterResponse struct {
	Group       string
	Artifact    string
	BuildTool   []string
	Version     string
	JavaVersion []string
	StarterCode bool
	Extensions  []Extension `json:"extensions"`
	Presets     []Preset    `json:"presents"`
}
