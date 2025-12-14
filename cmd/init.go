package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/extension"
	"github.com/subrotokumar/stackctl/cmd/ui/inputtext"
	"github.com/subrotokumar/stackctl/cmd/ui/listview"
	"github.com/subrotokumar/stackctl/cmd/ui/selector"
	"github.com/subrotokumar/stackctl/internal/quarkus"
	"github.com/subrotokumar/stackctl/internal/spring"
)

const logo = `
  _________ __                 __           __  .__   
 /   _____//  |______    ____ |  | __ _____/  |_|  |  
 \_____  \\   __\__  \ _/ ___\|  |/ // ___\   __\  |  
 /        \|  |  / __ \\  \___|    <\  \___|  | |  |__
/_______  /|__| (____  /\___  >__|_ \\___  >__| |____/
        \/           \/     \/     \/    \/           
`

type ProjectType string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Spring Boot project",
	Long: `Initialize a new Spring Boot project with the specified configuration.

You can customize the project by using various flags to specify dependencies,
build tool, Java version, and other project metadata.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// reader := bufio.NewReader(os.Stdin)
		fmt.Println()
		if core.ShowLogo {
			fmt.Println(core.LogoStyle.Render(logo))
		}

		project := "spring"
		if len(args) == 1 {
			project = strings.ToLower(args[0])
		}
		switch project {
		case "spring", "springboot":
			fmt.Println(core.GreenStyle.Render("SPRING"))
			SpringStarter(cmd, args)
		case "quarkus":
			fmt.Println(core.GreenStyle.Render("QUARKUS"))
			QuarkusStarter(cmd, args)
		}
	},
}

func SpringStarter(cmd *cobra.Command, args []string) {
	initializr, err := spring.Run()
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(0)
	}

	projectMetadata := spring.ProjectMetadata{
		GroupID:       initializr.GroupID.Default,
		ArtifactID:    initializr.ArtifactID.Default,
		Name:          initializr.Name.Default,
		Description:   initializr.Description.Default,
		PackageName:   initializr.PackageName.Default,
		Packaging:     initializr.Packaging.Default,
		Configuration: initializr.ConfigurationFileFormat.Default,
		JavaVersion:   initializr.JavaVersion.Default,
	}
	projectInitializr := spring.ProjectInitializr{
		Project:           "maven-project",
		Language:          initializr.Language.Default,
		SpringBootVersion: initializr.BootVersion.Default,
	}

	if core.EnableMetadataInput {

		title := "Project"
		projectInitializr.Project = selector.New(title, []string{"maven-project", "gradle-project-kotlin", "gradle-project"}).Run()
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render(title), projectInitializr.Project)

		title = "Language"
		projectInitializr.Language = selector.New(title, initializr.GetLanguages()).Run()
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render(title), projectInitializr.Language)

		title = "Spring Boot Version"
		projectInitializr.SpringBootVersion = selector.New(title, initializr.GetBootVersions()).Run()
		fmt.Printf("%s: %s\n", core.QuestionStyle.Render(title), projectInitializr.SpringBootVersion)

		title = "Group"
		projectMetadata.GroupID = inputtext.New(title, initializr.GroupID.Default).Run()
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Metadata: "))
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.GroupID)

		title = "Artifact"
		projectMetadata.ArtifactID = inputtext.New(title, initializr.ArtifactID.Default).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.ArtifactID)

		title = "Name"
		projectMetadata.Name = inputtext.New(title, initializr.ArtifactID.Default).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.Name)

		title = "Description"
		projectMetadata.Description = inputtext.New(title, initializr.Description.Default).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.Description)

		title = "Package name"
		projectMetadata.PackageName = inputtext.New(title, initializr.PackageName.Default).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.PackageName)

		title = "Packaging"
		projectMetadata.Packaging = selector.New(title, initializr.GetPackagingTypes()).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.Packaging)

		title = "Configuration"
		projectMetadata.Configuration = selector.New(title, initializr.GetConfigurationFileFormat()).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.Configuration)

		title = "Java"
		projectMetadata.JavaVersion = selector.New(title, initializr.GetJavaVersions()).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), projectMetadata.JavaVersion)

	}

	projectInitializr.ProjectMetadata = projectMetadata
	if core.EnableDependencySelection {
		projectInitializr.Dependencies = listview.New(initializr.Dependencies.Values).Run()
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Selected Dependencies:"))
		for _, dep := range projectInitializr.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	if err := projectInitializr.Generate(); err != nil {
		panic(err)
	}
}

func QuarkusStarter(cmd *cobra.Command, args []string) {
	starter, err := quarkus.Run()
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(0)
	}

	ext := []string{}
	initializr := quarkus.ProjectInitializr{
		Group:       starter.Group,
		Artifact:    starter.Artifact,
		BuildTool:   starter.BuildTool[0],
		Version:     starter.Version,
		JavaVersion: starter.JavaVersion[0],
		StarterCode: true,
		Extension:   ext,
	}
	title := ""

	if core.EnableMetadataInput {

		title = "Group"
		initializr.Group = inputtext.New(title, initializr.Group).Run()
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Metadata: "))
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), initializr.Group)

		title = "Artifact ID"
		initializr.Artifact = inputtext.New(title, initializr.Artifact).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), initializr.Artifact)

		title = "Build Tool"
		initializr.BuildTool = selector.New(title, starter.BuildTool).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), initializr.BuildTool)

		title = "Version"
		initializr.Version = inputtext.New(title, initializr.Version).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), initializr.Version)

		title = "Java Version"
		initializr.JavaVersion = selector.New(title, starter.JavaVersion).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), initializr.JavaVersion)

		title = "Starter Code"
		starterCode := selector.New(title, []string{"Yes", "No"}).Run()
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(title), starterCode)
		if strings.ToLower(starterCode) == "yes" {
			initializr.StarterCode = true
		} else {
			initializr.StarterCode = false
		}
	}

	if core.EnableExtensionSelection {
		selectedExtension := extension.New(starter.Extensions).Run()
		initializr.Extension = selectedExtension
		fmt.Printf("\n%s\n", core.QuestionStyle.Render("Selected Dependencies:"))
		for _, dep := range selectedExtension {
			fmt.Printf("  - %s\n", dep)
		}
	}
	if err := initializr.Generate(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
