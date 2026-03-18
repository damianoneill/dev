package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/damianoneill/dev/internal/output"
)

// Executor runs shell commands.
type Executor interface {
	Run(ctx context.Context, cmd string, env map[string]string) error
}

// Real is an Executor that runs commands via os/exec.
type Real struct {
	out *output.Writer
	dir string
}

// New returns a Real executor.
func New(out *output.Writer, dir string) *Real {
	return &Real{out: out, dir: dir}
}

func (r *Real) Run(ctx context.Context, cmd string, env map[string]string) error {
	r.out.Command(cmd)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	c := exec.CommandContext(ctx, parts[0], parts[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = r.dir
	c.Env = mergeEnv(os.Environ(), env)

	if err := c.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w", cmd, err)
	}
	return nil
}

// DryRun is an Executor that prints commands without executing them.
type DryRun struct {
	out *output.Writer
}

// NewDryRun returns a DryRun executor.
func NewDryRun(out *output.Writer) *DryRun {
	return &DryRun{out: out}
}

func (d *DryRun) Run(_ context.Context, cmd string, _ map[string]string) error {
	d.out.Info("[dry-run] " + cmd)
	return nil
}

func mergeEnv(base []string, extra map[string]string) []string {
	result := make([]string, len(base))
	copy(result, base)
	for k, v := range extra {
		result = append(result, k+"="+v)
	}
	return result
}
