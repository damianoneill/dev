package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format source code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "fmt") {
			return runTask(cmd.Context(), "fmt", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Fmt(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(fmtCmd) }
