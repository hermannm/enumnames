// Package enumnames provides efficient mapping between integer enum values and string names.
package enumnames

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Map is an immutable mapping of integer enum values to string names.
// It must be instantiated with NewMap.
type Map[Enum IntegerEnum] struct {
	enumNames       []string
	lowestEnumValue Enum
}

type IntegerEnum interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// NewMap transforms the given enum-value-to-name map into a more efficient representation, which
// also provides utility methods for getting and marshaling enum names and values.
//
// Panics if:
//   - the range of enum values in the map is not contiguous
//   - there are duplicate names in the map
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
		index := enumValue - lowestEnumValue
		if int(index) >= capacity {
			panic("non-contiguous enum values given to enumnames.NewMap")
		}

		if slices.Contains(enumMap.enumNames, name) {
			panic(fmt.Sprintf("duplicate enum name '%s' given to enumnames.NewMap", name))
		}

		enumMap.enumNames[index] = name
	}

	return enumMap
}

// GetName gets the mapped name for the given enum value, or ok=false if it is not mapped.
func (enumMap Map[Enum]) GetName(enumValue Enum) (name string, ok bool) {
	index, inBounds := enumMap.index(enumValue)
	if !inBounds {
		return "", false
	}
	return enumMap.enumNames[index], true
}

// GetNameOrFallback gets the mapped name for the given enum value.
// If it is not mapped, returns the fallback.
func (enumMap Map[Enum]) GetNameOrFallback(enumValue Enum, fallback string) (name string) {
	if name, ok := enumMap.GetName(enumValue); ok {
		return name
	} else {
		return fallback
	}
}

// EnumValueFromName gets the corresponding enum value for the given name, or ok=false if no enum
// value is mapped to the name.
func (enumMap Map[Enum]) EnumValueFromName(name string) (enumValue Enum, ok bool) {
	for candidate, candidateName := range enumMap.enumNames {
		if candidateName == name {
			return Enum(candidate) + enumMap.lowestEnumValue, true
		}
	}

	return 0, false
}

// Contains checks if the given enum value exists in the map.
func (enumMap Map[Enum]) Contains(enumValue Enum) bool {
	_, inBounds := enumMap.index(enumValue)
	return inBounds
}

// Size returns the number of enum values in the map.
func (enumMap Map[Enum]) Size() int {
	return len(enumMap.enumNames)
}

// EnumValues returns a slice of all enum values in the map, sorted by their integer value.
// Mutating it will not affect the map.
func (enumMap Map[Enum]) EnumValues() []Enum {
	values := make([]Enum, len(enumMap.enumNames))
	for i := range enumMap.enumNames {
		values[i] = Enum(i) + enumMap.lowestEnumValue
	}
	return values
}

// Names returns a slice of all enum names in the map, sorted by their mapped enum integer value.
// Mutating it will not affect the map.
func (enumMap Map[Enum]) Names() []string {
	names := make([]string, len(enumMap.enumNames))
	copy(names, enumMap.enumNames)
	return names
}

// String returns a string repsentation of the map, mapping integer enum values to their names.
func (enumMap Map[Enum]) String() string {
	var builder strings.Builder

	lowestEnumValue := int(enumMap.lowestEnumValue)
	lastIndex := len(enumMap.enumNames) - 1

	builder.WriteString("enumnames.Map[")
	for i, name := range enumMap.enumNames {
		builder.WriteString(strconv.Itoa(i + lowestEnumValue))
		builder.WriteRune(':')
		builder.WriteString(name)
		if i != lastIndex {
			builder.WriteRune(' ')
		}
	}
	builder.WriteRune(']')

	return builder.String()
}

func (enumMap Map[Enum]) index(enumValue Enum) (index Enum, inBounds bool) {
	index = enumValue - enumMap.lowestEnumValue
	if index < 0 || int(index) >= len(enumMap.enumNames) {
		return 0, false
	}
	return index, true
}

// MarshalToNameJSON marshals the given enum value to its mapped name.
// It errors if the given enum value is not mapped.
func (enumMap Map[Enum]) MarshalToNameJSON(enumValue Enum) ([]byte, error) {
	if name, ok := enumMap.GetName(enumValue); ok {
		return json.Marshal(name)
	} else {
		return nil, fmt.Errorf("enum value '%d' not registered in enum name map", enumValue)
	}
}

// UnmarshalFromNameJSON unmarshals the given bytes with an enum name to the enum value pointed to
// by dest. It errors if it failed to unmarshal to string, or if the given enum name is not mapped.
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
