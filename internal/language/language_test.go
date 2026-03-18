package language_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damianoneill/dev/internal/language"
	// Import language implementations to trigger their init() registrations.
	_ "github.com/damianoneill/dev/internal/language/golang"
	_ "github.com/damianoneill/dev/internal/language/python"
)

func TestResolve_Go(t *testing.T) {
	lang, err := language.Resolve("go")
	require.NoError(t, err)
	assert.Equal(t, "go", lang.Name())
}

func TestResolve_Python(t *testing.T) {
	lang, err := language.Resolve("python")
	require.NoError(t, err)
	assert.Equal(t, "python", lang.Name())
}

func TestResolve_Unknown(t *testing.T) {
	_, err := language.Resolve("rust")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rust")
}
