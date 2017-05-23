package structmapper

import (
	"fmt"
	"reflect"
)

// This file contains utility functions

// IsNilOrEmpty checks if a passed interface is either nil or of
// the type's zero value.
//
// The zero value depends on the passed type. For example, the zero value of a string is an empty string.
func IsNilOrEmpty(i interface{}, v reflect.Value) bool {
	// Simple case: interface is nil
	if i == nil {
		return true
	}

	// Hard case: check if interface has "zero" value (ie. empty string, zero integer, etc.)
	return reflect.DeepEqual(i, reflect.Zero(v.Type()).Interface())
}

// ForceStringMapKeys takes a map[string]interface{} and ensures that all maps which are nested
// in this top-level map are of type map[string]interface{}.
// The map that is passed in is expected to contain only primitive types, maps, slices and arrays.
//
// The purpose of this function is preparing a map returned by ToMap() to something encodings
// like the standard library's JSON encoding can work with.
//
// Keys which are not strings already are converted to strings by either using the key's String() method,
// if available, converting via reflect conversion or falling back to using fmt.Sprint() for conversion.
func ForceStringMapKeys(in map[string]interface{}) (out map[string]interface{}, err error) {
	out = make(map[string]interface{}, len(in))
	for key, value := range in {
		if out[key], err = convertToStringKeys(value); err != nil {
			out = nil
			return
		}
	}

	return
}

func convertToStringKeys(in interface{}) (out interface{}, err error) {
	return convertValueToStringKeys(reflect.ValueOf(in))
}

func convertValueToStringKeys(in reflect.Value) (out interface{}, err error) {
	if !in.IsValid() || in.Interface() == nil {
		return nil, nil
	}

	inKind := in.Kind()
	switch inKind {
	case reflect.Map:
		out, err = convertMapToStringKeys(in)
		return
	case reflect.Slice:
		out, err = convertSliceToStringKeys(in)
		return
	case reflect.Array:
		out, err = convertArrayToStringKeys(in)
		return
	case reflect.Interface:
		out, err = convertValueToStringKeys(in.Elem())
		return
	case reflect.Struct, reflect.Chan, reflect.Func:
		err = fmt.Errorf("Conversion of type %s (kind: %s) is not supported", in.Type().String(),
			inKind.String())
	case reflect.Ptr:
		if !in.IsNil() {
			out, err = convertValueToStringKeys(in.Elem())
			return
		}

		// If we received a nil-pointer we fall through to the default case
	}

	// Default: Any case not handled above is a primitive type which does not require conversion
	out = in.Interface()
	return
}

func convertSliceToStringKeys(in reflect.Value) (out interface{}, err error) {
	// No-op for empty slices
	if in.Len() == 0 {
		return in.Interface(), nil
	}

	outSlice := reflect.MakeSlice(in.Type(), in.Len(), in.Cap())

	for i := 0; i < in.Len(); i++ {
		var val interface{}
		if val, err = convertValueToStringKeys(in.Index(i)); err != nil {
			return
		}

		outSlice.Index(i).Set(reflect.ValueOf(val))
	}

	out = outSlice.Interface()
	return
}

func convertArrayToStringKeys(in reflect.Value) (out interface{}, err error) {
	// No-op for empty arrays
	if in.Len() == 0 {
		return in.Interface(), nil
	}

	outArray := reflect.New(in.Type()).Elem()

	for i := 0; i < in.Len(); i++ {
		var val interface{}
		if val, err = convertValueToStringKeys(in.Index(i)); err != nil {
			return
		}

		outArray.Index(i).Set(reflect.ValueOf(val))
	}

	out = outArray.Interface()
	return
}

func convertMapToStringKeys(in reflect.Value) (out map[string]interface{}, err error) {
	stringType := reflect.TypeOf("")
	inKeys := in.MapKeys()
	out = make(map[string]interface{}, len(inKeys))
	for _, key := range inKeys {
		var keyString string
		wasConverted := false
		keyInterface := key.Interface()

		if stringer, ok := keyInterface.(fmt.Stringer); ok {
			// Key implements fmt.Stringer: use value returned by String()
			keyString = stringer.String()
			wasConverted = true
		} else if goStringer, ok := keyInterface.(fmt.GoStringer); ok {
			// Key implements fmt.GoStringer: use value returned by GoString()
			keyString = goStringer.GoString()
			wasConverted = true
		} else if key.Kind() == reflect.String {
			// Key is already a string or a type based on string: use key.String() to obtain the string
			// value
			keyString = key.String()
			wasConverted = true
		} else if key.Kind() == reflect.Interface && key.IsValid() &&
			key.Elem().Type().ConvertibleTo(stringType) {
			// Key is an interface, but has a type that is convertible to string underneath
			switch key.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// No-op, as conversion from integer-types to string is possible using reflect,
				// but gives us the unicode character corresponding to the integer's value
			default:
				keyString = key.Elem().Convert(stringType).Interface().(string)
				wasConverted = true
			}
		}

		if !wasConverted {
			// Last resort: use fmt.Sprint to obtain a value
			keyString = fmt.Sprint(keyInterface)
		}

		if out[keyString], err = convertValueToStringKeys(in.MapIndex(key)); err != nil {
			out = nil
			return
		}
	}

	return
}
