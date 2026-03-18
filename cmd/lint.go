package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run linters",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "lint") {
			return runTask(cmd.Context(), "lint", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Lint(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(lintCmd) }
