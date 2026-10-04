package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/v4/cmd/core"
	"github.com/subrotokumar/stackctl/v4/internal/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := task.FindTaskFile(taskFile)
		if err != nil {
			return err
		}

		tf, err := task.LoadTaskfile(path)
		if err != nil {
			return err
		}
		fmt.Println(core.PurpleStyle.Render("Available tasks:"))
		for name, t := range tf.Tasks {
			fmt.Printf("%-20s %s\n", name, core.GreyStyle.Render(t.Desc))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
