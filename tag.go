package structmapper

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// DefaultTagName defines the default tag name used by Mapper
const DefaultTagName = "mapper"

// Option defines the type used by Mapper Option functions
type Option func(*Mapper) error

// OptionTagName sets the tag name the mapper uses
func OptionTagName(tagName string) Option {
	return func(m *Mapper) error {
		if tagName == "" {
			return ErrTagNameEmpty
		}

		m.tagName = tagName
		return nil
	}
}

// Default options for Mapper
var defaultOptions = []Option{
	OptionTagName(DefaultTagName),
}

var _ error = (*InvalidTag)(nil)

// InvalidTag is an error that indicates that the tag value was invalid
type InvalidTag struct {
	tag string
}

// Error returns the error string and causes InvalidTag to implement the error interface
func (it *InvalidTag) Error() string {
	return fmt.Sprintf("Invalid tag: '%s'", it.tag)
}

// Tag returns the tag name
func (it *InvalidTag) Tag() string {
	return it.tag
}

func newErrorInvalidTag(tag string) error {
	return &InvalidTag{
		tag: tag,
	}
}

// IsInvalidTag checks if the given error is an InvalidTag error
// and returns the InvalidTag error along with a boolean that defines
// if it is indeed an invalid tag error.
// The returned *InvalidTag may be nil, if the flag is false
func IsInvalidTag(err error) (*InvalidTag, bool) {
	it, ok := err.(*InvalidTag)
	return it, ok
}

// parseTagFromStructField is a helper that calls parseTag given a reflect.StructField and a tag name
func parseTagFromStructField(f reflect.StructField, tagName string) (name string, omitEmpty bool, err error) {
	tag := f.Tag.Get(tagName)
	name, omitEmpty, err = parseTag(tag)
	if name == "" {
		name = f.Name
	}
	return
}

// parseTag parses a tag string and returns the corresponding name, omitEmpty flag and a possible
// error
func parseTag(tag string) (name string, omitEmpty bool, err error) {
	name = tag

	// Handle the "ignore me" tag value
	if name == "-" {
		return
	}

	// Check if the tag has the omitempty suffix set
	if strings.HasSuffix(tag, ",omitempty") {
		// Update the omitEmpty flag and strip the suffix from the tag name
		omitEmpty = true
		name = strings.TrimSuffix(tag, ",omitempty")
	}

	// Check if the rest of the tag does not contain any symbols
	for _, letter := range name {
		if letter != '_' && !unicode.IsLetter(letter) && !unicode.IsDigit(letter) {
			err = newErrorInvalidTag(tag)
			return
		}
	}

	return
}
