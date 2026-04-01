package convert

import (
	"strings"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// convertToJsonBytes converts result data to JSON bytes using specialProcessKeys.
func convertToJsonBytes(specialProcessKeys []string, result []any) []any {
	return []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
}

// convertToRowBytes concatenates all items as strings and returns as bytes.
func convertToRowBytes(result []any) []any {
	var builder strings.Builder
	for _, item := range result {
		builder.WriteString(CastTypes(item, "toString").(string))
	}
	return []any{[]byte(builder.String())}
}

// processByMessageMode processes data based on message mode (json/row).
func processByMessageMode(messageMode string, resultData map[string][]any, sinkName string, specialProcessKeys []string, result []any) {
	switch messageMode {
	case utils.TransformJsonMode:
		resultData[sinkName] = convertToJsonBytes(specialProcessKeys, result)
	case utils.TransformRowMode:
		resultData[sinkName] = convertToRowBytes(result)
	}
}

// getMessageMode extracts the message mode from a sink config.
func getMessageMode(sinkCfg *config.SinkConfig) string {
	switch sinkCfg.Type {
	case utils.SinkKafkaTagName:
		return sinkCfg.Kafka.MessageMode
	case utils.SinkRocketMQTagName:
		return sinkCfg.RocketMQ.MessageMode
	case utils.SinkRabbitMQTagName:
		return sinkCfg.RabbitMQ.MessageMode
	case utils.SinkPulsarTagName:
		return sinkCfg.Pulsar.MessageMode
	default:
		return ""
	}
}

// processBySink Processing of data according to Sink
func processBySink(resultData map[string][]any, streamConfig *config.StreamConfig, specialProcessKeys []string) {
	for sinkName, result := range resultData {
		sinkCfg := streamConfig.GetSinkBySinkName(sinkName)
		if sinkCfg == nil {
			continue
		}

		switch sinkCfg.Type {
		case utils.SinkHTTPTagName:
			if sinkCfg.HTTP.ContentType == utils.HTTPContentTypeJSON {
				resultData[sinkName] = convertToJsonBytes(specialProcessKeys, result)
			}
		case utils.SinkKafkaTagName, utils.SinkRocketMQTagName, utils.SinkRabbitMQTagName, utils.SinkPulsarTagName:
			processByMessageMode(getMessageMode(sinkCfg), resultData, sinkName, specialProcessKeys, result)
		}
	}
}

// JsonMode JSON mode
func JsonMode(beforeConvertData *models.TransformBeforeConvert, runningCfg []config.TransformSchema, streamConfig *config.StreamConfig) *models.TransformAfterConvert {
	// init
	specialProcessKeys := make([]string, 0, len(runningCfg))
	resultData := make(map[string][]any, len(streamConfig.Sink))
	// Ensure that the final data is in order
	jsonData := beforeConvertData.BeforeConvertData.(map[string]any)
	for _, schema := range runningCfg {
		// Whether to ignore the field (remove it from the sender's data based on the sourceKey, currently only the top-level key is supported).
		if schema.IsIgnore {
			delete(jsonData, schema.SourceKey)
			continue
		}
		// Whether to use strict mode for this field (sender must provide)
		if schema.IsStrictMode {
			if _, keyExist := jsonData[schema.SourceKey]; !keyExist {
				// In the future, after upgrading to 1.23, you can add a maps.Keys to locate the problem in the logs.
				logger.Logger.Error(utils.LogServiceName + "[JsonMode]Data transfer field does not match the given parse, Parse field: " + schema.SourceKey)
				return nil
			}
		}
		// Special record keys (for special scenarios)
		if schema.IsKeepKeys {
			specialProcessKeys = append(specialProcessKeys, schema.SourceKey)
		}
		// Whether it is an extended field (allow sender not to provide it, add it in transform module)
		if schema.IsExpand && jsonData[schema.SourceKey] == nil {
			jsonData[schema.SourceKey] = CastFunctionNameToFunctionResult(schema.ExpandValue)
		}
		// Data conversion and generation of corresponding structures
		generateConvertResultData(schema.SinkName, resultData,
			CastTypes(jsonData[schema.SourceKey], schema.Converter))
	}
	// Determine if SinkTag to see if special handling is needed.
	if len(specialProcessKeys) == 0 {
		return &models.TransformAfterConvert{
			SourceOutput:     beforeConvertData.SourceOutput,
			AfterConvertData: resultData,
		}
	}
	// Special handling logic
	processBySink(resultData, streamConfig, specialProcessKeys)

	return &models.TransformAfterConvert{
		SourceOutput:     beforeConvertData.SourceOutput,
		AfterConvertData: resultData,
	}
}
