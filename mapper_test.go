package structmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMapper(t *testing.T) {
	// Initialize Mapper without options
	sm, err := NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Check if default tag name is set
	require.EqualValues(t, DefaultTagName, sm.tagName)
}

func TestNewMapperBadDefaultOption(t *testing.T) {
	// Override default options and reset them using a defer, so we do not break other tests
	defOptions := defaultOptions
	defer func() {
		defaultOptions = defOptions
	}()

	// Use invalid OptionTagName
	defaultOptions = []Option{OptionTagName("")}

	// Ensure NewMapper panics with a bad default option
	require.Panics(t, func() {
		NewMapper()
	})

}
