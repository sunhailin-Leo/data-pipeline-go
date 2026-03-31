package convert

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/gjson"
	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// generateConvertResultData generate result
func generateConvertResultData(sinkNames string, resultData map[string][]any, afterCastData any) {
	for _, sinkName := range strings.Split(sinkNames, ",") {
		resultData[sinkName] = append(resultData[sinkName], afterCastData)
	}
}

// JsonPathToMap json path to map
func JsonPathToMap(data []byte, paths []config.TransformJsonPath) map[string]any {
	result := make(map[string]any)
	dataStr := gjson.ParseBytes(data)
	for _, path := range paths {
		pathRes := dataStr.Get(path.Path).Value()
		if pathRes == nil || path.DestField == "" {
			logger.Logger.Error(utils.LogServiceName + "[JsonPathToMap]Failed to get jsonPath result! Error Path: " + path.Path)
			continue
		}
		result[path.DestField] = pathRes
	}

	return result
}

// JsonToMap json to map
func JsonToMap(data []byte) map[string]any {
	result := make(map[string]any)
	unmarshalErr := sonic.Unmarshal(data, &result)
	if unmarshalErr != nil {
		logger.Logger.Error(utils.LogServiceName + "[JsonToMap]Data conversion Json-Map failed! Error Cause: " + unmarshalErr.Error())
		return nil
	}
	return result
}

// MapToJson map to json
func MapToJson(data map[string]any) []byte {
	result, marshalErr := sonic.Marshal(data)
	if marshalErr != nil {
		logger.Logger.Error(utils.LogServiceName + "[MapToJson]Data conversion Map-Json failed! Error Cause: " + marshalErr.Error())
		return nil
	}
	return result
}

// castFunc is a type converter function that returns the converted value and an error.
type castFunc func(any) (any, error)

// castFuncMap maps converter names to their corresponding cast functions.
var castFuncMap = map[string]castFunc{
	"toBool":      func(v any) (any, error) { return cast.ToBoolE(v) },
	"toFloat64":   func(v any) (any, error) { return cast.ToFloat64E(v) },
	"toFloat32":   func(v any) (any, error) { return cast.ToFloat32E(v) },
	"toInt64":     func(v any) (any, error) { return cast.ToInt64E(v) },
	"toInt32":     func(v any) (any, error) { return cast.ToInt32E(v) },
	"toInt16":     func(v any) (any, error) { return cast.ToInt16E(v) },
	"toInt8":      func(v any) (any, error) { return cast.ToInt8E(v) },
	"toInt":       func(v any) (any, error) { return cast.ToIntE(v) },
	"toUint":      func(v any) (any, error) { return cast.ToUintE(v) },
	"toUint64":    func(v any) (any, error) { return cast.ToUint64E(v) },
	"toUint32":    func(v any) (any, error) { return cast.ToUint32E(v) },
	"toUint16":    func(v any) (any, error) { return cast.ToUint16E(v) },
	"toUint8":     func(v any) (any, error) { return cast.ToUint8E(v) },
	"toString":    func(v any) (any, error) { return cast.ToStringE(v) },
	"toStringMap": func(v any) (any, error) { return cast.ToStringMapE(v) },
}

// CastTypes data type conversion
func CastTypes(data any, convertorName string) any {
	if convertorName == "" {
		return data
	}

	fn, exists := castFuncMap[convertorName]
	if !exists {
		logger.Logger.Error(utils.LogServiceName + "[CastTypes]unknown convertor: " + convertorName)
		return nil
	}

	result, err := fn(data)
	if err != nil {
		logger.Logger.Error(utils.LogServiceName + "[CastTypes-" + convertorName + "]Data conversion failed! Reason for error: " + err.Error())
		return nil
	}
	return result
}

// CastTypesDefaultValue data type default value
func CastTypesDefaultValue(convertorName string) any {
	if convertorName == "" {
		return nil
	}

	switch convertorName {
	case "toBool":
		return false
	case "toFloat64", "toFloat32":
		return 0.0
	case "toInt64", "toInt32", "toInt16", "toInt8", "toInt", "toUint", "toUint64", "toUint32", "toUint16", "toUint8":
		return 0
	case "toString":
		return ""
	default:
		return nil
	}
}

// CastFunctionNameToFunctionResult if value is special function then return function result
func CastFunctionNameToFunctionResult(value any) any {
	switch value {
	case "$.UUID()":
		return utils.GetUUID()
	default:
		return value
	}
}
