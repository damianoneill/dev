package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "test") {
			return runTask(cmd.Context(), "test", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Test(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(testCmd) }
