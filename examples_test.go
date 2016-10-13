package structmapper_test

import (
	"fmt"
	"net"
	"strings"

	"github.com/anexia-it/go-structmapper"
)

func ExampleMapper_ToMap() {
	type ExampleStruct struct {
		A       string
		B       int `mapper:"bee"`
		C       []string
		D       map[string]int
		E       net.IP
		Ignored string `mapper:"-"`
	}

	// Create new structmapper instance with default configuration
	sm, err := structmapper.NewMapper()
	if err != nil {
		panic(err)
	}

	source := &ExampleStruct{
		A: "a",
		B: 5,
		C: []string{"c0", "c1"},
		D: map[string]int{
			"d0": 1,
			"d1": 2,
		},
		E:       net.ParseIP("127.0.0.1"),
		Ignored: "should be ignored",
	}

	m, err := sm.ToMap(source)
	if err != nil {
		panic(err)
	}

	c := []string{}

	for _, v := range m["C"].([]interface{}) {
		c = append(c, v.(string))
	}
	d := m["D"].(map[interface{}]interface{})

	fmt.Printf("A=%s,B=%d,C=%s,D=d0=%d,d1=%d,E=%s,Ignored=%s", m["A"], m["bee"], strings.Join(c, ","),
		d["d0"], d["d1"], m["E"], m["Ignored"])
	// Output: A=a,B=5,C=c0,c1,D=d0=1,d1=2,E=127.0.0.1,Ignored=%!s(<nil>)

}

func ExampleMapper_ToStruct() {
	type ExampleStruct struct {
		A       string
		B       int `mapper:"bee"`
		C       []string
		D       map[string]int
		E       net.IP
		Ignored string `mapper:"-"`
	}

	// Create new structmapper instance with default configuration
	sm, err := structmapper.NewMapper()
	if err != nil {
		panic(err)
	}

	source := map[string]interface{}{
		"A":   "a",
		"bee": 5,
		"C":   []string{"c0", "c1"},
		"D": map[interface{}]interface{}{
			"d0": 1,
			"d1": 2,
		},
		"E": "127.0.0.1",
	}

	target := &ExampleStruct{}

	if err := sm.ToStruct(source, target); err != nil {
		panic(err)
	}

	fmt.Printf("A=%s,B=%d,C=%s,D=d0=%d,d1=%d,E=%s,Ignored=%s", target.A, target.B, strings.Join(target.C, ","),
		target.D["d0"], target.D["d1"], target.E, target.Ignored)
	// Output: A=a,B=5,C=c0,c1,D=d0=1,d1=2,E=127.0.0.1,Ignored=
}

func ExampleMapper_roundtrip() {
	type ExampleStruct struct {
		A       string
		B       int `mapper:"bee"`
		C       []string
		D       map[string]int
		E       net.IP
		Ignored string `mapper:"-"`
	}

	// Create new structmapper instance with default configuration
	sm, err := structmapper.NewMapper()
	if err != nil {
		panic(err)
	}

	source := &ExampleStruct{
		A: "a",
		B: 5,
		C: []string{"c0", "c1"},
		D: map[string]int{
			"d0": 1,
			"d1": 2,
		},
		E:       net.ParseIP("127.0.0.1"),
		Ignored: "should be ignored",
	}

	// Convert to map
	m, err := sm.ToMap(source)
	if err != nil {
		panic(err)
	}

	target := &ExampleStruct{}

	// Convert back to struct
	if err := sm.ToStruct(m, target); err != nil {
		panic(err)
	}

	fmt.Printf("A=%s,B=%d,C=%s,D=d0=%d,d1=%d,E=%s,Ignored=%s", target.A, target.B, strings.Join(target.C, ","),
		target.D["d0"], target.D["d1"], target.E, target.Ignored)
	// Output: A=a,B=5,C=c0,c1,D=d0=1,d1=2,E=127.0.0.1,Ignored=
}
