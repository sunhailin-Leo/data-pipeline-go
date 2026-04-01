package convert

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// MatchFilters checks if the data matches all filter conditions
// data is a map[string]any (parsed JSON data)
// Returns true if all filters match (AND logic), false otherwise
func MatchFilters(data map[string]any, filters []config.TransformFilter) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		fieldValue, exists := data[filter.Field]
		if !exists {
			return false
		}
		if !matchSingleFilter(fieldValue, filter.Operator, filter.Value) {
			return false
		}
	}
	return true
}

func matchSingleFilter(fieldValue any, operator string, filterValue any) bool {
	switch operator {
	case "eq":
		return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", filterValue)
	case "neq":
		return fmt.Sprintf("%v", fieldValue) != fmt.Sprintf("%v", filterValue)
	case "gt":
		return cast.ToFloat64(fieldValue) > cast.ToFloat64(filterValue)
	case "gte":
		return cast.ToFloat64(fieldValue) >= cast.ToFloat64(filterValue)
	case "lt":
		return cast.ToFloat64(fieldValue) < cast.ToFloat64(filterValue)
	case "lte":
		return cast.ToFloat64(fieldValue) <= cast.ToFloat64(filterValue)
	case "contains":
		return strings.Contains(cast.ToString(fieldValue), cast.ToString(filterValue))
	case "not_contains":
		return !strings.Contains(cast.ToString(fieldValue), cast.ToString(filterValue))
	case "regex":
		pattern, err := regexp.Compile(cast.ToString(filterValue))
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error(utils.LogServiceName + "[Transform-Filter]invalid regex pattern: " + cast.ToString(filterValue))
			}
			return false
		}
		return pattern.MatchString(cast.ToString(fieldValue))
	case "in":
		return matchIn(fieldValue, filterValue)
	case "not_in":
		return !matchIn(fieldValue, filterValue)
	default:
		if logger.Logger != nil {
			logger.Logger.Error(utils.LogServiceName + "[Transform-Filter]unknown operator: " + operator)
		}
		return false
	}
}

func matchIn(fieldValue, filterValue any) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	switch values := filterValue.(type) {
	case []any:
		for _, v := range values {
			if fmt.Sprintf("%v", v) == fieldStr {
				return true
			}
		}
	case []string:
		for _, v := range values {
			if v == fieldStr {
				return true
			}
		}
	case string:
		// support comma-separated string: "a,b,c"
		for _, v := range strings.Split(values, ",") {
			if strings.TrimSpace(v) == fieldStr {
				return true
			}
		}
	}
	return false
}
