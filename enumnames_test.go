package enumnames_test

import (
	"encoding/json"
	"testing"

	"hermannm.dev/enumnames"
)

type TestEnum uint8

const (
	Test1 TestEnum = iota + 1
	Test2
	Test3
)

var testEnumMap = map[TestEnum]string{
	Test1: "FIRST",
	Test2: "SECOND",
	Test3: "THIRD",
}

var testEnumNames = enumnames.NewMap(testEnumMap)

func TestGetName(t *testing.T) {
	for enumValue, expectedName := range testEnumMap {
		name, ok := testEnumNames.GetName(enumValue)
		if !ok {
			t.Fatalf("expected '%s', got ok=false", expectedName)
		}
		if name != expectedName {
			t.Fatalf("expected '%s', got '%s'", expectedName, name)
		}
	}

	invalidEnumValue := TestEnum(100)
	_, ok := testEnumNames.GetName(invalidEnumValue)
	if ok {
		t.Fatal("expected GetName to return ok=false with invalid enum value")
	}
}

func TestGetNameorFallback(t *testing.T) {
	invalidEnumValue := TestEnum(100)
	fallback := "fallback"
	name := testEnumNames.GetNameOrFallback(invalidEnumValue, fallback)
	if name != fallback {
		t.Fatalf("expected '%s', got '%s'", fallback, name)
	}
}

func TestEnumValueFromName(t *testing.T) {
	for expectedEnumValue, name := range testEnumMap {
		enumValue, ok := testEnumNames.EnumValueFromName(name)
		if !ok {
			t.Fatalf("expected '%d', got ok=false", expectedEnumValue)
		}
		if enumValue != expectedEnumValue {
			t.Fatalf("expected '%d', got '%d'", expectedEnumValue, enumValue)
		}
	}

	invalidName := "garbage user input"
	_, ok := testEnumNames.EnumValueFromName(invalidName)
	if ok {
		t.Fatal("expected EnumValueFromName to return ok=false with invalid enum name")
	}
}

func TestContains(t *testing.T) {
	for enumValue := range testEnumMap {
		if !testEnumNames.Contains(enumValue) {
			t.Fatalf("expected enum names to contain entry for enum value '%d'", enumValue)
		}
	}

	invalidEnumValue := TestEnum(100)
	if testEnumNames.Contains(invalidEnumValue) {
		t.Fatal("expected Contains to return false for invalid enum value")
	}
}

func TestSize(t *testing.T) {
	size := testEnumNames.Size()
	if size != 3 {
		t.Fatalf("expected enum map size to be 3, got %d", size)
	}
}

func TestEnumValues(t *testing.T) {
	enumValues := testEnumNames.EnumValues()

	if len(enumValues) != 3 {
		t.Fatalf("expected enum values with length 3, got %+v", enumValues)
	}

	if enumValues[0] != Test1 || enumValues[1] != Test2 || enumValues[2] != Test3 {
		t.Fatalf("expected [Test1, Test2, Test3], got %+v", enumValues)
	}
}

func TestNames(t *testing.T) {
	names := testEnumNames.Names()

	if len(names) != 3 {
		t.Fatalf("expected enum names with length 3, got %+v", names)
	}

	if names[0] != "FIRST" || names[1] != "SECOND" || names[2] != "THIRD" {
		t.Fatalf("expected ['FIRST', 'SECOND', 'THIRD'], got %+v", names)
	}
}

func TestString(t *testing.T) {
	expected := "enumnames.Map[1:FIRST 2:SECOND 3:THIRD]"
	actual := testEnumNames.String()
	if expected != actual {
		t.Fatalf("expected '%s', got '%s'", expected, actual)
	}
}

func TestNegativeEnum(t *testing.T) {
	type NegativeEnum int8

	const (
		Neg1 NegativeEnum = -2
		Neg2 NegativeEnum = -1
		Neg3 NegativeEnum = 0
	)

	negativeEnumNames := enumnames.NewMap(map[NegativeEnum]string{
		Neg1: "FIRST",
		Neg2: "SECOND",
		Neg3: "THIRD",
	})

	name, ok := negativeEnumNames.GetName(Neg2)
	expectedName := "SECOND"
	if !ok {
		t.Fatalf("expected '%s', got ok=false", expectedName)
	}
	if name != expectedName {
		t.Fatalf("expected '%s', got '%s'", expectedName, name)
	}
}

func (test TestEnum) MarshalJSON() ([]byte, error) {
	return testEnumNames.MarshalToNameJSON(test)
}

func (test *TestEnum) UnmarshalJSON(bytes []byte) error {
	return testEnumNames.UnmarshalFromNameJSON(bytes, test)
}

type JSONExample struct {
	EnumField TestEnum `json:"enumField"`
}

func TestMarshalToNameJSON(t *testing.T) {
	example := JSONExample{EnumField: Test1}

	bytes, err := json.Marshal(example)
	if err != nil {
		t.Fatalf("expected JSON marshaling of enum value to succeed, but got error: %v", err)
	}

	expectedMarshalValue := `{"enumField":"FIRST"}`
	if string(bytes) != expectedMarshalValue {
		t.Fatalf("expected '%s', got '%s'", expectedMarshalValue, string(bytes))
	}

	example.EnumField = TestEnum(100)
	_, err = json.Marshal(example)
	if err == nil {
		t.Fatal("expected JSON marshaling to fail for invalid enum value")
	}
}

func TestUnmarshalFromNameJSON(t *testing.T) {
	jsonInput := []byte(`{"enumField":"FIRST"}`)

	var result JSONExample
	if err := json.Unmarshal(jsonInput, &result); err != nil {
		t.Fatalf("expected JSON unmarshaling of enum value to succeed, but got error: %v", err)
	}

	if result.EnumField != Test1 {
		t.Fatalf("expected '%d', got '%d'", Test1, result.EnumField)
	}

	invalidJSONInput := []byte(`{"enumField":"garbage user input"}`)
	var result2 JSONExample
	if err := json.Unmarshal(invalidJSONInput, &result2); err == nil {
		t.Fatal("expected JSON unmarshaling to fail for invalid enum name")
	}
}
