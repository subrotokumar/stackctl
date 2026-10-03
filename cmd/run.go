package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/cmd/ui/picker"
	"github.com/subrotokumar/stackctl/internal/task"
)

var (
	taskFile  string
	runDry    bool
	runSilent bool
)

var runCmd = &cobra.Command{
	Use:   "run [task] [VAR=value ...] [-- extra args]",
	Short: "Run a task from the task file (no task name opens an interactive selector)",
	Args:  cobra.ArbitraryArgs,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		path, err := task.FindTaskFile(taskFile)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		tf, err := task.LoadTaskfile(path)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, 0, len(tf.Tasks))
		for n := range tf.Tasks {
			names = append(names, n)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		extra := []string{}
		if dash := cmd.ArgsLenAtDash(); dash >= 0 {
			extra = args[dash:]
			args = args[:dash]
		}

		path, err := task.FindTaskFile(taskFile)
		if err != nil {
			return err
		}
		tf, err := task.LoadTaskfile(path)
		if err != nil {
			return err
		}

		// No task name -> interactive selector
		var name string
		var varArgs []string
		if len(args) == 0 {
			name, err = picker.PickTask(tf)
			if err != nil {
				return err
			}
			if name == "" { // cancelled
				return nil
			}
		} else {
			name, varArgs = args[0], args[1:]
		}

		cli := map[string]string{"CLI_ARGS": strings.Join(extra, " ")}
		for _, kv := range varArgs {
			k, v, ok := strings.Cut(kv, "=")
			if !ok {
				return fmt.Errorf("invalid argument %q (expected VAR=value)", kv)
			}
			cli[k] = v
		}

		return task.NewRunner(tf, path, runDry, runSilent).Run(name, cli)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.PersistentFlags().StringVarP(&taskFile, "file", "f", "", "task file (default: stackctl.* or Taskfile.* in current directory)")
	runCmd.Flags().BoolVarP(&runDry, "dry", "n", false, "print commands without running them")
	runCmd.Flags().BoolVarP(&runSilent, "silent", "s", false, "don't echo commands")
}
