package utils

import (
	"reflect"
	"testing"
)

func TestPopSlice(t *testing.T) {
	slice := []any{1, 2, 3}
	slice, pop := PopSlice(slice)
	if pop != 1 {
		t.Errorf("Expected 1, got %v", pop)
	}
	if len(slice) != 2 {
		t.Errorf("Expected 2, got %v", len(slice))
	}
	if slice[0] != 2 {
		t.Errorf("Expected 2, got %v", slice[0])
	}
	if slice[1] != 3 {
		t.Errorf("Expected 3, got %v", slice[1])
	}
}

func TestInterfaceSliceToStringSlice(t *testing.T) {
	slice := []interface{}{1, 2, 3}
	result := InterfaceSliceToStringSlice(slice)
	if len(result) != 3 {
		t.Errorf("Expected 3, got %v", len(result))
	}
	if result[0] != "1" {
		t.Errorf("Expected 1, got %v", result[0])
	}
	if result[1] != "2" {
		t.Errorf("Expected 2, got %v", result[1])
	}
	if result[2] != "3" {
		t.Errorf("Expected 3, got %v", result[2])
	}
}

func TestStringSliceToMap(t *testing.T) {
	t.Helper()

	type testCaseUnit struct {
		Name     string
		Keys     []string
		Value    []any
		Expected map[string]any
	}

	testUnits := []testCaseUnit{
		{
			Name:  "test-1",
			Keys:  []string{"a", "b"},
			Value: []any{1, 2},
			Expected: map[string]any{
				"a": 1,
				"b": 2,
			},
		},
	}

	for _, unit := range testUnits {
		actual := StringSliceToMap(unit.Keys, unit.Value)
		if !reflect.DeepEqual(actual, unit.Expected) {
			t.Fatalf("TestArrayToMap failed. TestName: %s, Expected: %v, actual: %v", unit.Name, unit.Expected, actual)
		}
	}
}

func TestStringSliceToInterface(t *testing.T) {
	tests := []struct {
		input    []string
		expected []interface{}
	}{
		{
			input:    []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			input:    []string{},
			expected: []interface{}{},
		},
		{
			input:    []string{"1", "2", "3"},
			expected: []interface{}{"1", "2", "3"},
		},
		{
			input:    []string{"hello", "world"},
			expected: []interface{}{"hello", "world"},
		},
	}

	for _, test := range tests {
		result := StringSliceToInterface(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input %v, expected %v, but got %v", test.input, test.expected, result)
		}
	}
}
