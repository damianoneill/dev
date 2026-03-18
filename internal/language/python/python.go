package python

import (
	"context"
	"path/filepath"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/scaffold"
)

func init() {
	language.Register(&Python{})
}

// Python implements language.Language for Python projects.
type Python struct{}

func (p *Python) Name() string { return "python" }

func (p *Python) DefaultTasks() map[string]config.Task {
	return map[string]config.Task{
		"build": {Cmd: "python -m build"},
		"test":  {Cmd: "pytest"},
		"lint":  {Cmd: "ruff check ."},
		"fmt":   {Cmd: "ruff format ."},
		"clean": {Cmd: "find . -type d -name __pycache__ -exec rm -rf {} +"},
	}
}

func (p *Python) Build(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "python -m build", nil)
}

func (p *Python) Run(ctx context.Context, ex executor.Executor, args []string) error {
	cmd := "python -m"
	if len(args) > 0 {
		cmd = "python " + args[0]
	}
	return ex.Run(ctx, cmd, nil)
}

func (p *Python) Test(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "pytest", nil)
}

func (p *Python) Lint(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "ruff check .", nil)
}

func (p *Python) Fmt(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "ruff format .", nil)
}

func (p *Python) Clean(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "find . -type d -name __pycache__ -exec rm -rf {} +", nil)
}

func (p *Python) Setup(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "pip install -e .[dev]", nil)
}

func (p *Python) Init(ctx context.Context, ex executor.Executor, dir string, params scaffold.Params) error {
	return scaffold.WriteFile(filepath.Join(dir, "pyproject.toml"), scaffold.PyprojectTmpl, params)
}
