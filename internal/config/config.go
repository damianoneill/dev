package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "dev.yaml"

// Config is the merged result of global + project configuration.
type Config struct {
	Project *Project
	Global  *Global
	// Dir is the resolved project root directory.
	Dir string
}

// Load finds dev.yaml by walking up from dir, merges with global config,
// and runs language detection if the language field is absent.
func Load(dir string) (*Config, error) {
	global, err := LoadGlobal()
	if err != nil {
		return nil, err
	}

	project, projectDir, err := findAndLoadProject(dir)
	if err != nil {
		return nil, err
	}

	if project.Language == "" {
		project.Language = DetectLanguage(projectDir)
	}

	// Apply global defaults for tasks not defined in the project.
	if project.Tasks == nil {
		project.Tasks = make(map[string]Task)
	}
	for name, cmd := range global.Defaults {
		if _, exists := project.Tasks[name]; !exists {
			project.Tasks[name] = Task{Cmd: cmd}
		}
	}

	return &Config{
		Project: project,
		Global:  global,
		Dir:     projectDir,
	}, nil
}

func findAndLoadProject(startDir string) (*Project, string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, configFileName)
		data, err := os.ReadFile(candidate)
		if err == nil {
			var p Project
			if err := yaml.Unmarshal(data, &p); err != nil {
				return nil, "", err
			}
			return &p, dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// No dev.yaml found; return empty project rooted at startDir.
	return &Project{}, startDir, nil
}
