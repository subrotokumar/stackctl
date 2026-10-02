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

var logos = []string{
	"",
	" ████ █████  ███   ███  █   █  ███  █████ █",
	"█       █   █   █ █     █  █  █       █   █",
	" ███    █   █████ █     ███   █       █   █",
	"    █   █   █   █ █     █  █  █       █   █",
	"████    █   █   █  ███  █   █  ███    █   █████",
	"",
}

type ProjectType string

var initCmd = &cobra.Command{
	Use:   "init <project-type>",
	Short: "Initialize a new Java project",
	Long: `Initialize a new Java project with the specified configuration.

You must specify a project type as an argument.

Supported project types:
  - spring | springboot
  - quarkus

Example:
  stackctl init spring
  stackctl init quarkus
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if core.ShowLogo {
			block := core.LogoStyle.Render(strings.Join(logos, "\n"))
			fmt.Print("\r\n" + strings.ReplaceAll(block, "\n", "\r\n") + "\r\n")
		}
		var project string
		if len(args) == 1 {
			project = strings.ToLower(args[0])
		} else {
			choice, back := selector.New("Project type", []string{"Spring Boot", "Quarkus"}).RunWithBack()
			if back {
				os.Exit(0) // nothing to go back to at the first prompt
			}
			if choice == "Quarkus" {
				project = "quarkus"
			} else {
				project = "spring"
			}
		}
		switch project {
		case "spring", "springboot":
			fmt.Print(strings.ReplaceAll(core.GreenStyle.Render("SPRING"), "\n", "\r\n") + "\r\n")
			initializr.SpringStarter()
		case "quarkus":
			fmt.Print(strings.ReplaceAll(core.GreenStyle.Render("QUARKUS"), "\n", "\r\n") + "\r\n")
			initializr.QuarkusStarter()
		default:
			fmt.Println(core.RedStyle.Render("❌ Unknown project type: " + project))
			fmt.Println("Supported types: spring, springboot, quarkus")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
