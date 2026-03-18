package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/task"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if t, ok := ac.Config.Project.Tasks["test"]; ok && t.Cmd != "" {
			return task.New(ac.Config.Project.Tasks, ac.Executor).Run(cmd.Context(), "test")
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Test(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(testCmd) }
