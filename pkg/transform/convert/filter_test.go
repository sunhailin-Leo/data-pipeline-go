package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
)

// TestMatchFilters_EmptyFilters tests that empty filters return true
func TestMatchFilters_EmptyFilters(t *testing.T) {
	data := map[string]any{"name": "test", "age": 25}
	filters := []config.TransformFilter{}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Empty filters should return true")
	}
}

// TestMatchFilters_EqOperator tests the eq operator
func TestMatchFilters_EqOperator(t *testing.T) {
	data := map[string]any{"name": "test", "age": 25}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "test"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Eq operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "other"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Eq operator should not match")
	}
}

// TestMatchFilters_NeqOperator tests the neq operator
func TestMatchFilters_NeqOperator(t *testing.T) {
	data := map[string]any{"name": "test", "age": 25}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "neq", Value: "other"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Neq operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "name", Operator: "neq", Value: "test"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Neq operator should not match")
	}
}

// TestMatchFilters_GtOperator tests the gt operator
func TestMatchFilters_GtOperator(t *testing.T) {
	data := map[string]any{"age": 25}
	filters := []config.TransformFilter{
		{Field: "age", Operator: "gt", Value: 20},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Gt operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "age", Operator: "gt", Value: 25},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Gt operator should not match")
	}
}

// TestMatchFilters_GteOperator tests the gte operator
func TestMatchFilters_GteOperator(t *testing.T) {
	data := map[string]any{"age": 25}
	filters := []config.TransformFilter{
		{Field: "age", Operator: "gte", Value: 25},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Gte operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "age", Operator: "gte", Value: 26},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Gte operator should not match")
	}
}

// TestMatchFilters_LtOperator tests the lt operator
func TestMatchFilters_LtOperator(t *testing.T) {
	data := map[string]any{"age": 25}
	filters := []config.TransformFilter{
		{Field: "age", Operator: "lt", Value: 30},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Lt operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "age", Operator: "lt", Value: 25},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Lt operator should not match")
	}
}

// TestMatchFilters_LteOperator tests the lte operator
func TestMatchFilters_LteOperator(t *testing.T) {
	data := map[string]any{"age": 25}
	filters := []config.TransformFilter{
		{Field: "age", Operator: "lte", Value: 25},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Lte operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "age", Operator: "lte", Value: 24},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Lte operator should not match")
	}
}

// TestMatchFilters_ContainsOperator tests the contains operator
func TestMatchFilters_ContainsOperator(t *testing.T) {
	data := map[string]any{"message": "hello world"}
	filters := []config.TransformFilter{
		{Field: "message", Operator: "contains", Value: "hello"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Contains operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "message", Operator: "contains", Value: "goodbye"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Contains operator should not match")
	}
}

// TestMatchFilters_NotContainsOperator tests the not_contains operator
func TestMatchFilters_NotContainsOperator(t *testing.T) {
	data := map[string]any{"message": "hello world"}
	filters := []config.TransformFilter{
		{Field: "message", Operator: "not_contains", Value: "goodbye"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Not_contains operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "message", Operator: "not_contains", Value: "hello"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Not_contains operator should not match")
	}
}

// TestMatchFilters_RegexOperator tests the regex operator
func TestMatchFilters_RegexOperator(t *testing.T) {
	data := map[string]any{"email": "test@example.com"}
	filters := []config.TransformFilter{
		{Field: "email", Operator: "regex", Value: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Regex operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "email", Operator: "regex", Value: "^invalid$"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Regex operator should not match")
	}
}

// TestMatchFilters_InOperator tests the in operator
func TestMatchFilters_InOperator(t *testing.T) {
	data := map[string]any{"status": "active"}
	filters := []config.TransformFilter{
		{Field: "status", Operator: "in", Value: []any{"active", "pending", "completed"}},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("In operator should match with []any")
	}

	filters = []config.TransformFilter{
		{Field: "status", Operator: "in", Value: []string{"active", "pending", "completed"}},
	}
	result = MatchFilters(data, filters)
	if !result {
		t.Error("In operator should match with []string")
	}

	filters = []config.TransformFilter{
		{Field: "status", Operator: "in", Value: "active,pending,completed"},
	}
	result = MatchFilters(data, filters)
	if !result {
		t.Error("In operator should match with comma-separated string")
	}

	filters = []config.TransformFilter{
		{Field: "status", Operator: "in", Value: []any{"pending", "completed"}},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("In operator should not match")
	}
}

// TestMatchFilters_NotInOperator tests the not_in operator
func TestMatchFilters_NotInOperator(t *testing.T) {
	data := map[string]any{"status": "inactive"}
	filters := []config.TransformFilter{
		{Field: "status", Operator: "not_in", Value: []any{"active", "pending", "completed"}},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("Not_in operator should match")
	}

	filters = []config.TransformFilter{
		{Field: "status", Operator: "not_in", Value: []any{"inactive", "pending"}},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Not_in operator should not match")
	}
}

// TestMatchFilters_FieldNotExists tests when field does not exist
func TestMatchFilters_FieldNotExists(t *testing.T) {
	data := map[string]any{"name": "test"}
	filters := []config.TransformFilter{
		{Field: "age", Operator: "eq", Value: 25},
	}
	result := MatchFilters(data, filters)
	if result {
		t.Error("Should return false when field does not exist")
	}
}

// TestMatchFilters_UnknownOperator tests unknown operator
func TestMatchFilters_UnknownOperator(t *testing.T) {
	data := map[string]any{"name": "test"}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "unknown", Value: "test"},
	}
	result := MatchFilters(data, filters)
	if result {
		t.Error("Should return false for unknown operator")
	}
}

// TestMatchFilters_MultipleFiltersAndLogic tests AND logic with multiple filters
func TestMatchFilters_MultipleFiltersAndLogic(t *testing.T) {
	data := map[string]any{"name": "test", "age": 25, "status": "active"}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "test"},
		{Field: "age", Operator: "gt", Value: 20},
		{Field: "status", Operator: "eq", Value: "active"},
	}
	result := MatchFilters(data, filters)
	if !result {
		t.Error("All filters should match with AND logic")
	}

	filters = []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "test"},
		{Field: "age", Operator: "gt", Value: 30},
		{Field: "status", Operator: "eq", Value: "active"},
	}
	result = MatchFilters(data, filters)
	if result {
		t.Error("Should return false when one filter does not match")
	}
}

// TestMatchIn tests matchIn function directly
func TestMatchIn(t *testing.T) {
	tests := []struct {
		name        string
		fieldValue  any
		filterValue any
		expected    bool
	}{
		{
			name:        "match with []any",
			fieldValue:  "active",
			filterValue: []any{"active", "pending"},
			expected:    true,
		},
		{
			name:        "match with []string",
			fieldValue:  "active",
			filterValue: []string{"active", "pending"},
			expected:    true,
		},
		{
			name:        "match with comma-separated string",
			fieldValue:  "active",
			filterValue: "active,pending",
			expected:    true,
		},
		{
			name:        "no match with []any",
			fieldValue:  "inactive",
			filterValue: []any{"active", "pending"},
			expected:    false,
		},
		{
			name:        "no match with []string",
			fieldValue:  "inactive",
			filterValue: []string{"active", "pending"},
			expected:    false,
		},
		{
			name:        "no match with comma-separated string",
			fieldValue:  "inactive",
			filterValue: "active,pending",
			expected:    false,
		},
		{
			name:        "match with spaces in comma-separated string",
			fieldValue:  "active",
			filterValue: "active, pending, completed",
			expected:    true,
		},
		{
			name:        "numeric value match",
			fieldValue:  123,
			filterValue: []any{123, 456, 789},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchIn(tt.fieldValue, tt.filterValue)
			if result != tt.expected {
				t.Errorf("matchIn(%v, %v) = %v, want %v", tt.fieldValue, tt.filterValue, result, tt.expected)
			}
		})
	}
}

// Benchmark functions
func BenchmarkMatchFilters(b *testing.B) {
	b.ReportAllocs()
	data := map[string]any{
		"name":   "test",
		"age":    25,
		"status": "active",
	}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "test"},
		{Field: "age", Operator: "gt", Value: 20},
		{Field: "status", Operator: "eq", Value: "active"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MatchFilters(data, filters)
	}
}

func BenchmarkMatchFilters_SingleFilter(b *testing.B) {
	b.ReportAllocs()
	data := map[string]any{"name": "test"}
	filters := []config.TransformFilter{
		{Field: "name", Operator: "eq", Value: "test"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MatchFilters(data, filters)
	}
}

func BenchmarkMatchFilters_InOperator(b *testing.B) {
	b.ReportAllocs()
	data := map[string]any{"status": "active"}
	filters := []config.TransformFilter{
		{Field: "status", Operator: "in", Value: []any{"active", "pending", "completed"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MatchFilters(data, filters)
	}
}

func TestMatchSingleFilter_Regex_InvalidPattern(t *testing.T) {
	initLogger()
	// Invalid regex pattern should return false
	result := matchSingleFilter("test@example.com", "regex", "(?P<name")
	assert.False(t, result)
}

func TestMatchSingleFilter_NotIn(t *testing.T) {
	// Test not_in with value not in list
	result := matchSingleFilter("inactive", "not_in", []any{"active", "pending", "completed"})
	assert.True(t, result)

	// Test not_in with value in list
	result = matchSingleFilter("active", "not_in", []any{"active", "pending", "completed"})
	assert.False(t, result)

	// Test not_in with []string
	result = matchSingleFilter("inactive", "not_in", []string{"active", "pending", "completed"})
	assert.True(t, result)

	// Test not_in with comma-separated string
	result = matchSingleFilter("inactive", "not_in", "active,pending,completed")
	assert.True(t, result)

	result = matchSingleFilter("active", "not_in", "active,pending,completed")
	assert.False(t, result)
}
