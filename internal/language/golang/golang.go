package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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
		"coverage": {Cmd: "go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"},
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

func (g *Go) Coverage(ctx context.Context, ex executor.Executor, minCoverage float64) error {
	if err := ex.Run(ctx, "go test -coverprofile=coverage.out ./...", nil); err != nil {
		return err
	}
	// coverage.out won't exist in dry-run mode; skip threshold check.
	if _, err := os.Stat("coverage.out"); os.IsNotExist(err) {
		return nil
	}
	pct, err := goCoverageTotal(ctx, "coverage.out")
	if err != nil {
		return err
	}
	if pct < minCoverage {
		return fmt.Errorf("coverage %.1f%% is below minimum %.1f%%", pct, minCoverage)
	}
	return nil
}

// goCoverageTotal runs "go tool cover -func" and returns the total coverage percentage.
func goCoverageTotal(ctx context.Context, profile string) (float64, error) {
	out, err := exec.CommandContext(ctx, "go", "tool", "cover", "-func="+profile).Output()
	if err != nil {
		return 0, fmt.Errorf("go tool cover: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		return 0, fmt.Errorf("go tool cover produced no output")
	}
	// Last line: "total:    (statements)    73.9%"
	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) < 3 {
		return 0, fmt.Errorf("unexpected go tool cover output: %s", lines[len(lines)-1])
	}
	pctStr := strings.TrimSuffix(fields[len(fields)-1], "%")
	return strconv.ParseFloat(pctStr, 64)
}
