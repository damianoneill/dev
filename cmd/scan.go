package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run security scans (trivy + opengrep)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "scan") {
			return runTask(cmd.Context(), "scan", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Scan(cmd.Context(), ac.Executor)
	},
}

func init() { rootCmd.AddCommand(scanCmd) }
