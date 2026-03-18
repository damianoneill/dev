package golang

import (
	"context"
	"path/filepath"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
	"github.com/damianoneill/dev/internal/scaffold"
)

func init() {
	language.Register(&Go{})
}

// Go implements language.Language for Go projects.
type Go struct{}

func (g *Go) Name() string { return "go" }

func (g *Go) DefaultTasks() map[string]config.Task {
	return map[string]config.Task{
		"build": {Cmd: "go build ./..."},
		"test":  {Cmd: "go test ./..."},
		"lint":  {Cmd: "golangci-lint run"},
		"fmt":   {Cmd: "gofmt -w ."},
		"clean": {Cmd: "go clean ./..."},
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

func (g *Go) Init(ctx context.Context, ex executor.Executor, dir string, p scaffold.Params) error {
	if err := ex.Run(ctx, "go mod init "+p.Module, nil); err != nil {
		return err
	}
	if err := scaffold.WriteFile(filepath.Join(dir, "main.go"), scaffold.GoMainTmpl, p); err != nil {
		return err
	}
	return scaffold.WriteFile(filepath.Join(dir, ".golangci.yml"), scaffold.GolangciTmpl, p)
}
