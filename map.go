package structmapper

import (
	"encoding"
	"reflect"
	"unicode"

	"github.com/hashicorp/go-multierror"
)

// This file contains the struct to map functionality of Mapper

func (sm *Mapper) mapMap(v reflect.Value) (m map[interface{}]interface{}, err error) {
	keys := v.MapKeys()
	m = make(map[interface{}]interface{}, len(keys))

	for i := 0; i < len(keys); i++ {
		keyV := keys[i]
		keyI := keyV.Interface()
		valueV := v.MapIndex(keyV)

		valueI, mapErr := sm.mapValue(valueV.Interface(), valueV)

		if mapErr != nil {
			err = multierror.Append(err, mapErr)
			continue
		}
		m[keyI] = valueI
	}

	return
}

func (sm *Mapper) mapSlice(v reflect.Value) (s []interface{}, err error) {
	s = make([]interface{}, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		valueV := v.Index(i)
		valueI := valueV.Interface()

		mappedValueI, mapErr := sm.mapValue(valueI, valueV)
		if mapErr != nil {
			err = multierror.Append(err, mapErr)
			continue
		}
		s = append(s, mappedValueI)
	}

	return
}

func (sm *Mapper) mapValue(i interface{}, v reflect.Value) (value interface{}, err error) {
	// Check if the passed interface implements encoding.TextMarshaler, in which case we use the marshaler
	// for generating the value
	if marshaler, ok := i.(encoding.TextMarshaler); ok {
		text, marshalErr := marshaler.MarshalText()
		if marshalErr != nil {
			err = marshalErr
		}
		// As MarshalText is expected to return a textual representation, convert this back to a string
		value = string(text)
		return
	}

	// At this point it is safe to get rid of a possible pointer...
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	} else if v.Kind() == reflect.Ptr {
		// No-op for nil-pointers
		return
	}

	// Per-type handling
	switch v.Kind() {
	case reflect.Struct:
		// Handle struct
		value, err = sm.mapStruct(v)
	case reflect.Slice, reflect.Array:
		value, err = sm.mapSlice(v)
	case reflect.Map:
		value, err = sm.mapMap(v)
	default:
		// All other types are mapped as-is
		value = i
	}

	return
}

func (sm *Mapper) mapAnonymousField(m map[string]interface{}, v reflect.Value) error {
	if v.Kind() == reflect.Ptr && v.IsValid() && !v.IsNil() {
		v = v.Elem()
	}

	// Call mapStruct on anoynmous field
	mappedFields, err := sm.mapStruct(v)
	if err != nil {
		return err
	}

	// Merge onto map of struct containing anonymous field
	for key, value := range mappedFields {
		m[key] = value
	}
	return nil
}

func (sm *Mapper) mapStruct(v reflect.Value) (m map[string]interface{}, err error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}

	t := v.Type()

	// Create a new map that is pre-allocated with the number of fields v contains
	m = make(map[string]interface{}, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldD := t.Field(i)
		fieldV := v.Field(i)

		if fieldD.Anonymous {
			if anonErr := sm.mapAnonymousField(m, fieldV); anonErr != nil {
				err = multierror.Append(err, anonErr)
			}
			continue
		}

		fieldName := fieldD.Name

		if !unicode.IsUpper([]rune(fieldName)[0]) {
			// Ignore private fields
			continue
		}

		fieldName, omitEmpty, tagErr := parseTagFromStructField(fieldD, sm.tagName)
		if tagErr != nil {
			// Parsing the tag failed, ignore the field and carry on
			err = multierror.Append(err, tagErr)
			continue
		}

		fieldI := fieldV.Interface()

		if fieldName == "-" || omitEmpty && IsNilOrEmpty(fieldI, fieldV) {
			// Tag defines that field shall be ignored or omitEmpty is set
			// and the field is nil or empty
			continue
		} else if fieldI != nil {
			// If field is non-nil, map it...
			mappedFieldI, mappingErr := sm.mapValue(fieldI, fieldV)
			if mappingErr != nil {
				// If mapping failed, add an error
				err = multierror.Append(err, multierror.Prefix(mappingErr, fieldName+":"))
				continue
			}

			if omitEmpty && IsNilOrEmpty(mappedFieldI, reflect.ValueOf(mappedFieldI)) {
				// If omitEmpty is set and the mapped value is nil or zero carry on
				continue
			}
			// Override fieldI with the mapped value
			fieldI = mappedFieldI
		}

		m[fieldName] = fieldI
	}

	return
}

func (sm *Mapper) toMap(s interface{}) (map[string]interface{}, error) {
	if s == nil {
		// If the input struct is nil, return an empty map
		return map[string]interface{}{}, nil
	}

	// Verify that we are working on a struct...
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, ErrNotAStruct
	}

	return sm.mapStruct(v)
}
