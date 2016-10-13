package structmapper

import (
	"testing"

	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

func TestOptionTagName(t *testing.T) {
	sm, err := NewMapper(OptionTagName("test"))

	require.NoError(t, err)
	require.NotNil(t, sm)

	// Check if the supplied tag name was set
	require.EqualValues(t, "test", sm.tagName)

	// Try setting an empty tag name
	sm, err = NewMapper(OptionTagName(""))
	require.Error(t, err)
	require.Nil(t, sm)

	multiErr, ok := err.(*multierror.Error)
	require.EqualValues(t, ok, true, "Returned error is not a multierror.Error")
	require.Len(t, multiErr.WrappedErrors(), 1)
	require.EqualError(t, multiErr.WrappedErrors()[0], ErrTagNameEmpty.Error())
}

func TestParseTag(t *testing.T) {
	// Check if the special-case ignore-me tag ("-") gives the correct result
	name, omitEmpty, err := parseTag("-")
	require.EqualValues(t, "-", name)
	require.EqualValues(t, false, omitEmpty)
	require.NoError(t, err)

	// Check if ",omitEmpty" alone works
	name, omitEmpty, err = parseTag(",omitempty")
	require.EqualValues(t, "", name)
	require.EqualValues(t, true, omitEmpty)
	require.NoError(t, err)

	// Check if "name,omitEmpty" returns the correct tag name
	name, omitEmpty, err = parseTag("test,omitempty")
	require.EqualValues(t, "test", name)
	require.EqualValues(t, true, omitEmpty)
	require.NoError(t, err)

	// Check if a punctation inside the tag name gives an error
	name, omitEmpty, err = parseTag("test.,omitempty")
	require.EqualValues(t, "test.", name)
	require.EqualValues(t, true, omitEmpty)
	require.Error(t, err)

	invalidTagErr, ok := err.(*InvalidTag)
	require.EqualValues(t, true, ok, "Not an InvalidTag error")

	require.EqualValues(t, "test.,omitempty", invalidTagErr.Tag())

	// Check if whitespace inside the tag name gives an error
	name, omitEmpty, err = parseTag("test ,omitempty")
	require.EqualValues(t, "test ", name)
	require.EqualValues(t, true, omitEmpty)
	require.Error(t, err)

	invalidTagErr, ok = err.(*InvalidTag)
	require.EqualValues(t, true, ok, "Not an InvalidTag error")

	require.EqualValues(t, "test ,omitempty", invalidTagErr.Tag())
}

func TestIsInvalidTag(t *testing.T) {
	// Test if IsInvalidTag works correctly for an invalid tag
	err := newErrorInvalidTag("test")
	require.NotNil(t, err)

	it, ok := IsInvalidTag(err)
	require.EqualValues(t, true, ok, "IsInvalidTag should return true")
	require.NotNil(t, it)
	require.EqualValues(t, err, it)

	// Test if IsInvalidTag works correctly for non-invalid tag error
	err = errors.New("test")
	it, ok = IsInvalidTag(err)
	require.EqualValues(t, false, ok, "IsInvalidTag should return false")
	require.Nil(t, it)
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
