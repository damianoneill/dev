package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/task"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run the full CI pipeline locally (lint → test → build)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)

		// If the project defines a ci task, use it.
		if _, ok := ac.Config.Project.Tasks["ci"]; ok {
			return task.New(ac.Config.Project.Tasks, ac.Executor).Run(cmd.Context(), "ci")
		}

		// Otherwise synthesise a ci task from lint → test → build.
		tasks := mergeTasks(ac.Config.Project.Tasks, map[string]config.Task{
			"ci": {Deps: []string{"lint", "test", "build"}},
		})
		return task.New(tasks, ac.Executor).Run(cmd.Context(), "ci")
	},
}

// mergeTasks returns a new map containing base with overrides applied.
func mergeTasks(base, overrides map[string]config.Task) map[string]config.Task {
	merged := make(map[string]config.Task, len(base)+len(overrides))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}
	return merged
}

func init() { rootCmd.AddCommand(ciCmd) }
