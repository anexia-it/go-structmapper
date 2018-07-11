package structmapper

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionTagName(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		sm, err := NewMapper(OptionTagName("test"))

		require.NoError(t, err)
		require.NotNil(t, sm)

		// Check if the supplied tag name was set
		require.EqualValues(t, "test", sm.tagName)
	})

	t.Run("EmptyTag", func(t *testing.T) {
		// Try setting an empty tag name
		sm, err := NewMapper(OptionTagName(""))
		require.Error(t, err)
		require.Nil(t, sm)

		require.IsType(t, &multierror.Error{}, err)
		multiErr := err.(*multierror.Error)
		assert.Len(t, multiErr.WrappedErrors(), 1)
		assert.EqualError(t, multiErr.WrappedErrors()[0], ErrTagNameEmpty.Error())
	})

}

func TestParseTag(t *testing.T) {
	t.Run("Dash", func(t *testing.T) {
		// Check if the special-case ignore-me tag ("-") gives the correct result
		name, omitEmpty, err := parseTag("-")
		assert.EqualValues(t, "-", name)
		assert.EqualValues(t, false, omitEmpty)
		assert.NoError(t, err)
	})

	t.Run("OmitEmptyNoTagName", func(t *testing.T) {
		// Check if ",omitEmpty" alone works
		name, omitEmpty, err := parseTag(",omitempty")
		assert.EqualValues(t, "", name)
		assert.EqualValues(t, true, omitEmpty)
		assert.NoError(t, err)
	})

	t.Run("OmitEmpty", func(t *testing.T) {
		// Check if "name,omitEmpty" returns the correct tag name
		name, omitEmpty, err := parseTag("test,omitempty")
		assert.EqualValues(t, "test", name)
		assert.EqualValues(t, true, omitEmpty)
		assert.NoError(t, err)
	})

	t.Run("Puncation", func(t *testing.T) {
		// Check if a punctation inside the tag name gives an error
		name, omitEmpty, err := parseTag("test.,omitempty")
		assert.EqualValues(t, "test.", name)
		assert.EqualValues(t, true, omitEmpty)
		assert.Error(t, err)

		require.IsType(t, &InvalidTag{}, err)
		invalidTagErr := err.(*InvalidTag)
		assert.EqualValues(t, "test.,omitempty", invalidTagErr.Tag())
	})

	t.Run("Whitespace", func(t *testing.T) {
		// Check if whitespace inside the tag name gives an error
		name, omitEmpty, err := parseTag("test ,omitempty")
		assert.EqualValues(t, "test ", name)
		assert.EqualValues(t, true, omitEmpty)
		assert.Error(t, err)

		require.IsType(t, &InvalidTag{}, err)
		invalidTagErr := err.(*InvalidTag)
		assert.EqualValues(t, "test ,omitempty", invalidTagErr.Tag())
	})

	t.Run("Underscores", func(t *testing.T) {
		// Check if underscores are allowed
		name, omitEmpty, err := parseTag("test_tag")
		assert.NoError(t, err)
		assert.EqualValues(t, "test_tag", name)
		assert.EqualValues(t, false, omitEmpty)
	})

}

func TestIsInvalidTag(t *testing.T) {
	// Test if IsInvalidTag works correctly for an invalid tag
	err := newErrorInvalidTag("test")
	require.NotNil(t, err)

	t.Run("Yes", func(t *testing.T) {
		it, ok := IsInvalidTag(err)
		require.EqualValues(t, true, ok, "IsInvalidTag should return true")
		require.NotNil(t, it)
		require.EqualValues(t, err, it)
	})

	t.Run("No", func(t *testing.T) {
		// Test if IsInvalidTag works correctly for non-invalid tag error
		err = errors.New("test")
		it, ok := IsInvalidTag(err)
		require.EqualValues(t, false, ok, "IsInvalidTag should return false")
		require.Nil(t, it)
	})

}

func TestInvalidTag_Tag(t *testing.T) {
	// Check if the tag is preserved
	err := newErrorInvalidTag("test tag")
	require.NotNil(t, err)

	it, ok := IsInvalidTag(err)
	require.EqualValues(t, true, ok, "IsInvalidTag should return true")
	require.NotNil(t, it)
	require.EqualValues(t, err, it)
	require.EqualValues(t, "test tag", it.Tag())
}

func TestInvalidTag_Error(t *testing.T) {
	// Check if the error message is correct
	err := newErrorInvalidTag("test tag")
	require.NotNil(t, err)
	require.EqualError(t, err, "Invalid tag: 'test tag'")

}
