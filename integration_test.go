package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// devBinary builds the dev binary once per test run and returns its path.
func devBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "dev")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))
	return bin
}

func run(t *testing.T, bin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			code = 1
		}
	}
	return outBuf.String(), errBuf.String(), code
}

func golden(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "golden", name))
	require.NoError(t, err, "missing golden file %s — run tests with -update to generate", name)
	return string(data)
}

func TestVersion(t *testing.T) {
	bin := devBinary(t)
	stdout, _, code := run(t, bin, "version")
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout, "dev version")
}

func TestHelp(t *testing.T) {
	bin := devBinary(t)
	stdout, _, _ := run(t, bin)
	assert.Equal(t, golden(t, "help.txt"), stdout)
}

func TestBuildDryRun_Go(t *testing.T) {
	bin := devBinary(t)
	cwd, _ := filepath.Abs("testdata/projects/go_project")
	stdout, _, code := run(t, bin, "build", "--dry-run", "--cwd", cwd)
	assert.Equal(t, 0, code)
	assert.Equal(t, golden(t, "build_dry_run_go.txt"), stdout)
}

func TestCIDryRun_Go(t *testing.T) {
	bin := devBinary(t)
	cwd, _ := filepath.Abs("testdata/projects/go_project")
	stdout, _, code := run(t, bin, "ci", "--dry-run", "--cwd", cwd)
	assert.Equal(t, 0, code)
	assert.Equal(t, golden(t, "ci_dry_run_go.txt"), stdout)
}

func TestBuildDryRun_Python(t *testing.T) {
	bin := devBinary(t)
	cwd, _ := filepath.Abs("testdata/projects/python_project")
	stdout, _, code := run(t, bin, "build", "--dry-run", "--cwd", cwd)
	assert.Equal(t, 0, code)
	assert.Equal(t, golden(t, "build_dry_run_python.txt"), stdout)
}

func TestCIDryRun_Python(t *testing.T) {
	bin := devBinary(t)
	cwd, _ := filepath.Abs("testdata/projects/python_project")
	stdout, _, code := run(t, bin, "ci", "--dry-run", "--cwd", cwd)
	assert.Equal(t, 0, code)
	assert.Equal(t, golden(t, "ci_dry_run_python.txt"), stdout)
}

func TestExitCode_UnknownCommand(t *testing.T) {
	bin := devBinary(t)
	_, _, code := run(t, bin, "doesnotexist")
	assert.NotEqual(t, 0, code)
}

func TestExitCode_UnknownLanguage(t *testing.T) {
	bin := devBinary(t)
	dir := t.TempDir() // no dev.yaml, no indicator files → language = "unknown"
	_, stderr, code := run(t, bin, "build", "--cwd", dir)
	assert.NotEqual(t, 0, code)
	assert.Contains(t, stderr, "unsupported language")
}

func TestCwdFlag(t *testing.T) {
	bin := devBinary(t)
	cwd, _ := filepath.Abs("testdata/projects/go_project")
	// Run from a completely different directory; --cwd should redirect to the project.
	_, _, code := run(t, bin, "build", "--dry-run", "--cwd", cwd)
	assert.Equal(t, 0, code)
}
