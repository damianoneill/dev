package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/output"

	// Register language implementations.
	_ "github.com/damianoneill/dev/internal/language/golang"
	_ "github.com/damianoneill/dev/internal/language/python"
)

// AppContext is threaded through all subcommands via cobra's context.
type AppContext struct {
	Config   *config.Config
	Executor executor.Executor
	Out      *output.Writer
	DryRun   bool
	Verbose  bool
}

type ctxKey struct{}

func appCtx(cmd *cobra.Command) *AppContext {
	return cmd.Context().Value(ctxKey{}).(*AppContext)
}

var (
	flagVerbose bool
	flagDryRun  bool
	flagCwd     string
)

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "A consistent developer CLI across all your projects",
	Long: `dev provides a unified command interface for build, test, lint, fmt,
and more — regardless of language or project structure.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dir := flagCwd
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		out := output.New(flagVerbose)

		cfg, err := config.Load(dir)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		var ex executor.Executor
		if flagDryRun {
			ex = executor.NewDryRun(out)
		} else {
			ex = executor.New(out, dir)
		}

		ac := &AppContext{
			Config:   cfg,
			Executor: ex,
			Out:      out,
			DryRun:   flagDryRun,
			Verbose:  flagVerbose,
		}
		// Store AppContext on the cobra command's context for subcommands to retrieve.
		ctx := context.WithValue(cmd.Context(), ctxKey{}, ac)
		cmd.SetContext(ctx)
		return nil
	},
}

// newOutput creates an output.Writer from the current flag state.
func newOutput() *output.Writer {
	return output.New(flagVerbose)
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "print commands without executing")
	rootCmd.PersistentFlags().StringVar(&flagCwd, "cwd", "", "working directory (default: current directory)")
}
