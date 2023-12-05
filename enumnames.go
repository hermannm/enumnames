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
	names     []string
	lowestKey Enum
}

type IntegerEnum interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// NewMap transforms the given enum-value-to-name map into a more efficient representation, which
// also provides utility methods for getting and marshaling enum names and values.
//
// Panics if:
//   - the range of integer enum keys in the map is not contiguous
//   - there are duplicate names in the map
func NewMap[Enum IntegerEnum](enumNames map[Enum]string) Map[Enum] {
	var lowestKey Enum
	first := true
	for key := range enumNames {
		if first || key < lowestKey {
			lowestKey = key
			first = false
		}
	}

	size := len(enumNames)
	enumMap := Map[Enum]{
		names:     make([]string, size),
		lowestKey: lowestKey,
	}

	for key, name := range enumNames {
		index := enumMap.keyToIndex(key)
		if int(index) >= size {
			panic("non-contiguous enum keys given to enumnames.NewMap")
		}
		if slices.Contains(enumMap.names, name) {
			panic(fmt.Sprintf("duplicate enum name '%s' given to enumnames.NewMap", name))
		}
		enumMap.names[index] = name
	}

	return enumMap
}

// GetName returns the mapped name for the given enum key, or ok=false if no mapping is found.
func (enumMap Map[Enum]) GetName(key Enum) (name string, ok bool) {
	if !enumMap.ContainsKey(key) {
		return "", false
	}

	index := enumMap.keyToIndex(key)
	return enumMap.names[index], true
}

// GetNameOrFallback returns the mapped name for the given enum key, or the fallback if no mapping
// is found.
func (enumMap Map[Enum]) GetNameOrFallback(key Enum, fallback string) (name string) {
	if !enumMap.ContainsKey(key) {
		return fallback
	}

	index := enumMap.keyToIndex(key)
	return enumMap.names[index]
}

// GetKey returns the enum key mapped to the given name, or ok=false if no mapping is found.
func (enumMap Map[Enum]) GetKey(name string) (key Enum, ok bool) {
	for i, candidate := range enumMap.names {
		if candidate == name {
			return enumMap.indexToKey(i), true
		}
	}

	return 0, false
}

// ContainsKey checks if the given enum key exists in the map.
func (enumMap Map[Enum]) ContainsKey(key Enum) bool {
	return key >= enumMap.lowestKey &&
		key < Enum(len(enumMap.names))+enumMap.lowestKey
}

// ContainsName checks if any enum key maps to the given name.
func (enumMap Map[Enum]) ContainsName(name string) bool {
	for _, candidate := range enumMap.names {
		if candidate == name {
			return true
		}
	}
	return false
}

// Size returns the number of enum-to-name entries in the map.
func (enumMap Map[Enum]) Size() int {
	return len(enumMap.names)
}

// Keys returns a slice of all enum keys in the map, sorted by their integer value.
// Mutating it will not affect the map.
func (enumMap Map[Enum]) Keys() []Enum {
	keys := make([]Enum, len(enumMap.names))
	for i := range enumMap.names {
		keys[i] = enumMap.indexToKey(i)
	}
	return keys
}

// Names returns a slice of all enum names in the map, sorted by the integer value of their keys.
// Mutating it will not affect the map.
func (enumMap Map[Enum]) Names() []string {
	names := make([]string, len(enumMap.names))
	copy(names, enumMap.names)
	return names
}

// String returns a string representation of the map, mapping integer enum keys to their names.
func (enumMap Map[Enum]) String() string {
	var builder strings.Builder

	lowestKey := int(enumMap.lowestKey)
	lastIndex := len(enumMap.names) - 1

	builder.WriteString("enumnames.Map[")
	for i, name := range enumMap.names {
		builder.WriteString(strconv.Itoa(i + lowestKey))
		builder.WriteRune(':')
		builder.WriteString(name)
		if i != lastIndex {
			builder.WriteRune(' ')
		}
	}
	builder.WriteRune(']')

	return builder.String()
}

// MarshalToNameJSON marshals the given enum key to its mapped name.
// It errors if the key is not mapped.
func (enumMap Map[Enum]) MarshalToNameJSON(key Enum) ([]byte, error) {
	if name, ok := enumMap.GetName(key); ok {
		return json.Marshal(name)
	} else {
		return nil, fmt.Errorf("invalid value '%d': key not found in enum name map", key)
	}
}

// UnmarshalFromNameJSON unmarshals the given enum name JSON to string, and sets dest to the enum
// key mapped to the name.
// It errors if string unmarshaling fails, or if the unmarshaled enum name is not mapped.
func (enumMap Map[Enum]) UnmarshalFromNameJSON(nameJSON []byte, dest *Enum) error {
	var name string
	if err := json.Unmarshal(nameJSON, &name); err != nil {
		return err
	}

	if key, ok := enumMap.GetKey(name); ok {
		*dest = key
		return nil
	} else {
		return fmt.Errorf(
			"invalid value '%s', expected one of: '%s'",
			name,
			strings.Join(enumMap.names, "', '"),
		)
	}
}

func (enumMap Map[Enum]) keyToIndex(key Enum) (index Enum) {
	return key - enumMap.lowestKey
}

func (enumMap Map[Enum]) indexToKey(index int) (key Enum) {
	return Enum(index) + enumMap.lowestKey
}
