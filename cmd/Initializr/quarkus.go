package initializr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/extension"
	"github.com/subrotokumar/stackctl/internal/init/quarkus"
)

const presetNone = "None (skip)"

func applyPreset(current, previous, next []string, known map[string]bool) []string {
	remove := make(map[string]bool, len(previous))
	for _, id := range previous {
		remove[id] = true
	}

	seen := make(map[string]bool)
	out := make([]string, 0, len(current)+len(next))
	for _, id := range current {
		if remove[id] || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	for _, id := range next {
		if !known[id] || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func QuarkusStarter() {
	fmt.Println(core.GreyStyle.Render("Fetching Quarkus metadata..."))
	starter, err := quarkus.Run()
	if err != nil {
		fmt.Println(core.RedStyle.Render("Error: " + err.Error()))
		return
	}

	initializr := quarkus.ProjectInitializr{
		Group:       starter.Group,
		Artifact:    starter.Artifact,
		BuildTool:   starter.BuildTool[0],
		Version:     starter.Version,
		JavaVersion: starter.JavaVersion[0],
		StarterCode: true,
		Extension:   []string{},
	}

	var steps []Step
	starterCode := "Yes"
	presetChoice := presetNone

	if core.EnableMetadataInput {
		steps = append(steps,
			func() bool { return AskText("Group", &initializr.Group) },
			func() bool { return AskText("Artifact ID", &initializr.Artifact) },
			func() bool { return AskSelect("Build Tool", starter.BuildTool, &initializr.BuildTool) },
			func() bool { return AskText("Version", &initializr.Version) },
			func() bool { return AskSelect("Java Version", starter.JavaVersion, &initializr.JavaVersion) },
			func() bool {
				if AskSelect("Starter Code", []string{"Yes", "No"}, &starterCode) {
					return true
				}
				initializr.StarterCode = strings.EqualFold(starterCode, "yes")
				return false
			},
		)
	}

	// Optional preset step (non-fatal if the API is unreachable).
	if core.EnableExtensionSelection {
		presets, err := quarkus.FetchPresets()
		if err != nil {
			fmt.Println(core.GreyStyle.Render("Skipping presets: " + err.Error()))
		} else if len(presets) > 0 {
			known := make(map[string]bool, len(starter.Extensions))
			for _, e := range starter.Extensions {
				known[e.ID] = true
			}

			titles := []string{presetNone}
			byTitle := map[string]quarkus.Preset{}
			for _, p := range presets {
				titles = append(titles, p.Title)
				byTitle[p.Title] = p
			}

			var applied []string
			steps = append(steps, func() bool {
				if AskSelect("Preset (optional)", titles, &presetChoice) {
					return true
				}
				var next []string
				if p, ok := byTitle[presetChoice]; ok {
					next = p.Extensions
				}
				initializr.Extension = applyPreset(initializr.Extension, applied, next, known)
				applied = next
				return false
			})
		}
	}

	if core.EnableExtensionSelection {
		steps = append(steps, func() bool {
			exts, back := extension.New(starter.Extensions).
				WithSelected(initializr.Extension).
				RunWithBack()
			if back {
				return true
			}
			sort.Strings(exts)
			initializr.Extension = exts
			return false
		})
	}

	RunSteps(steps)

	if core.EnableMetadataInput {
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Metadata: "))
		PrintKV([][2]string{
			{"Group", initializr.Group},
			{"Artifact ID", initializr.Artifact},
			{"Build Tool", initializr.BuildTool},
			{"Version", initializr.Version},
			{"Java Version", initializr.JavaVersion},
			{"Starter Code", starterCode},
		})
	}

	if core.EnableExtensionSelection && presetChoice != presetNone {
		fmt.Printf("\n%s %s\n", core.QuestionStyle.Render("Preset:"), presetChoice)
	}

	if core.EnableExtensionSelection {
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Selected Extensions:"))
		if len(initializr.Extension) == 0 {
			fmt.Println("  (none)")
		}
		for _, dep := range initializr.Extension {
			fmt.Printf("  - %s\n", dep)
		}
	}

	fmt.Printf("\n%s\n", core.GreyStyle.Render("Generating project..."))
	if err := initializr.Generate(); err != nil {
		fmt.Println(core.RedStyle.Render("Error: " + err.Error()))
		return
	}

	runCmd := "./mvnw quarkus:dev"
	if initializr.BuildTool != quarkus.BuildMaven {
		runCmd = "./gradlew quarkusDev"
	}
	fmt.Printf("\n%s\n", core.EndingMsgStyle.Render("Quarkus project created in ./"+initializr.Artifact))
	fmt.Printf("%s\n", core.TipMsgStyle.Render("Next: cd "+initializr.Artifact+" && stackctl run   (or "+runCmd+")"))
}
