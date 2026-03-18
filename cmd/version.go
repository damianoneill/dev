package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is overridden at build time via -ldflags for goreleaser builds.
// For `go install` builds, resolveVersion() reads it from embedded build info.
var version = "dev"

func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the dev version",
	// version does not need config; skip PersistentPreRunE.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("dev version", resolveVersion())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
