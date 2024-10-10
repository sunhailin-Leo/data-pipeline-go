package utils

import (
	"github.com/spf13/cast"
)

// PopSlice return slice after pop and pop element
func PopSlice(slice []any) ([]any, any) {
	if len(slice) == 0 {
		return slice, nil
	}
	return slice[1:], slice[0]
}

// InterfaceSliceToStringSlice []any to []string
func InterfaceSliceToStringSlice(slice []any) []string {
	s := make([]string, len(slice))
	for i, x := range slice {
		s[i] = cast.ToString(x)
	}
	return s
}

// StringSliceToMap string slice to map
func StringSliceToMap(keys []string, value []any) map[string]any {
	result := make(map[string]any)
	for i, key := range keys {
		result[key] = value[i]
	}
	return result
}
