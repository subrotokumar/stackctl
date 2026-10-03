package micronaut

const (
	DefaultGroup    = "com.example"
	DefaultArtifact = "demo"

	BuildMaven        = "MAVEN"
	BuildGradle       = "GRADLE"
	BuildGradleKotlin = "GRADLE_KOTLIN"
)

// Option is a label shown to the user and the value sent to the API.
type Option struct {
	Label string
	Value string
}

var (
	AppTypes = []Option{
		{"Application", "DEFAULT"},
		{"Command Line Application", "CLI"},
		{"Serverless Function", "FUNCTION"},
		{"gRPC Application", "GRPC"},
		{"Messaging Application", "MESSAGING"},
	}
	Languages = []Option{
		{"Java", "JAVA"},
		{"Kotlin", "KOTLIN"},
		{"Groovy", "GROOVY"},
	}
	BuildTools = []Option{
		{"Gradle (Groovy DSL)", BuildGradle},
		{"Gradle (Kotlin DSL)", BuildGradleKotlin},
		{"Maven", BuildMaven},
	}
	TestFrameworks = []Option{
		{"JUnit", "JUNIT"},
		{"Spock", "SPOCK"},
		{"Kotest", "KOTEST"},
	}
	JavaVersions = []Option{
		{"21", "JDK_21"},
		{"17", "JDK_17"},
	}
)

func Labels(opts []Option) []string {
	out := make([]string, len(opts))
	for i, o := range opts {
		out[i] = o.Label
	}
	return out
}

// ValueOf returns the API value for a label (or the label itself if unknown).
func ValueOf(opts []Option, label string) string {
	for _, o := range opts {
		if o.Label == label {
			return o.Value
		}
	}
	return label
}

// Feature is a single entry of /application-types/{type}/features.
type Feature struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Preview     bool   `json:"preview"`
	Community   bool   `json:"community"`
}

type featuresResponse struct {
	Features []Feature `json:"features"`
}

// Project is the user's final configuration.
type Project struct {
	Type        string
	Group       string
	Artifact    string
	Lang        string
	Build       string
	Test        string
	JavaVersion string
	Features    []string
}
