package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Global is the schema for ~/.dev/config.yaml.
type Global struct {
	Defaults map[string]string `yaml:"defaults"`
	Env      map[string]string `yaml:"env"`
}

// LoadGlobal reads ~/.dev/config.yaml, returning an empty Global if absent.
func LoadGlobal() (*Global, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Global{}, nil
	}
	path := filepath.Join(home, ".dev", "config.yaml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Global{}, nil
	}
	if err != nil {
		return nil, err
	}
	var g Global
	if err := yaml.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	return &g, nil
}
