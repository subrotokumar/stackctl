package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackctl",
	Short: "Spring Boot project initializer CLI",
	Long: `Stackctl is a command-line tool for generating Spring Boot projects quickly.
It provides an easy way to bootstrap new Spring Boot applications with your
preferred dependencies, build tools, and project structure.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
