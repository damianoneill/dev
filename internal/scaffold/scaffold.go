// Package scaffold writes the standard project files created by `dev init`.
package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Params holds values interpolated into scaffold templates.
type Params struct {
	Name     string // directory / project name
	Module   string // Go module path or Python package name
	Language string // "go" or "python"
}

// ParamsFromDir derives scaffold params from an existing directory.
// It uses the directory base name as the project name, and attempts to
// infer a module path from the git remote (e.g. github.com/org/repo).
func ParamsFromDir(dir, language string) Params {
	name := filepath.Base(dir)
	module := name
	if remote := gitRemoteModule(dir); remote != "" {
		module = remote
	}
	return Params{Name: name, Module: module, Language: language}
}

// WriteCommon writes files shared across all languages into dir.
func WriteCommon(dir string, p Params) error {
	files := map[string]string{
		".gitignore":               gitignoreTmpl,
		"dev.yaml":                 devYAMLTmpl,
		"README.md":                readmeTmpl,
		".github/workflows/ci.yml": ciTmpl,
	}
	for path, tmpl := range files {
		full := filepath.Join(dir, path)
		if _, err := os.Stat(full); err == nil {
			// skip files that already exist — don't overwrite
			continue
		}
		if err := writeTemplate(full, tmpl, p); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// WriteFile writes a single file from a template string into path,
// skipping if it already exists.
func WriteFile(path, tmpl string, p Params) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return writeTemplate(path, tmpl, p)
}

func writeTemplate(path, tmpl string, p Params) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, p)
}

// gitRemoteModule tries to derive a module path from the git remote URL.
// e.g. git@github.com:org/repo.git → github.com/org/repo
func gitRemoteModule(dir string) string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	raw := strings.TrimSpace(string(out))
	// SSH: git@github.com:org/repo.git
	raw = strings.TrimPrefix(raw, "git@")
	raw = strings.ReplaceAll(raw, ":", "/")
	// HTTPS: https://github.com/org/repo.git
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimSuffix(raw, ".git")
	return raw
}
