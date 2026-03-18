package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/doctor"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		ac := appCtx(cmd)
		results := doctor.Run(cmd.Context())
		allOK := true
		for _, r := range results {
			if r.OK {
				ac.Out.Success(fmt.Sprintf("%-20s %s", r.Name, r.Message))
			} else {
				ac.Out.Error(fmt.Sprintf("%-20s %s", r.Name, r.Message))
				if r.Fix != "" {
					ac.Out.Info("  fix: " + r.Fix)
				}
				allOK = false
			}
		}
		if !allOK {
			return fmt.Errorf("some checks failed")
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(doctorCmd) }
