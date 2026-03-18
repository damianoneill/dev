package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/task"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format source code",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if t, ok := ac.Config.Project.Tasks["fmt"]; ok && t.Cmd != "" {
			return task.New(ac.Config.Project.Tasks, ac.Executor).Run(cmd.Context(), "fmt")
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Fmt(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(fmtCmd) }
