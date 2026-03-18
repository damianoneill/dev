package config

import (
	"os"
	"path/filepath"
)

// DetectLanguage infers the project language from well-known indicator files
// found in dir. Returns "unknown" if nothing matches.
func DetectLanguage(dir string) string {
	probes := []struct {
		file string
		lang string
	}{
		{"go.mod", "go"},
		{"pyproject.toml", "python"},
		{"setup.py", "python"},
		{"requirements.txt", "python"},
	}
	for _, p := range probes {
		if _, err := os.Stat(filepath.Join(dir, p.file)); err == nil {
			return p.lang
		}
	}
	return "unknown"
}
