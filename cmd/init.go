package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	initializr "github.com/subrotokumar/stackctl/cmd/Initializr"
	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/cmd/ui/selector"
)

type ProjectType string

var initCmd = &cobra.Command{
	Use:   "init <project-type>",
	Short: "Initialize a new Java project",
	Long: `Initialize a new Java project with the specified configuration.

You must specify a project type as an argument.

Supported project types:
  - spring | springboot
  - quarkus
  - micronaut

Example:
  stackctl init spring
  stackctl init quarkus
  stackctl init micronaut

Example:
  stackctl init spring
  stackctl init quarkus
`,
	Args:    cobra.MinimumNArgs(0),
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		if core.ShowLogo {
			block := core.LogoStyle.Render(strings.Join(core.Logo, "\n"))
			fmt.Print("\r\n" + strings.ReplaceAll(block, "\n", "\r\n") + "\r\n")
		}
		var project string
		if len(args) == 1 {
			project = strings.ToLower(args[0])
		} else {
			choice, back := selector.New("Project type", []string{"Spring Boot", "Quarkus", "Micronaut"}).RunWithBack()
			if back {
				os.Exit(0)
			}
			switch choice {
			case "Quarkus":
				project = "quarkus"
			case "Micronaut":
				project = "micronaut"
			default:
				project = "spring"
			}
		}
		switch project {
		case "spring", "springboot":
			fmt.Print(strings.ReplaceAll(core.GreenStyle.Render("Spring"), "\n", "\r\n") + "\r\n")
			initializr.SpringStarter()
		case "quarkus":
			fmt.Print(strings.ReplaceAll(core.GreenStyle.Render("Quarkus"), "\n", "\r\n") + "\r\n")
			initializr.QuarkusStarter()
		case "micronaut":
			fmt.Print(strings.ReplaceAll(core.GreenStyle.Render("Micronaut"), "\n", "\r\n") + "\r\n")
			initializr.MicronautStarter()
		default:
			fmt.Println(core.RedStyle.Render("Unknown project type: " + project))
			fmt.Println("Supported types: spring, springboot, quarkus")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
