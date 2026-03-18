package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/task"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if t, ok := ac.Config.Project.Tasks["build"]; ok && t.Cmd != "" {
			return task.New(ac.Config.Project.Tasks, ac.Executor).Run(cmd.Context(), "build")
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Build(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(buildCmd) }
