package python

import (
	"context"
	"fmt"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
)

func init() {
	language.Register(&Python{})
}

// Python implements language.Language for Python projects.
type Python struct{}

func (p *Python) Name() string { return "python" }

func (p *Python) DefaultTasks() map[string]config.Task {
	return map[string]config.Task{
		"build":    {Cmd: "python -m build"},
		"run":      {Cmd: "python ."},
		"test":     {Cmd: "pytest"},
		"lint":     {Cmd: "ruff check ."},
		"fmt":      {Cmd: "ruff format ."},
		"clean":    {Cmd: "find . -type d -name __pycache__ -exec rm -rf {} +"},
		"setup":    {Cmd: "pip install -e .[dev]"},
		"sync":     {Cmd: "uv sync"},
		"trivy":    {Cmd: "trivy fs ."},
		"opengrep": {Cmd: "opengrep scan ."},
		"scan":     {Deps: []string{"trivy", "opengrep"}},
		"coverage": {Cmd: "pytest --cov --cov-fail-under=75"},
		"ci":       {Deps: []string{"lint", "test", "build"}},
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

func (p *Python) Sync(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "uv sync", nil)
}

func (p *Python) Scan(ctx context.Context, ex executor.Executor) error {
	if err := ex.Run(ctx, "trivy fs .", nil); err != nil {
		return err
	}
	return ex.Run(ctx, "opengrep scan .", nil)
}

func (p *Python) Coverage(ctx context.Context, ex executor.Executor, minCoverage float64) error {
	return ex.Run(ctx, fmt.Sprintf("pytest --cov --cov-fail-under=%.1f", minCoverage), nil)
}
