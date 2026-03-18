package python_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
	_ "github.com/damianoneill/dev/internal/language/python"
	"github.com/damianoneill/dev/internal/output"
)

func dryRun() *executor.DryRun {
	return executor.NewDryRun(output.NewWithWriters(&bytes.Buffer{}, &bytes.Buffer{}, false))
}

func pyLang(t *testing.T) language.Language {
	t.Helper()
	lang, err := language.Resolve("python")
	require.NoError(t, err)
	return lang
}

func TestPython_Name(t *testing.T) {
	assert.Equal(t, "python", pyLang(t).Name())
}

func TestPython_DefaultTasks(t *testing.T) {
	tasks := pyLang(t).DefaultTasks()
	assert.Equal(t, "python -m build", tasks["build"].Cmd)
	assert.Equal(t, "python .", tasks["run"].Cmd)
	assert.Equal(t, "pytest", tasks["test"].Cmd)
	assert.Equal(t, "ruff check .", tasks["lint"].Cmd)
	assert.Equal(t, "ruff format .", tasks["fmt"].Cmd)
	assert.Equal(t, "uv sync", tasks["sync"].Cmd)
	assert.Equal(t, []string{"trivy", "opengrep"}, tasks["scan"].Deps)
	assert.Equal(t, []string{"lint", "test", "build"}, tasks["ci"].Deps)
}

func TestPython_Methods_DryRun(t *testing.T) {
	lang := pyLang(t)
	ex := dryRun()
	ctx := context.Background()

	require.NoError(t, lang.Build(ctx, ex))
	require.NoError(t, lang.Test(ctx, ex))
	require.NoError(t, lang.Lint(ctx, ex))
	require.NoError(t, lang.Fmt(ctx, ex))
	require.NoError(t, lang.Clean(ctx, ex))
	require.NoError(t, lang.Setup(ctx, ex))
	require.NoError(t, lang.Sync(ctx, ex))
	require.NoError(t, lang.Scan(ctx, ex))
	require.NoError(t, lang.Run(ctx, ex, nil))
	require.NoError(t, lang.Run(ctx, ex, []string{"myapp"}))
}
