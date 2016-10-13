package structmapper_test

import (
	"testing"

	"net"

	"github.com/anexia-it/go-structmapper"
	"github.com/stretchr/testify/require"
)

type mapperTestStructSimple struct {
	A string `mapper:"eff,omitempty"`
}

type mapperTestStructNested struct {
	// Even though a tag is set, this should be ignored
	privateTest string                  `mapper:"private"`
	A           string                  `mapper:"a"`
	B           int                     `mapper:"b"`
	C           float64                 `mapper:"c"`
	D           uint64                  `mapper:"dee,omitempty"`
	E           *mapperTestStructSimple `mapper:"e"`
}

type mapperTestStructArraySlice struct {
	A []string                  `mapper:"a"`
	B []*mapperTestStructSimple `mapper:"b,omitempty"`
	C [2]string                 `mapper:"c"`
}

type mapperTestStructTextMarshaler struct {
	IP net.IP
}

type mapperTestStructMap struct {
	A map[int]string  `mapper:"a"`
	B map[int]float32 `mapper:"bee"`
	C map[string]int  `mapper:"z"`
}

type mapperTestStructBool struct {
	A bool
}

func TestMapper_Roundtrip_Map(t *testing.T) {
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

	target := &mapperTestStructMap{}

	// Convert struct to map
	m, err := sm.ToMap(source)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)
}

func TestMapper_Roundtrip_Simple(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructSimple{
		A: "test value",
	}

	target := &mapperTestStructSimple{}

	// Convert struct to map
	m, err := sm.ToMap(source)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)
}

func TestMapper_Roundtrip_Nested(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructNested{
		A: "0",
		B: 1,
		C: 2.1,
		D: 3,
		E: &mapperTestStructSimple{
			A: "4",
		},
	}

	target := &mapperTestStructNested{}

	// Convert struct to map
	m, err := sm.ToMap(source)
	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)

	// Define second source
	source2 := &mapperTestStructArraySlice{
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

	target2 := &mapperTestStructArraySlice{}

	// Convert struct to map
	m, err = sm.ToMap(source2)
	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target2))

	// Check if source and target are equal
	require.EqualValues(t, source2, target2)

}

func TestMapper_Roundtrip_ArraySlice(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructArraySlice{
		A: []string{"test value", "test value 1"},
		B: []*mapperTestStructSimple{
			{
				A: "test0",
			},
			{
				A: "test1",
			},
		},
		C: [2]string{"a", "b"},
	}

	target := &mapperTestStructArraySlice{}

	// Convert struct to map
	m, err := sm.ToMap(source)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)
}

func TestMapper_Roundtrip_Bool(t *testing.T) {
	// Initialize Mapper without options
	sm, err := structmapper.NewMapper()
	require.NoError(t, err)
	require.NotNil(t, sm)

	source := &mapperTestStructBool{
		A: true,
	}

	target := &mapperTestStructBool{}

	// Convert struct to map
	m, err := sm.ToMap(source)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)

	source = &mapperTestStructBool{
		A: false,
	}

	target = &mapperTestStructBool{}

	// Convert struct to map
	m, err = sm.ToMap(source)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Convert back to struct
	require.NoError(t, sm.ToStruct(m, target))

	// Check if source and target are equal
	require.EqualValues(t, source, target)
}
