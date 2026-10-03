package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/cmd/core"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version information",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		if core.ShowLogo {
			block := core.LogoStyle.Render(strings.Join(core.Logo, "\n"))
			fmt.Print("\r\n" + strings.ReplaceAll(block, "\n", "\r\n") + "\r\n")
		}
		fmt.Printf("Version: %s\n", core.Version)
		fmt.Printf("Creator: %s (%s)\n", core.Creator, core.Website)
		fmt.Printf("Source Code: %s\n", core.SourceCode)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
