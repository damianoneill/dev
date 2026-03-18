package golang

import (
	"context"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
)

func init() {
	language.Register(&Go{})
}

// Go implements language.Language for Go projects.
type Go struct{}

func (g *Go) Name() string { return "go" }

func (g *Go) DefaultTasks() map[string]config.Task {
	return map[string]config.Task{
		"build":    {Cmd: "go build ./..."},
		"run":      {Cmd: "go run ."},
		"test":     {Cmd: "go test ./..."},
		"lint":     {Cmd: "golangci-lint run"},
		"fmt":      {Cmd: "gofmt -w ."},
		"clean":    {Cmd: "go clean ./..."},
		"setup":    {Cmd: "go mod download"},
		"sync":     {Cmd: "go mod tidy"},
		"trivy":    {Cmd: "trivy fs ."},
		"opengrep": {Cmd: "opengrep scan ."},
		"scan":     {Deps: []string{"trivy", "opengrep"}},
		"ci":       {Deps: []string{"lint", "test", "build"}},
	}
}

func (g *Go) Build(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "go build ./...", nil)
}

func (g *Go) Run(ctx context.Context, ex executor.Executor, args []string) error {
	cmd := "go run ."
	if len(args) > 0 {
		cmd = "go run " + args[0]
	}
	return ex.Run(ctx, cmd, nil)
}

func (g *Go) Test(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "go test ./...", nil)
}

func (g *Go) Lint(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "golangci-lint run", nil)
}

func (g *Go) Fmt(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "gofmt -w .", nil)
}

func (g *Go) Clean(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "go clean ./...", nil)
}

func (g *Go) Setup(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "go mod download", nil)
}

func (g *Go) Sync(ctx context.Context, ex executor.Executor) error {
	return ex.Run(ctx, "go mod tidy", nil)
}

func (g *Go) Scan(ctx context.Context, ex executor.Executor) error {
	if err := ex.Run(ctx, "trivy fs .", nil); err != nil {
		return err
	}
	return ex.Run(ctx, "opengrep scan .", nil)
}
