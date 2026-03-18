package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/task"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove build artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if t, ok := ac.Config.Project.Tasks["clean"]; ok && t.Cmd != "" {
			return task.New(ac.Config.Project.Tasks, ac.Executor).Run(cmd.Context(), "clean")
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Clean(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(cleanCmd) }
