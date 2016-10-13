package structmapper_test

import (
	"reflect"
	"testing"

	"github.com/anexia-it/go-structmapper"
	"github.com/stretchr/testify/require"
)

func TestIsNilOrEmpty(t *testing.T) {
	zeroString := ""
	nonZeroString := "test"

	// Check if nil value returns true
	require.EqualValues(t, true, structmapper.IsNilOrEmpty(nil, reflect.ValueOf(zeroString)))

	// Check if empty string returns true
	require.EqualValues(t, true, structmapper.IsNilOrEmpty(zeroString, reflect.ValueOf(zeroString)))
	// Check if non-empty string returns false
	require.EqualValues(t, false, structmapper.IsNilOrEmpty(nonZeroString, reflect.ValueOf(zeroString)))

	// Check if a pointer to an empty string returns false (because the pointer is non-nil)
	require.EqualValues(t, false, structmapper.IsNilOrEmpty(&zeroString, reflect.ValueOf(&zeroString)))
	// Check if a pointer to an empty string returns false (because the pointer is non-nil)
	require.EqualValues(t, false, structmapper.IsNilOrEmpty(&nonZeroString, reflect.ValueOf(&zeroString)))

	// TODO: add additional test cases for types other than string
}
