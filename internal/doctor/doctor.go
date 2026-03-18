package doctor

import (
	"context"
	"fmt"
	"os/exec"
)

// Result holds the outcome of a single check.
type Result struct {
	Name    string
	OK      bool
	Message string
	Fix     string
}

// Check is a single environment validation.
type Check interface {
	Name() string
	Run(ctx context.Context) Result
}

// Run executes all registered checks and returns the results.
func Run(ctx context.Context) []Result {
	results := make([]Result, 0, len(checks))
	for _, c := range checks {
		results = append(results, c.Run(ctx))
	}
	return results
}

var checks []Check

func register(c Check) { checks = append(checks, c) }

// binaryCheck verifies that a binary is present in PATH.
type binaryCheck struct {
	name   string
	binary string
	fix    string
}

func (b *binaryCheck) Name() string { return b.name }

func (b *binaryCheck) Run(_ context.Context) Result {
	path, err := exec.LookPath(b.binary)
	if err != nil {
		return Result{
			Name:    b.name,
			OK:      false,
			Message: fmt.Sprintf("%q not found in PATH", b.binary),
			Fix:     b.fix,
		}
	}
	return Result{
		Name:    b.name,
		OK:      true,
		Message: path,
	}
}

func init() {
	register(&binaryCheck{"go", "go", "https://go.dev/dl/"})
	register(&binaryCheck{"golangci-lint", "golangci-lint", "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"})
	register(&binaryCheck{"python", "python3", "https://www.python.org/downloads/"})
	register(&binaryCheck{"git", "git", "https://git-scm.com/downloads"})
	register(&binaryCheck{"trivy", "trivy", "https://aquasecurity.github.io/trivy/latest/getting-started/installation/"})
	register(&binaryCheck{"opengrep", "opengrep", "https://github.com/opengrep/opengrep"})
}
