package enumnames_test

import (
	"fmt"
	"testing"

	"hermannm.dev/enumnames"
)

var enumMap, nameMap, valueMap = makeBenchmarkMaps(255)

// Global variables to avoid the compiler optimizing away our benchmarked function calls
// (see https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
var globalName string
var globalValue uint8

func BenchmarkGetName(b *testing.B) {
	var name string

	for i := 0; i < b.N; i++ {
		name, _ = enumMap.GetName(1)
	}

	globalName = name
}

func BenchmarkGetNameWithMap(b *testing.B) {
	var name string

	for i := 0; i < b.N; i++ {
		name = nameMap[1]
	}

	globalName = name
}

func BenchmarkEnumValueFromName(b *testing.B) {
	var value uint8

	for i := 0; i < b.N; i++ {
		value, _ = enumMap.EnumValueFromName("Test 1")
	}

	globalValue = value
}

func BenchmarkEnumValueFromNameWithMap(b *testing.B) {
	var value uint8

	for i := 0; i < b.N; i++ {
		value = valueMap["Test 1"]
	}

	globalValue = value
}

func makeBenchmarkMaps(size uint8) (
	enumMap enumnames.Map[uint8],
	nameMap map[uint8]string,
	valueMap map[string]uint8,
) {
	nameMap = make(map[uint8]string, int(size))
	for i := uint8(0); i < size; i++ {
		nameMap[i] = fmt.Sprintf("Test %d", i)
	}

	valueMap = make(map[string]uint8, len(nameMap))
	for enumValue, name := range nameMap {
		valueMap[name] = enumValue
	}

	return enumnames.NewMap(nameMap), nameMap, valueMap
}
