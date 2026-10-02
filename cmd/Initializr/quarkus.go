package initializr

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/extension"
	"github.com/subrotokumar/stackctl/internal/quarkus"
)

func QuarkusStarter() {
	fmt.Println(core.GreyStyle.Render("Fetching Quarkus metadata..."))
	starter, err := quarkus.Run()
	if err != nil {
		fmt.Println(core.RedStyle.Render("❌ " + err.Error()))
		os.Exit(1)
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

	if core.EnableExtensionSelection {
		steps = append(steps, func() bool {
			// WithSelected restores earlier picks when coming back to this Step.
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
		fmt.Println(core.RedStyle.Render("❌ " + err.Error()))
		os.Exit(1)
	}

	runCmd := "./mvnw quarkus:dev"
	if initializr.BuildTool != quarkus.BuildMaven {
		runCmd = "./gradlew quarkusDev"
	}
	fmt.Printf("\n%s\n", core.EndingMsgStyle.Render("✅ Quarkus project created in ./"+initializr.Artifact))
	fmt.Printf("%s\n", core.TipMsgStyle.Render("Next: cd "+initializr.Artifact+" && stackctl run   (or "+runCmd+")"))
}
