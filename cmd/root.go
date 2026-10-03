package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/cmd/core"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackctl",
	Short: "Spring Boot project initializer CLI",
	Long: `Stackctl is a fast, interactive Terminal User Interface (TUI) command-line for creating application projects from framework starters.
Select your framework, choose dependencies and project options, and generate a ready to run project directly from your terminal.
Built for a keyboard first workflow, Stackctl removes the need to open a browser just to bootstrap a new application.
Powered by Go, Bubble Tea, and official framework APIs.
No browser. No boilerplate. Just ship. ⚡`,
	Version: core.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
