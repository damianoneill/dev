package golang_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/language"
	_ "github.com/damianoneill/dev/internal/language/golang"
	"github.com/damianoneill/dev/internal/output"
)

func dryRun() *executor.DryRun {
	return executor.NewDryRun(output.NewWithWriters(&bytes.Buffer{}, &bytes.Buffer{}, false))
}

func goLang(t *testing.T) language.Language {
	t.Helper()
	lang, err := language.Resolve("go")
	require.NoError(t, err)
	return lang
}

func TestGo_Name(t *testing.T) {
	assert.Equal(t, "go", goLang(t).Name())
}

func TestGo_DefaultTasks(t *testing.T) {
	tasks := goLang(t).DefaultTasks()
	assert.Equal(t, "go build ./...", tasks["build"].Cmd)
	assert.Equal(t, "go run .", tasks["run"].Cmd)
	assert.Equal(t, "go test ./...", tasks["test"].Cmd)
	assert.Equal(t, "golangci-lint run", tasks["lint"].Cmd)
	assert.Equal(t, "gofmt -w .", tasks["fmt"].Cmd)
	assert.Equal(t, "go clean ./...", tasks["clean"].Cmd)
	assert.Equal(t, "go mod download", tasks["setup"].Cmd)
	assert.Equal(t, "go mod tidy", tasks["sync"].Cmd)
	assert.Equal(t, []string{"trivy", "opengrep"}, tasks["scan"].Deps)
	assert.Equal(t, []string{"lint", "test", "build"}, tasks["ci"].Deps)
}

func TestGo_Methods_DryRun(t *testing.T) {
	lang := goLang(t)
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
	require.NoError(t, lang.Run(ctx, ex, []string{"./cmd/dev"}))
}
