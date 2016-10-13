package structmapper_test

import (
	"testing"

	"net"

	"github.com/anexia-it/go-structmapper"
	"github.com/stretchr/testify/require"
)

func TestMapper_ToMap_Errors(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Call ToMap with nil value
	m, err := sm.ToMap(nil)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Len(t, m, 0)

	// Call ToMap with non-struct
	m, err = sm.ToMap("test")
	require.EqualError(t, err, structmapper.ErrNotAStruct.Error())
	require.Nil(t, m)

	// Call ToMap with pointer to non-struct
	testValue := "test"
	m, err = sm.ToMap(&testValue)
	require.EqualError(t, err, structmapper.ErrNotAStruct.Error())
	require.Nil(t, m)
}

func TestMapper_ToMap_Simple(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Test if mapping a simple structure works
	source := &mapperTestStructSimple{
		A: "test value",
	}

	expected := map[string]interface{}{
		"eff": "test value",
	}

	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)
}

func TestMapper_ToMap_NestedSimple(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Construct simple test case: all fields present
	source := &mapperTestStructNested{
		A: "0",
		B: 1,
		C: 2.1,
		D: 3,
		E: &mapperTestStructSimple{
			A: "4",
		},
	}

	expected := map[string]interface{}{
		"a":   "0",
		"b":   1,
		"c":   2.1,
		"dee": uint64(3),
		"e": map[string]interface{}{
			"eff": "4",
		},
	}

	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)

	// Test if omission of fields works
	source = &mapperTestStructNested{
		A: "0",
		B: 1,
		C: 2.1,
		E: &mapperTestStructSimple{},
	}

	expected = map[string]interface{}{
		"a": "0",
		"b": 1,
		"c": 2.1,
		"e": map[string]interface{}{},
	}

	m, err = sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)
}

func TestMapper_ToMap_ArraySlice(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	testStructArraySlice := &mapperTestStructArraySlice{
		A: []string{"0.0", "0.1"},
		B: []*mapperTestStructSimple{
			{
				A: "1.0",
			},
			{
				A: "",
			},
		},
		C: [2]string{"2.0", ""},
	}

	expectedMap := map[string]interface{}{
		"a": []interface{}{"0.0", "0.1"},
		"b": []interface{}{
			map[string]interface{}{
				"eff": "1.0",
			},
			map[string]interface{}{},
		},
		"c": []interface{}{"2.0", ""},
	}

	m, err := sm.ToMap(testStructArraySlice)
	require.NoError(t, err)
	require.EqualValues(t, expectedMap, m)
}

func TestMapper_ToMap_TextMarshaler(t *testing.T) {
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	ip := net.ParseIP("127.0.0.1")

	source := &mapperTestStructTextMarshaler{
		IP: ip,
	}

	expected := map[string]interface{}{
		"IP": ip.String(),
	}

	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)
}

func TestMapper_ToMap_Map(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructMap{
		A: map[int]string{
			10:   "a",
			1024: "b",
			30:   "c",
		},
		B: map[int]float32{
			1: 1.1,
			2: 2.2,
		},
		C: map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}

	expected := map[string]interface{}{
		"a": map[interface{}]interface{}{
			10:   "a",
			1024: "b",
			30:   "c",
		},
		"bee": map[interface{}]interface{}{
			1: float32(1.1),
			2: float32(2.2),
		},
		"z": map[interface{}]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}

	// Convert struct to map
	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)
}

func TestMapper_ToMap_Anonymous(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructAnonymousOuter{
		mapperTestStructAnonymousInner: mapperTestStructAnonymousInner{
			A: "inner",
		},
		A: "outer",
	}

	expected := map[string]interface{}{
		"a_inner": "inner",
		"a_outer": "outer",
	}

	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.EqualValues(t, expected, m)
}
