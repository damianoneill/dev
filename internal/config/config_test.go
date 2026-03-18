package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/damianoneill/dev/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected string
	}{
		{"go project", []string{"go.mod"}, "go"},
		{"python pyproject", []string{"pyproject.toml"}, "python"},
		{"python setup", []string{"setup.py"}, "python"},
		{"python requirements", []string{"requirements.txt"}, "python"},
		{"unknown", []string{"README.md"}, "unknown"},
		{"go takes precedence", []string{"go.mod", "pyproject.toml"}, "go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				require.NoError(t, os.WriteFile(filepath.Join(dir, f), []byte{}, 0o644))
			}
			assert.Equal(t, tt.expected, config.DetectLanguage(dir))
		})
	}
}

func TestLoad_WithDevYAML(t *testing.T) {
	dir := t.TempDir()
	yaml := `language: go
tasks:
  build:
    cmd: go build ./...
  test:
    cmd: go test ./...
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dev.yaml"), []byte(yaml), 0o644))

	cfg, err := config.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "go", cfg.Project.Language)
	assert.Equal(t, "go build ./...", cfg.Project.Tasks["build"].Cmd)
	assert.Equal(t, "go test ./...", cfg.Project.Tasks["test"].Cmd)
	assert.Equal(t, dir, cfg.Dir)
}

func TestLoad_DetectsLanguage(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/foo\n"), 0o644))

	cfg, err := config.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "go", cfg.Project.Language)
}

func TestLoad_WalksUp(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "sub", "pkg")
	require.NoError(t, os.MkdirAll(sub, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "dev.yaml"), []byte("language: go\n"), 0o644))

	cfg, err := config.Load(sub)
	require.NoError(t, err)
	assert.Equal(t, "go", cfg.Project.Language)
	assert.Equal(t, root, cfg.Dir)
}
