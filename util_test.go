package structmapper_test

import (
	"reflect"
	"testing"

	"time"

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

func TestForceStringMapKeys(t *testing.T) {
	// Simple map test: should be returned as-is
	simpleMap := map[string]interface{}{
		"0":  "0",
		"1":  int(1),
		"2":  uint(2),
		"3":  int8(3),
		"4":  int16(4),
		"5":  int32(5),
		"6":  int64(6),
		"7":  uint8(7),
		"8":  uint16(8),
		"9":  uint32(9),
		"10": uint64(10),
		"11": float32(11.1),
		"12": float64(12.2),
		"13": []byte{1, 3},
		"14": complex(1.0, 4.0),
		"15": true,
		"16": nil,
		"17": (*string)(nil),
	}

	res, err := structmapper.ForceStringMapKeys(simpleMap)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(simpleMap))
	require.EqualValues(t, simpleMap, res)

	// Simple map test: zero values
	simpleMap = map[string]interface{}{
		"0":  "",
		"1":  int(0),
		"2":  uint(0),
		"3":  int8(0),
		"4":  int16(0),
		"5":  int32(0),
		"6":  int64(0),
		"7":  uint8(0),
		"8":  uint16(0),
		"9":  uint32(0),
		"10": uint64(0),
		"11": float32(0.0),
		"12": float64(0.0),
		"13": []byte{},
		"14": complex(0.0, 0.0),
		"15": false,
		"16": nil,
		"17": (*string)(nil),
	}
	res, err = structmapper.ForceStringMapKeys(simpleMap)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(simpleMap))
	require.EqualValues(t, simpleMap, res)

	// Test nesting with nested map[string]interface{}
	nestedMapStringInterface := map[string]interface{}{
		"0": map[string]interface{}{
			"0":  "0",
			"1":  int(1),
			"2":  uint(2),
			"3":  int8(3),
			"4":  int16(4),
			"5":  int32(5),
			"6":  int64(6),
			"7":  uint8(7),
			"8":  uint16(8),
			"9":  uint32(9),
			"10": uint64(10),
			"11": float32(11.1),
			"12": float64(12.2),
			"13": []byte{1, 3},
			"14": complex(1.0, 4.0),
			"15": true,
			"16": nil,
			"17": (*string)(nil),
		},
		"1": map[string]interface{}{
			"0":  "0",
			"1":  int(1),
			"2":  uint(2),
			"3":  int8(3),
			"4":  int16(4),
			"5":  int32(5),
			"6":  int64(6),
			"7":  uint8(7),
			"8":  uint16(8),
			"9":  uint32(9),
			"10": uint64(10),
			"11": float32(11.1),
			"12": float64(12.2),
			"13": []byte{1, 3},
			"14": complex(1.0, 4.0),
			"15": true,
			"16": nil,
			"17": (*string)(nil),
		},
	}

	res, err = structmapper.ForceStringMapKeys(nestedMapStringInterface)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(nestedMapStringInterface))
	require.EqualValues(t, nestedMapStringInterface, res)

	// Test nesting with slices and arrays
	nestedSlicesArrays := map[string]interface{}{
		"0": []string{"0.0", "0.1", "0.2"},
		"1": []int{0, 1, 2},
		"2": [3]uint{0, 1, 2},
		"4": []interface{}{
			map[string]interface{}{
				"0":  "0",
				"1":  int(1),
				"2":  uint(2),
				"3":  int8(3),
				"4":  int16(4),
				"5":  int32(5),
				"6":  int64(6),
				"7":  uint8(7),
				"8":  uint16(8),
				"9":  uint32(9),
				"10": uint64(10),
				"11": float32(11.1),
				"12": float64(12.2),
				"13": []byte{1, 3},
				"14": complex(1.0, 4.0),
				"15": true,
				"16": nil,
				"17": (*string)(nil),
			},
		},
		"5": [2]interface{}{
			map[string]interface{}{
				"0":  "0",
				"1":  int(1),
				"2":  uint(2),
				"3":  int8(3),
				"4":  int16(4),
				"5":  int32(5),
				"6":  int64(6),
				"7":  uint8(7),
				"8":  uint16(8),
				"9":  uint32(9),
				"10": uint64(10),
				"11": float32(11.1),
				"12": float64(12.2),
				"13": []byte{1, 3},
				"14": complex(1.0, 4.0),
				"15": true,
				"16": nil,
				"17": (*string)(nil),
			},
			map[string]interface{}{
				"0":  "0",
				"1":  int(1),
				"2":  uint(2),
				"3":  int8(3),
				"4":  int16(4),
				"5":  int32(5),
				"6":  int64(6),
				"7":  uint8(7),
				"8":  uint16(8),
				"9":  uint32(9),
				"10": uint64(10),
				"11": float32(11.1),
				"12": float64(12.2),
				"13": []byte{1, 3},
				"14": complex(1.0, 4.0),
				"15": true,
				"16": nil,
				"17": (*string)(nil),
			},
		},
	}

	res, err = structmapper.ForceStringMapKeys(nestedSlicesArrays)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(nestedSlicesArrays))
	require.EqualValues(t, nestedSlicesArrays, res)

	// Test nesting with empty arrays and slices
	nestedEmptyArraysSlices := map[string]interface{}{
		"0": [0]string{},
		"1": []int{},
		"2": "test",
	}

	res, err = structmapper.ForceStringMapKeys(nestedEmptyArraysSlices)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(nestedEmptyArraysSlices))
	require.EqualValues(t, nestedEmptyArraysSlices, res)

	// Test key conversion: Stringer
	now := time.Now()
	duration := time.Millisecond * 10
	keyConversionStringer := map[string]interface{}{
		"0": map[interface{}]interface{}{
			now:      0,
			duration: 1,
		},
		"1": map[interface{}]interface{}{
			now:      2,
			duration: 3,
		},
	}
	res, err = structmapper.ForceStringMapKeys(keyConversionStringer)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(keyConversionStringer))
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			now.String():      0,
			duration.String(): 1,
		},
		"1": map[string]interface{}{
			now.String():      2,
			duration.String(): 3,
		},
	}, res)

	// Test key conversion: GoStringer
	keyConversionGoStringer := map[string]interface{}{
		"0": map[interface{}]interface{}{
			goStringerString("0"): 0,
			goStringerString("1"): 1,
		},
		"1": map[interface{}]interface{}{
			goStringerString("2"): 2,
			goStringerString("3"): 3,
		},
	}
	res, err = structmapper.ForceStringMapKeys(keyConversionGoStringer)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(keyConversionGoStringer))
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			"gostringer-0": 0,
			"gostringer-1": 1,
		},
		"1": map[string]interface{}{
			"gostringer-2": 2,
			"gostringer-3": 3,
		},
	}, res)

	// Test key conversion: convertible to string
	keyConversionConvertible := map[string]interface{}{
		"0": map[interface{}]interface{}{
			convertibleToString("00"): 0,
			convertibleToString("01"): 1,
		},
		"1": map[interface{}]interface{}{
			convertibleToString("10"): 2,
			convertibleToString("11"): 3,
		},
	}
	res, err = structmapper.ForceStringMapKeys(keyConversionConvertible)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(keyConversionConvertible))
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			"00": 0,
			"01": 1,
		},
		"1": map[string]interface{}{
			"10": 2,
			"11": 3,
		},
	}, res)

	// Test key conversion: convertible to string, map has key type based on string
	keyConversionConvertible2 := map[string]interface{}{
		"0": map[convertibleToString]interface{}{
			convertibleToString("00"): 0,
			convertibleToString("01"): 1,
		},
		"1": map[convertibleToString]interface{}{
			convertibleToString("10"): 2,
			convertibleToString("11"): 3,
		},
	}

	res, err = structmapper.ForceStringMapKeys(keyConversionConvertible2)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(keyConversionConvertible2))
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			"00": 0,
			"01": 1,
		},
		"1": map[string]interface{}{
			"10": 2,
			"11": 3,
		},
	}, res)

	// Test key conversion: last resort using fmt.Sprint
	keyConversionLastResort := map[string]interface{}{
		"0": map[interface{}]interface{}{
			0.5:     0,
			uint(1): 1,
			true:    2,
		},
		"1": map[interface{}]interface{}{
			0.5:       0,
			uint64(1): 1,
			true:      2,
		},
	}

	res, err = structmapper.ForceStringMapKeys(keyConversionLastResort)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, len(keyConversionLastResort))
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			"0.5":  0,
			"1":    1,
			"true": 2,
		},
		"1": map[string]interface{}{
			"0.5":  0,
			"1":    1,
			"true": 2,
		},
	}, res)

	// Test: struct value should give us an error
	structValueMap := map[string]interface{}{
		"0": testStruct{},
	}

	res, err = structmapper.ForceStringMapKeys(structValueMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type structmapper_test.testStruct (kind: struct) is not supported")

	// Test: chan value should give us an error
	ch := make(chan bool)
	defer close(ch)
	chanValueMap := map[string]interface{}{
		"0": ch,
	}

	res, err = structmapper.ForceStringMapKeys(chanValueMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type chan bool (kind: chan) is not supported")

	// Test: func value should give us an error
	funcValueMap := map[string]interface{}{
		"0": testFn,
	}

	res, err = structmapper.ForceStringMapKeys(funcValueMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type func() (kind: func) is not supported")

	// Test: pointer conversion
	testString := "test"
	pointerConversionMap := map[string]interface{}{
		"0": map[interface{}]interface{}{
			"0": &testString,
		},
	}

	res, err = structmapper.ForceStringMapKeys(pointerConversionMap)
	require.NoError(t, err)
	require.EqualValues(t, map[string]interface{}{
		"0": map[string]interface{}{
			"0": testString,
		},
	}, res)

	// Test: slice element value conversion error
	sliceElementStructMap := map[string]interface{}{
		"0": []interface{}{
			testStruct{},
		},
	}
	res, err = structmapper.ForceStringMapKeys(sliceElementStructMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type structmapper_test.testStruct (kind: struct) is not supported")

	// Test: array element value conversion error
	arrayElementStructMap := map[string]interface{}{
		"0": [1]interface{}{
			testStruct{},
		},
	}
	res, err = structmapper.ForceStringMapKeys(arrayElementStructMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type structmapper_test.testStruct (kind: struct) is not supported")

	// Test: map element value conversion error
	// Additionally use a pointer to testStruct to also have pointer conversion in this test
	mapElementStructMap := map[string]interface{}{
		"0": map[string]interface{}{
			"0": &testStruct{},
		},
	}
	res, err = structmapper.ForceStringMapKeys(mapElementStructMap)
	require.Nil(t, res)
	require.EqualError(t, err, "Conversion of type structmapper_test.testStruct (kind: struct) is not supported")
}

type testStruct struct {
}

func testFn() {

}

type goStringerString string

func (s goStringerString) GoString() string {
	return "gostringer-" + string(s)
}

type convertibleToString string
