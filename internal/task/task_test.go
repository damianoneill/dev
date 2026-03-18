package task_test

import (
	"context"
	"testing"

	"github.com/damianoneill/dev/internal/config"
	"github.com/damianoneill/dev/internal/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingExecutor struct {
	ran []string
}

func (r *recordingExecutor) Run(_ context.Context, cmd string, _ map[string]string) error {
	r.ran = append(r.ran, cmd)
	return nil
}

func TestRunner_SimpleDep(t *testing.T) {
	tasks := map[string]config.Task{
		"lint":  {Cmd: "golangci-lint run"},
		"test":  {Cmd: "go test ./..."},
		"build": {Cmd: "go build ./...", Deps: []string{"lint", "test"}},
	}
	ex := &recordingExecutor{}
	r := task.New(tasks, ex)
	require.NoError(t, r.Run(context.Background(), "build"))
	assert.Equal(t, []string{"golangci-lint run", "go test ./...", "go build ./..."}, ex.ran)
}

func TestRunner_NoDuplication(t *testing.T) {
	tasks := map[string]config.Task{
		"lint":  {Cmd: "lint"},
		"test":  {Cmd: "test", Deps: []string{"lint"}},
		"build": {Cmd: "build", Deps: []string{"lint", "test"}},
	}
	ex := &recordingExecutor{}
	r := task.New(tasks, ex)
	require.NoError(t, r.Run(context.Background(), "build"))
	assert.Equal(t, []string{"lint", "test", "build"}, ex.ran)
}

func TestRunner_CycleDetected(t *testing.T) {
	tasks := map[string]config.Task{
		"a": {Cmd: "a", Deps: []string{"b"}},
		"b": {Cmd: "b", Deps: []string{"a"}},
	}
	ex := &recordingExecutor{}
	r := task.New(tasks, ex)
	err := r.Run(context.Background(), "a")
	assert.ErrorContains(t, err, "cycle")
}

func TestRunner_UnknownTask(t *testing.T) {
	ex := &recordingExecutor{}
	r := task.New(map[string]config.Task{}, ex)
	err := r.Run(context.Background(), "missing")
	assert.ErrorContains(t, err, "not found")
}
