package enumnames_test

import (
	"fmt"
	"testing"

	"hermannm.dev/enumnames"
)

var enumMap, nameMap, reverseNameMap = makeBenchmarkMaps(255)

// Global variables to avoid the compiler optimizing away our benchmarked function calls
// (see https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
//
//goland:noinspection GoUnusedGlobalVariable
var (
	globalName string
	globalKey  uint8
)

func BenchmarkGetName(b *testing.B) {
	var name string

	for range b.N {
		name, _ = enumMap.GetName(1)
	}

	globalName = name
}

func BenchmarkGetNameWithMap(b *testing.B) {
	var name string

	for range b.N {
		name = nameMap[1]
	}

	globalName = name
}

func BenchmarkGetKey(b *testing.B) {
	var value uint8

	for range b.N {
		value, _ = enumMap.GetKey("Test 1")
	}

	globalKey = value
}

func BenchmarkGetKeyWithMap(b *testing.B) {
	var value uint8

	for range b.N {
		value = reverseNameMap["Test 1"]
	}

	globalKey = value
}

func makeBenchmarkMaps(size uint8) (
	enumMap enumnames.Map[uint8],
	nameMap map[uint8]string,
	reverseNameMap map[string]uint8,
) {
	nameMap = make(map[uint8]string, int(size))
	for i := range size {
		nameMap[i] = fmt.Sprintf("Test %d", i)
	}

	reverseNameMap = make(map[string]uint8, len(nameMap))
	for enumValue, name := range nameMap {
		reverseNameMap[name] = enumValue
	}

	return enumnames.NewMap(nameMap), nameMap, reverseNameMap
}
