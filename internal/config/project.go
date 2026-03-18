package config

// Task defines a single task entry from dev.yaml.
type Task struct {
	Cmd  string            `yaml:"cmd"`
	Deps []string          `yaml:"deps"`
	Env  map[string]string `yaml:"env"`
}

// Project is the schema for dev.yaml.
type Project struct {
	Version  string            `yaml:"version"`
	Language string            `yaml:"language"`
	Tasks    map[string]Task   `yaml:"tasks"`
	Env      map[string]string `yaml:"env"`
}
