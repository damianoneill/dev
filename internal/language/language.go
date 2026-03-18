package language

import (
	"context"
	"fmt"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/scaffold"
)

// Language defines the contract each language implementation must satisfy.
type Language interface {
	Name() string
	DefaultTasks() map[string]config.Task
	Build(ctx context.Context, ex executor.Executor) error
	Run(ctx context.Context, ex executor.Executor, args []string) error
	Test(ctx context.Context, ex executor.Executor) error
	Lint(ctx context.Context, ex executor.Executor) error
	Fmt(ctx context.Context, ex executor.Executor) error
	Clean(ctx context.Context, ex executor.Executor) error
	Setup(ctx context.Context, ex executor.Executor) error
	// Init scaffolds language-specific files into dir for `dev init`.
	Init(ctx context.Context, ex executor.Executor, dir string, p scaffold.Params) error
}

var registry = map[string]Language{}

// Register adds a language implementation to the registry.
func Register(l Language) {
	registry[l.Name()] = l
}

// Resolve returns the Language for the given name, or an error if unknown.
func Resolve(name string) (Language, error) {
	l, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unsupported language %q — add a dev.yaml or ensure the project is detected correctly", name)
	}
	return l, nil
}
