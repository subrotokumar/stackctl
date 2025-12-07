package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/springx/cmd/core"
	"github.com/subrotokumar/springx/cmd/ui/inputtext"
	"github.com/subrotokumar/springx/cmd/ui/listview"
	"github.com/subrotokumar/springx/cmd/ui/selector"
	"github.com/subrotokumar/springx/internal/spring"
)

const logo = `
 ____             _
/ ___| _ __  _ __(_)_ __   __ ___  __
\___ \| '_ \| '__| | '_ \ / _‛ \ \/ /
 ___) | |_) | |  | | | | | (_| |>  <
|____/| .__/|_|  |_|_| |_|\__, /_/\_\
      |_|                  |___/
`

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
		fmt.Println(core.LogoStyle.Render(logo))

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

		if true {

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

		if false {
			listview.New(initializr.Dependencies.Values).Run()
		}

		if err := projectInitializr.Starter(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
