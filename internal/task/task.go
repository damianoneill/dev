package task

import (
	"context"
	"fmt"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
)

// Runner executes tasks in dependency order.
type Runner struct {
	tasks map[string]config.Task
	ex    executor.Executor
}

// New returns a Runner for the given task map.
func New(tasks map[string]config.Task, ex executor.Executor) *Runner {
	return &Runner{tasks: tasks, ex: ex}
}

// Run executes the named task and all its dependencies in topological order.
func (r *Runner) Run(ctx context.Context, name string) error {
	order, err := resolve(r.tasks, name)
	if err != nil {
		return err
	}
	done := map[string]bool{}
	for _, t := range order {
		if done[t] {
			continue
		}
		task := r.tasks[t]
		if task.Cmd != "" {
			if err := r.ex.Run(ctx, task.Cmd, task.Env); err != nil {
				return fmt.Errorf("task %q: %w", t, err)
			}
		}
		done[t] = true
	}
	return nil
}

// resolve returns the execution order for a task using topological sort.
func resolve(tasks map[string]config.Task, name string) ([]string, error) {
	var order []string
	visited := map[string]bool{}
	inStack := map[string]bool{}

	var visit func(n string) error
	visit = func(n string) error {
		if inStack[n] {
			return fmt.Errorf("cycle detected in task dependencies involving %q", n)
		}
		if visited[n] {
			return nil
		}
		inStack[n] = true
		t, ok := tasks[n]
		if !ok {
			return fmt.Errorf("task %q not found", n)
		}
		for _, dep := range t.Deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		inStack[n] = false
		visited[n] = true
		order = append(order, n)
		return nil
	}

	if err := visit(name); err != nil {
		return nil, err
	}
	return order, nil
}
