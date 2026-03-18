package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/language"
)

// taskOrder defines the order tasks appear in the generated dev.yaml.
var taskOrder = []string{"build", "run", "test", "lint", "fmt", "clean", "setup", "sync", "trivy", "opengrep", "scan", "ci"}

var initCmd = &cobra.Command{
	Use:   "init <language>",
	Short: "Write a dev.yaml for the current project",
	Long: `Generate a dev.yaml pre-populated with every default task for the
given language. Edit the file to override any command.

Fails if dev.yaml already exists.`,
	Args: cobra.ExactArgs(1),
	// init does not need an existing config; skip PersistentPreRunE.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		langName := args[0]
		lang, err := language.Resolve(langName)
		if err != nil {
			return err
		}

		dir := flagCwd
		if dir == "" {
			dir, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		dest := filepath.Join(dir, "dev.yaml")
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("dev.yaml already exists — remove it first if you want to regenerate")
		}

		defaults := lang.DefaultTasks()
		proj := config.Project{
			Language: langName,
			Tasks:    make(map[string]config.Task, len(defaults)),
		}
		for _, name := range taskOrder {
			if t, ok := defaults[name]; ok {
				proj.Tasks[name] = t
			}
		}

		data, err := marshalDevYAML(proj)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return err
		}

		out := newOutput()
		out.Success("wrote dev.yaml — edit to customise any task")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// marshalDevYAML produces a readable dev.yaml with tasks in a stable order.
// gopkg.in/yaml.v3 marshals maps in insertion order when using yaml.Node,
// so we build the document manually.
func marshalDevYAML(p config.Project) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	root.Content = append(root.Content,
		scalarNode("language"), scalarNode(p.Language),
		scalarNode("tasks"), tasksNode(p.Tasks),
	)

	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	return yaml.Marshal(doc)
}

func tasksNode(tasks map[string]config.Task) *yaml.Node {
	m := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, name := range taskOrder {
		t, ok := tasks[name]
		if !ok {
			continue
		}
		taskMap := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		if t.Cmd != "" {
			taskMap.Content = append(taskMap.Content,
				scalarNode("cmd"), scalarNode(t.Cmd),
			)
		}
		if len(t.Deps) > 0 {
			seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Style: yaml.FlowStyle}
			for _, d := range t.Deps {
				seq.Content = append(seq.Content, scalarNode(d))
			}
			taskMap.Content = append(taskMap.Content,
				scalarNode("deps"), seq,
			)
		}
		m.Content = append(m.Content, scalarNode(name), taskMap)
	}
	return m
}

func scalarNode(val string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: val}
}
