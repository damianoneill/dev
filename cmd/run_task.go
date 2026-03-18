package cmd

import (
	"context"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/executor"
	"github.com/damianoneill/dev/internal/task"
)

// taskDefined returns true if the named task exists in the project config
// with a cmd or deps — i.e. the task runner should handle it.
func taskDefined(tasks map[string]config.Task, name string) bool {
	t, ok := tasks[name]
	return ok && (t.Cmd != "" || len(t.Deps) > 0)
}

// runTask runs a named task via the task runner.
func runTask(ctx context.Context, name string, tasks map[string]config.Task, ex executor.Executor) error {
	return task.New(tasks, ex).Run(ctx, name)
}
