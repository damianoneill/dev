package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove build artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "clean") {
			return runTask(cmd.Context(), "clean", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Clean(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(cleanCmd) }
