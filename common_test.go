package structmapper_test

import (
	"testing"

	"net"

	structmapper "github.com/anexia-it/go-structmapper"
	"github.com/stretchr/testify/require"
)

type mapperTestStructSimple struct {
	A string `mapper:"eff,omitempty"`
}

type mapperTestStructNested struct {
	// Even though a tag is set, this should be ignored
	privateTest string                  `mapper:"private"` //nolint:structcheck,unused
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

type MapperTestStructAnonymousInner struct {
	A string `mapper:"a_inner"`
}

type MapperTestStructAnonymousOuter struct {
	MapperTestStructAnonymousInner
	A string `mapper:"a_outer"`
}

type MapperTestStructAnonymousPtrOuter struct {
	*MapperTestStructAnonymousInner
	A string `mapper:"a_outer"`
}

type mapperTestStructMapStringString struct {
	A string
	B map[string]string
}

type mapperTestStructNestedStructSlice struct {
	A string
	B []mapperTestStructSimple
}

type mapperTestStructNestedNestedStructSlice struct {
	A string
	B mapperTestStructNestedStructSlice
}

type mapperTestStructConvertibleSource struct {
	A int
}

type mapperTestStructConvertibleDest struct {
	A uint64
}

type mapperTestStructMapInt16 struct {
	Data map[string]int16 `mapper:"data"`
}

func TestMapper_Roundtrip(t *testing.T) {
	// Test round-trips between map/unmap

	t.Run("Map", func(t *testing.T) {
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
	})

	t.Run("Simple", func(t *testing.T) {
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
	})

	t.Run("Nested", func(t *testing.T) {
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

	})

	t.Run("ArraySlice", func(t *testing.T) {
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
	})

	t.Run("Bool", func(t *testing.T) {
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
	})

	t.Run("Anonymous", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		source := &MapperTestStructAnonymousOuter{
			MapperTestStructAnonymousInner: MapperTestStructAnonymousInner{
				A: "inner",
			},
			A: "outer",
		}

		target := &MapperTestStructAnonymousOuter{}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))

		require.EqualValues(t, source, target)
	})

	t.Run("MapStringString", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		source := &mapperTestStructMapStringString{
			A: "test0",
			B: map[string]string{
				"b0": "1",
				"b1": "2",
			},
		}

		target := &mapperTestStructMapStringString{}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))

		require.EqualValues(t, source, target)
	})

	t.Run("NestedStructSlice", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		source := &mapperTestStructNestedStructSlice{
			A: "test0",
			B: []mapperTestStructSimple{
				{
					A: "test1",
				},
				{
					A: "test2",
				},
			},
		}

		target := &mapperTestStructNestedStructSlice{}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))

		require.EqualValues(t, source, target)
	})

	t.Run("NestedNestedStructSlice", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		source := &mapperTestStructNestedNestedStructSlice{
			A: "test0",
			B: mapperTestStructNestedStructSlice{
				A: "test1",
				B: []mapperTestStructSimple{
					{
						A: "test1",
					},
					{
						A: "test2",
					},
				},
			},
		}

		target := &mapperTestStructNestedNestedStructSlice{}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))

		require.EqualValues(t, source, target)
	})

	t.Run("Convertible", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		source := &mapperTestStructConvertibleSource{
			A: 123456,
		}

		target := &mapperTestStructConvertibleDest{}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))

		require.EqualValues(t, uint64(source.A), target.A)
	})

	t.Run("MapInterfaceInterface", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		target := &mapperTestStructNested{}
		source := &mapperTestStructNested{
			E: &mapperTestStructSimple{
				A: "test",
			},
		}

		m, err := sm.ToMap(source)
		require.NoError(t, err)
		require.NotNil(t, m)

		require.NoError(t, sm.ToStruct(m, target))
	})
}
