package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dependencies with the manifest or lockfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "sync") {
			return runTask(cmd.Context(), "sync", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Sync(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(syncCmd) }
