package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/scaffold"
)

var initCmd = &cobra.Command{
	Use:   "init <language>",
	Short: "Initialise a project in the current directory",
	Long: `Scaffold standard project files (dev.yaml, .gitignore, README.md,
CI workflow) and run language-specific initialisation (e.g. go mod init).

Existing files are never overwritten.`,
	Args: cobra.ExactArgs(1),
	// init does not need an existing config; skip PersistentPreRunE.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		langName := args[0]
		lang, err := language.Resolve(langName)
		if err != nil {
			return err
		}

		dir := flagCwd
		if dir == "" {
			dir, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		p := scaffold.ParamsFromDir(dir, langName)
		if flagModule != "" {
			p.Module = flagModule
		}

		out := newOutput()
		ex := newExecutor(out, dir)

		out.Info(fmt.Sprintf("initialising %s project %q in %s", langName, p.Name, dir))

		if err := scaffold.WriteCommon(dir, p); err != nil {
			return err
		}
		if err := lang.Init(cmd.Context(), ex, dir, p); err != nil {
			return err
		}

		out.Success("done — run `dev setup` to install dependencies")
		return nil
	},
}

var flagModule string

func init() {
	initCmd.Flags().StringVar(&flagModule, "module", "", "module path (Go: overrides inferred path, e.g. github.com/org/repo)")
	rootCmd.AddCommand(initCmd)
}
