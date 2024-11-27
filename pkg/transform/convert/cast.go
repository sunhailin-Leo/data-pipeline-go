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

// CastTypes data type conversion
func CastTypes(data any, convertorName string) any {
	if convertorName == "" {
		return data
	}
	// based on convertorName to convert
	switch convertorName {
	case "toBool":
		v, err := cast.ToBoolE(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toBool]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toFloat64":
		v, err := cast.ToFloat64E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toFloat64]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toFloat32":
		v, err := cast.ToFloat32E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toFloat32]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toInt64":
		v, err := cast.ToInt64E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toInt64]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toInt32":
		v, err := cast.ToInt32E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toInt32]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toInt16":
		v, err := cast.ToInt16E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toInt16]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toInt8":
		v, err := cast.ToInt8E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toInt8]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toInt":
		v, err := cast.ToIntE(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toInt]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toUint":
		v, err := cast.ToUintE(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toUint]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toUint64":
		v, err := cast.ToUint64E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toUint64]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toUint32":
		v, err := cast.ToUint32E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toUint32]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toUint16":
		v, err := cast.ToUint16E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toUint16]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toUint8":
		v, err := cast.ToUint8E(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toUint8]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toString":
		v, err := cast.ToStringE(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toString]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	case "toStringMap":
		v, err := cast.ToStringMapE(data)
		if err != nil {
			logger.Logger.Error(utils.LogServiceName + "[CastTypes-toStringMap]Data conversion failed! Reason for error: " + err.Error())
			return nil
		}
		return v
	default:
		logger.Logger.Error(utils.LogServiceName + "[CastTypes-convertMap]unknown convertor!")
		return nil
	}
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
