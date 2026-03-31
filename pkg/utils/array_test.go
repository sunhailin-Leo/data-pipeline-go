package utils

import (
	"reflect"
	"testing"
)

func TestPopSlice(t *testing.T) {
	t.Run("normal slice", func(t *testing.T) {
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
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []any{}
		slice, pop := PopSlice(slice)
		if pop != nil {
			t.Errorf("Expected nil, got %v", pop)
		}
		if len(slice) != 0 {
			t.Errorf("Expected 0, got %v", len(slice))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var slice []any = nil
		slice, pop := PopSlice(slice)
		if pop != nil {
			t.Errorf("Expected nil, got %v", pop)
		}
		if len(slice) != 0 {
			t.Errorf("Expected 0, got %v", len(slice))
		}
	})

	t.Run("single element", func(t *testing.T) {
		slice := []any{42}
		slice, pop := PopSlice(slice)
		if pop != 42 {
			t.Errorf("Expected 42, got %v", pop)
		}
		if len(slice) != 0 {
			t.Errorf("Expected 0, got %v", len(slice))
		}
	})
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
		{
			Name:  "keys more than values",
			Keys:  []string{"a", "b", "c"},
			Value: []any{1, 2},
			Expected: map[string]any{
				"a": 1,
				"b": 2,
			},
		},
		{
			Name:  "keys less than values",
			Keys:  []string{"a", "b"},
			Value: []any{1, 2, 3},
			Expected: map[string]any{
				"a": 1,
				"b": 2,
			},
		},
		{
			Name:     "empty keys and values",
			Keys:     []string{},
			Value:    []any{},
			Expected: map[string]any{},
		},
		{
			Name:     "empty keys with values",
			Keys:     []string{},
			Value:    []any{1, 2, 3},
			Expected: map[string]any{},
		},
		{
			Name:     "keys with empty values",
			Keys:     []string{"a", "b", "c"},
			Value:    []any{},
			Expected: map[string]any{},
		},
		{
			Name:  "nil values",
			Keys:  []string{"a", "b"},
			Value: []any{nil, nil},
			Expected: map[string]any{
				"a": nil,
				"b": nil,
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

// Benchmark functions
func BenchmarkPopSlice(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := []any{1, 2, 3, 4, 5}
		PopSlice(slice)
	}
}

func BenchmarkInterfaceSliceToStringSlice(b *testing.B) {
	b.ReportAllocs()
	slice := []any{1, 2, 3, 4, 5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InterfaceSliceToStringSlice(slice)
	}
}

func BenchmarkStringSliceToMap(b *testing.B) {
	b.ReportAllocs()
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	values := []any{1, 2, 3, 4, 5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringSliceToMap(keys, values)
	}
}

func BenchmarkStringSliceToInterface(b *testing.B) {
	b.ReportAllocs()
	slice := []string{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringSliceToInterface(slice)
	}
}
