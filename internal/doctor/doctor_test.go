package doctor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/damianoneill/dev/internal/doctor"
)

func TestRun_ReturnsResults(t *testing.T) {
	results := doctor.Run(context.Background())
	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.Name)
		// OK may be true or false depending on what's installed; just check shape.
		if r.OK {
			assert.NotEmpty(t, r.Message) // path to binary
		} else {
			assert.NotEmpty(t, r.Fix) // install hint
		}
	}
}

func TestRun_GoIsPresent(t *testing.T) {
	results := doctor.Run(context.Background())
	var goResult *doctor.Result
	for i := range results {
		if results[i].Name == "go" {
			goResult = &results[i]
			break
		}
	}
	assert.NotNil(t, goResult, "expected a 'go' check")
	assert.True(t, goResult.OK, "go binary should be present in the test environment")
}

func TestRun_GitIsPresent(t *testing.T) {
	results := doctor.Run(context.Background())
	var gitResult *doctor.Result
	for i := range results {
		if results[i].Name == "git" {
			gitResult = &results[i]
			break
		}
	}
	assert.NotNil(t, gitResult, "expected a 'git' check")
	assert.True(t, gitResult.OK, "git binary should be present in the test environment")
}
