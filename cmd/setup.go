package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install dependencies and bootstrap the environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Setup(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(setupCmd) }
