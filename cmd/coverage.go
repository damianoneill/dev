package cmd

import (
	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
)

var minCoverage float64

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Run tests with coverage report",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		if taskDefined(ac.Config.Project.Tasks, "coverage") {
			return runTask(cmd.Context(), "coverage", ac.Config.Project.Tasks, ac.Executor)
		}
		lang, err := language.Resolve(ac.Config.Project.Language)
		if err != nil {
			return err
		}
		return lang.Coverage(cmd.Context(), ac.Executor, minCoverage)
	},
}

func init() {
	coverageCmd.Flags().Float64VarP(&minCoverage, "min-coverage", "m", 75.0, "fail if total coverage is below this percentage")
	rootCmd.AddCommand(coverageCmd)
}
