package enumnames

import (
	"encoding/json"
	"fmt"
	"strings"
)

type IntegerEnum interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Map[Enum IntegerEnum] struct {
	enumNames       []string
	lowestEnumValue Enum
}

func NewMap[Enum IntegerEnum](enumValuesWithNames map[Enum]string) Map[Enum] {
	var lowestEnumValue Enum
	first := true
	for enumValue := range enumValuesWithNames {
		if first || enumValue < lowestEnumValue {
			lowestEnumValue = enumValue
			first = false
		}
	}

	capacity := len(enumValuesWithNames)

	enumMap := Map[Enum]{
		enumNames:       make([]string, capacity),
		lowestEnumValue: lowestEnumValue,
	}

	for enumValue, name := range enumValuesWithNames {
		if enumValue < 0 {
			panic("negative enum value given to enumnames.NewMap")
		}

		index := enumValue - lowestEnumValue
		if int(index) >= capacity {
			panic("non-contiguous enum values given to enumnames.NewMap")
		}

		enumMap.enumNames[index] = name
	}

	return enumMap
}

func (enumMap Map[Enum]) GetName(enumValue Enum) (name string, ok bool) {
	index, inBounds := enumMap.index(enumValue)
	if !inBounds {
		return "", false
	}
	return enumMap.enumNames[index], true
}

func (enumMap Map[Enum]) GetNameOrFallback(enumValue Enum, fallback string) (name string) {
	if name, ok := enumMap.GetName(enumValue); ok {
		return name
	} else {
		return fallback
	}
}

func (enumMap Map[Enum]) EnumValueFromName(name string) (enumValue Enum, ok bool) {
	for candidate, candidateName := range enumMap.enumNames {
		if candidateName == name {
			return Enum(candidate) + enumMap.lowestEnumValue, true
		}
	}

	return 0, false
}

func (enumMap Map[Enum]) Contains(enumValue Enum) bool {
	_, inBounds := enumMap.index(enumValue)
	return inBounds
}

func (enumMap Map[Enum]) Size() int {
	return len(enumMap.enumNames)
}

func (enumMap Map[Enum]) index(enumValue Enum) (index Enum, inBounds bool) {
	index = enumValue - enumMap.lowestEnumValue
	if index < 0 || int(index) >= len(enumMap.enumNames) {
		return 0, false
	}
	return index, true
}

func (enumMap Map[Enum]) MarshalToNameJSON(enumValue Enum) ([]byte, error) {
	if name, ok := enumMap.GetName(enumValue); ok {
		return json.Marshal(name)
	} else {
		return nil, fmt.Errorf("enum value '%d' not registered in enum name map", enumValue)
	}
}

func (enumMap Map[Enum]) UnmarshalFromNameJSON(bytes []byte, dest *Enum) error {
	var name string
	if err := json.Unmarshal(bytes, &name); err != nil {
		return err
	}

	if enumValue, ok := enumMap.EnumValueFromName(name); ok {
		*dest = enumValue
		return nil
	} else {
		return fmt.Errorf(
			"invalid value '%s', expected one of: '%s'",
			name,
			strings.Join(enumMap.enumNames, "', '"),
		)
	}
}
