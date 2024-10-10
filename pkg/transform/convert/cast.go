package convert

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// generateConvertResultData generate result
func generateConvertResultData(sinkNames string, resultData map[string][]any, afterCastData any) {
	for _, sinkName := range strings.Split(sinkNames, ",") {
		resultData[sinkName] = append(resultData[sinkName], afterCastData)
	}
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
	default:
		logger.Logger.Error(utils.LogServiceName + "[CastTypes-convertMap]unknown convertor!")
		return nil
	}
}
