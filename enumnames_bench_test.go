package enumnames_test

import (
	"fmt"
	"testing"

	"hermannm.dev/enumnames"
)

var enumMap, nameMap, valueMap = makeBenchmarkMaps(255)

func BenchmarkGetEnumName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		name, _ := enumMap.GetName(1)
		_ = name
	}
}

func BenchmarkGetEnumNameWithMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		name := nameMap[1]
		_ = name
	}
}

func BenchmarkEnumValueFromName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		value, _ := enumMap.EnumValueFromName("Test 1")
		_ = value
	}
}

func BenchmarkEnumValueFromNameWithMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		value := valueMap["Test 1"]
		_ = value
	}
}

func makeBenchmarkMaps(
	size uint8,
) (enumMap enumnames.Map[TestEnum], nameMap map[TestEnum]string, valueMap map[string]TestEnum) {
	nameMap = make(map[TestEnum]string, int(size))
	for i := uint8(0); i < size; i++ {
		nameMap[TestEnum(i)] = fmt.Sprintf("Test %d", i)
	}

	valueMap = make(map[string]TestEnum, len(nameMap))
	for enumValue, name := range nameMap {
		valueMap[name] = enumValue
	}

	return enumnames.NewMap(nameMap), nameMap, valueMap
}
