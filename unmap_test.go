package structmapper_test

import (
	"net"
	"testing"

	"github.com/anexia-it/go-structmapper"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mapperTestStructInterfaceField struct {
	A interface{} `mapper:"x"`
}

func TestMapper_ToStruct(t *testing.T) {
	t.Run("Errors", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		// Call ToStruct with nil map
		require.EqualError(t, sm.ToStruct(nil, &mapperTestStructSimple{}), structmapper.ErrMapIsNil.Error())

		testValue := "test"

		// Call ToStruct with non-struct pointer
		require.EqualError(t, sm.ToStruct(make(map[string]interface{}), &testValue), structmapper.ErrNotAStruct.Error())

		// Call ToStruct with non-struct pointer
		require.EqualError(t, sm.ToStruct(make(map[string]interface{}), mapperTestStructSimple{}),
			structmapper.ErrNotAStructPointer.Error())
	})

	t.Run("InterfaceField", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		m := map[string]interface{}{
			"x": "test",
		}

		target := &mapperTestStructInterfaceField{}

		err = sm.ToStruct(m, target)
		require.Error(t, err)
		me, ok := err.(*multierror.Error)
		require.EqualValues(t, true, ok, "Returned error is not a *multierror.Error")
		require.Len(t, me.Errors, 1)
		e := me.Errors[0]
		// Test if the error is correct...
		require.Error(t, e, multierror.Prefix(structmapper.ErrFieldIsInterface, "x: ").Error())
	})

	t.Run("Simple", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		// Simple test case: single field, no nesting
		expected := &mapperTestStructSimple{
			A: "test",
		}

		m := map[string]interface{}{
			"eff": "test",
		}

		target := &mapperTestStructSimple{}

		require.NoError(t, sm.ToStruct(m, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("NestedSimple", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		// Construct simple test case: all fields present
		expected := &mapperTestStructNested{
			A: "0",
			B: 1,
			C: 2.1,
			D: 3,
			E: &mapperTestStructSimple{
				A: "4",
			},
		}

		m := map[string]interface{}{
			"a":   "0",
			"b":   1,
			"c":   2.1,
			"dee": uint64(3),
			"e": map[string]interface{}{
				"eff": "4",
			},
		}

		target := &mapperTestStructNested{}

		require.NoError(t, sm.ToStruct(m, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("ArraySlice", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		expected := &mapperTestStructArraySlice{
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

		m := map[string]interface{}{
			"a": []interface{}{"0.0", "0.1"},
			"b": []interface{}{
				map[string]interface{}{
					"eff": "1.0",
				},
				map[string]interface{}{},
			},
			"c": []interface{}{"2.0", ""},
		}

		target := &mapperTestStructArraySlice{}

		require.NoError(t, sm.ToStruct(m, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("TextUnmarshaler", func(t *testing.T) {
		sm, err := structmapper.NewMapper()

		require.NoError(t, err)
		require.NotNil(t, sm)

		ip := net.ParseIP("127.0.0.1")
		require.NotNil(t, ip)

		expected := &mapperTestStructTextMarshaler{
			IP: ip,
		}

		m := map[string]interface{}{
			"IP": ip.String(),
		}

		target := &mapperTestStructTextMarshaler{}

		require.NoError(t, sm.ToStruct(m, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("Map", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		expected := &mapperTestStructMap{
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

		m := map[string]interface{}{
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

		target := &mapperTestStructMap{}

		// Convert map to struct
		require.NoError(t, sm.ToStruct(m, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("Anonymous", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		expected := &MapperTestStructAnonymousOuter{
			MapperTestStructAnonymousInner: MapperTestStructAnonymousInner{
				A: "inner",
			},
			A: "outer",
		}

		target := &MapperTestStructAnonymousOuter{}

		source := map[string]interface{}{
			"a_inner": "inner",
			"a_outer": "outer",
		}

		require.NoError(t, sm.ToStruct(source, target))

		require.EqualValues(t, expected, target)
	})

	t.Run("AnonymousPtr", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		expected := &MapperTestStructAnonymousPtrOuter{
			MapperTestStructAnonymousInner: &MapperTestStructAnonymousInner{
				A: "inner",
			},
			A: "outer",
		}

		target := &MapperTestStructAnonymousPtrOuter{}

		source := map[string]interface{}{
			"a_inner": "inner",
			"a_outer": "outer",
		}

		require.NoError(t, sm.ToStruct(source, target))

		require.EqualValues(t, expected, target)
	})

	t.Run("TypeMismatch", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		target := &mapperTestStructSimple{}

		source := map[string]interface{}{
			"eff": 3.1415,
		}

		err = sm.ToStruct(source, target) //
		require.Error(t, err)
		w, ok := err.(errwrap.Wrapper)
		require.EqualValues(t, true, ok, "returned error is not an errwrap.Wrapper")
		wrapped := w.WrappedErrors()
		require.Len(t, wrapped, 1)
		require.EqualError(t, wrapped[0], "eff: Type mismatch: string and float64 are incompatible")
	})

	t.Run("MapInterfaceInterface", func(t *testing.T) {
		// Initialize Mapper without options
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		target := &mapperTestStructNested{}
		expected := &mapperTestStructNested{
			E: &mapperTestStructSimple{
				A: "test",
			},
		}
		source := map[string]interface{}{
			"e": map[interface{}]interface{}{
				"eff": "test",
			},
		}

		require.NoError(t, sm.ToStruct(source, target))
		require.EqualValues(t, expected, target)
	})

	t.Run("MapValueTypeCompability", func(t *testing.T) {
		sm, err := structmapper.NewMapper()
		require.NoError(t, err)
		require.NotNil(t, sm)

		target := mapperTestStructMapInt16{}
		source := map[string]interface{}{
			"data": map[string]interface{}{
				"0": float32(1),
				"1": float32(2),
			},
		}
		expected := mapperTestStructMapInt16{
			Data: map[string]int16{
				"0": 1,
				"1": 2,
			},
		}

		assert.NoError(t, sm.ToStruct(source, &target))
		assert.EqualValues(t, expected, target)
	})
}
