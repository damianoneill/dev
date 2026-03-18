package executor_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/output"
)

func newWriter() *output.Writer {
	return output.NewWithWriters(&bytes.Buffer{}, &bytes.Buffer{}, false)
}

func TestDryRun_PrintsCommand(t *testing.T) {
	var out bytes.Buffer
	w := output.NewWithWriters(&out, &bytes.Buffer{}, false)
	ex := executor.NewDryRun(w)
	err := ex.Run(context.Background(), "go build ./...", nil)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "go build ./...")
	assert.Contains(t, out.String(), "dry-run")
}

func TestDryRun_DoesNotExecute(t *testing.T) {
	ex := executor.NewDryRun(newWriter())
	// This would fail if actually executed (nonexistent command).
	err := ex.Run(context.Background(), "this-command-does-not-exist", nil)
	require.NoError(t, err)
}

func TestReal_EmptyCommand(t *testing.T) {
	ex := executor.New(newWriter(), "")
	err := ex.Run(context.Background(), "", nil)
	assert.Error(t, err)
}

func TestReal_SimpleCommand(t *testing.T) {
	ex := executor.New(newWriter(), "")
	err := ex.Run(context.Background(), "true", nil)
	require.NoError(t, err)
}

func TestReal_FailingCommand(t *testing.T) {
	ex := executor.New(newWriter(), "")
	err := ex.Run(context.Background(), "false", nil)
	assert.Error(t, err)
}

func TestReal_ShellOperator(t *testing.T) {
	ex := executor.New(newWriter(), "")
	// Uses sh -c because of the pipe.
	err := ex.Run(context.Background(), "echo hello | cat", nil)
	require.NoError(t, err)
}

func TestReal_WithEnv(t *testing.T) {
	ex := executor.New(newWriter(), "")
	// mergeEnv is exercised; command itself just needs to succeed.
	err := ex.Run(context.Background(), "true", map[string]string{"DEV_TEST_VAR": "hello"})
	require.NoError(t, err)
}
