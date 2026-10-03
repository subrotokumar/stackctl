package initializr

import (
	"fmt"
	"sort"

	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/feature"
	"github.com/subrotokumar/stackctl/internal/micronaut"
)

func MicronautStarter() {
	project := micronaut.Project{
		Type:        micronaut.AppTypes[0].Value,
		Group:       micronaut.DefaultGroup,
		Artifact:    micronaut.DefaultArtifact,
		Lang:        micronaut.Languages[0].Value,
		Build:       micronaut.BuildTools[0].Value,
		Test:        micronaut.TestFrameworks[0].Value,
		JavaVersion: micronaut.JavaVersions[0].Value,
		Features:    []string{},
	}

	// Labels shown in the prompts; mapped to API values after each answer.
	typeLabel := micronaut.AppTypes[0].Label
	langLabel := micronaut.Languages[0].Label
	buildLabel := micronaut.BuildTools[0].Label
	testLabel := micronaut.TestFrameworks[0].Label
	javaLabel := micronaut.JavaVersions[0].Label

	ask := func(title string, opts []micronaut.Option, label, dst *string) Step {
		return func() bool {
			if AskSelect(title, micronaut.Labels(opts), label) {
				return true
			}
			*dst = micronaut.ValueOf(opts, *label)
			return false
		}
	}

	var steps []Step

	if core.EnableMetadataInput {
		steps = append(steps,
			ask("Application Type", micronaut.AppTypes, &typeLabel, &project.Type),
			func() bool { return AskText("Name (artifact)", &project.Artifact) },
			func() bool { return AskText("Group (base package)", &project.Group) },
			ask("Language", micronaut.Languages, &langLabel, &project.Lang),
			ask("Build Tool", micronaut.BuildTools, &buildLabel, &project.Build),
			ask("Test Framework", micronaut.TestFrameworks, &testLabel, &project.Test),
			ask("Java Version", micronaut.JavaVersions, &javaLabel, &project.JavaVersion),
		)
	}

	if core.EnableExtensionSelection {
		// Features depend on the application type, so they are fetched when this
		// step runs and cached per type.
		cache := map[string][]micronaut.Feature{}
		steps = append(steps, func() bool {
			feats, ok := cache[project.Type]
			if !ok {
				fmt.Println(core.GreyStyle.Render("Fetching Micronaut features..."))
				var err error
				feats, err = micronaut.FetchFeatures(project.Type)
				if err != nil {
					fmt.Println(core.GreyStyle.Render("Skipping features: " + err.Error()))
					return false
				}
				cache[project.Type] = feats
			}

			// Drop earlier picks that don't exist for the current application type.
			known := make(map[string]bool, len(feats))
			for _, f := range feats {
				known[f.Name] = true
			}
			keep := make([]string, 0, len(project.Features))
			for _, n := range project.Features {
				if known[n] {
					keep = append(keep, n)
				}
			}

			names, back := feature.New(feats).WithSelected(keep).RunWithBack()
			if back {
				return true
			}
			sort.Strings(names)
			project.Features = names
			return false
		})
	}

	RunSteps(steps)

	if core.EnableMetadataInput {
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Metadata: "))
		PrintKV([][2]string{
			{"Application Type", typeLabel},
			{"Name", project.Artifact},
			{"Group", project.Group},
			{"Language", langLabel},
			{"Build Tool", buildLabel},
			{"Test Framework", testLabel},
			{"Java Version", javaLabel},
		})
	}

	if core.EnableExtensionSelection {
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Selected Features:"))
		if len(project.Features) == 0 {
			fmt.Println("  (none, Micronaut defaults only)")
		}
		for _, f := range project.Features {
			fmt.Printf("  - %s\n", f)
		}
	}

	fmt.Printf("\n%s\n", core.GreyStyle.Render("Generating project..."))
	if err := project.Generate(); err != nil {
		fmt.Println(core.RedStyle.Render("Error: " + err.Error()))
		return
	}

	runCmd := "./gradlew run"
	if project.Build == micronaut.BuildMaven {
		runCmd = "./mvnw mn:run"
	}
	fmt.Printf("\n%s\n", core.EndingMsgStyle.Render("Micronaut project created in ./"+project.Artifact))
	fmt.Printf("%s\n", core.TipMsgStyle.Render("Next: cd "+project.Artifact+" && stackctl run   (or "+runCmd+")"))
}
