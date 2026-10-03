package initializr

import (
	"fmt"

	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/listview"
	"github.com/subrotokumar/stackctl/internal/init/spring"
)

func SpringStarter() {
	initializr, err := spring.Run()
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	md := spring.ProjectMetadata{
		GroupID:       initializr.GroupID.Default,
		ArtifactID:    initializr.ArtifactID.Default,
		Name:          initializr.ArtifactID.Default, // same default as your original code
		Description:   initializr.Description.Default,
		PackageName:   initializr.PackageName.Default,
		Packaging:     initializr.Packaging.Default,
		Configuration: initializr.ConfigurationFileFormat.Default,
		JavaVersion:   initializr.JavaVersion.Default,
	}
	pi := spring.ProjectInitializr{
		Project:           "maven-project",
		Language:          initializr.Language.Default,
		SpringBootVersion: initializr.BootVersion.Default,
	}

	var steps []Step

	if core.EnableMetadataInput {
		steps = append(steps,
			func() bool {
				return AskSelect("Project", []string{"maven-project", "gradle-project-kotlin", "gradle-project"}, &pi.Project)
			},
			func() bool { return AskSelect("Language", initializr.GetLanguages(), &pi.Language) },
			func() bool {
				return AskSelect("Spring Boot Version", initializr.GetBootVersions(), &pi.SpringBootVersion)
			},
			func() bool { return AskText("Group", &md.GroupID) },
			func() bool { return AskText("Artifact", &md.ArtifactID) },
			func() bool { return AskText("Name", &md.Name) },
			func() bool { return AskText("Description", &md.Description) },
			func() bool { return AskText("Package name", &md.PackageName) },
			func() bool { return AskSelect("Packaging", initializr.GetPackagingTypes(), &md.Packaging) },
			func() bool {
				return AskSelect("Configuration", initializr.GetConfigurationFileFormat(), &md.Configuration)
			},
			func() bool { return AskSelect("Java", initializr.GetJavaVersions(), &md.JavaVersion) },
		)
	}

	if core.EnableDependencySelection {
		steps = append(steps, func() bool {
			deps, back := listview.New(initializr.Dependencies.Values).RunWithBack()
			if back {
				return true
			}
			pi.Dependencies = deps
			return false
		})
	}

	RunSteps(steps)

	if core.EnableMetadataInput {
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render("Project"), pi.Project)
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render("Language"), pi.Language)
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render("Spring Boot Version"), pi.SpringBootVersion)
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Metadata: "))
		PrintKV([][2]string{
			{"Group", md.GroupID},
			{"Artifact", md.ArtifactID},
			{"Name", md.Name},
			{"Description", md.Description},
			{"Package name", md.PackageName},
			{"Packaging", md.Packaging},
			{"Configuration", md.Configuration},
			{"Java", md.JavaVersion},
		})
	}

	pi.ProjectMetadata = md

	if core.EnableDependencySelection {
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Selected Dependencies:"))
		for _, dep := range pi.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	if err := pi.Generate(); err != nil {
		panic(err)
	}
}
