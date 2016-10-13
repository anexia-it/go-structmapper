package structmapper

import "github.com/hashicorp/go-multierror"

// Mapper provides the mapping logic
type Mapper struct {
	tagName string
}

// ToStruct takes a source map[string]interface{} and maps its values onto a target struct.
func (mapper *Mapper) ToStruct(source map[string]interface{}, target interface{}) error {
	return mapper.toStruct(source, target)
}

// ToMap takes a source struct and maps its values onto a map[string]interface{}, which is then returned.
func (mapper *Mapper) ToMap(source interface{}) (map[string]interface{}, error) {
	return mapper.toMap(source)
}

// NewMapper initializes a new mapper instance.
// Optionally Mapper options may be passed to this function
func NewMapper(options ...Option) (*Mapper, error) {
	sm := &Mapper{}

	var err error

	// Apply default options first
	for _, opt := range defaultOptions {
		if optErr := opt(sm); optErr != nil {
			// Panic if default option could not be applied
			panic(optErr)
		}
	}

	// ... and passed options afterwards.
	// This way the passed options override the default options
	for _, opt := range options {
		if optErr := opt(sm); optErr != nil {
			err = multierror.Append(err, optErr)
		}
	}

	if err != nil {
		return nil, err
	}

	return sm, nil
}
