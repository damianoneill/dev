package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var runCmd = &cobra.Command{
	Use:   "run [args]",
	Short: "Run the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Run(cmd.Context(), ac.Executor, args)
	},
}

func init() { rootCmd.AddCommand(runCmd) }
