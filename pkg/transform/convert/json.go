package convert

import (
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

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
				resultData[sinkName] = []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
			}
		case utils.SinkKafkaTagName:
			switch sinkCfg.Kafka.MessageMode {
			case utils.TransformJsonMode:
				resultData[sinkName] = []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
			case utils.TransformRowMode:
				var row string
				for _, item := range result {
					row += CastTypes(item, "toString").(string)
				}
				resultData[sinkName] = []any{[]byte(row)}
			}
		case utils.SinkRocketMQTagName:
			switch sinkCfg.RocketMQ.MessageMode {
			case utils.TransformJsonMode:
				resultData[sinkName] = []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
			case utils.TransformRowMode:
				var row string
				for _, item := range result {
					row += CastTypes(item, "toString").(string)
				}
				resultData[sinkName] = []any{[]byte(row)}
			}
		case utils.SinkRabbitMQTagName:
			switch sinkCfg.RabbitMQ.MessageMode {
			case utils.TransformJsonMode:
				resultData[sinkName] = []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
			case utils.TransformRowMode:
				var row string
				for _, item := range result {
					row += CastTypes(item, "toString").(string)
				}
				resultData[sinkName] = []any{[]byte(row)}
			}
		case utils.SinkPulsarTagName:
			switch sinkCfg.Pulsar.MessageMode {
			case utils.TransformJsonMode:
				resultData[sinkName] = []any{MapToJson(utils.StringSliceToMap(specialProcessKeys, result))}
			case utils.TransformRowMode:
				var row string
				for _, item := range result {
					row += CastTypes(item, "toString").(string)
				}
				resultData[sinkName] = []any{[]byte(row)}
			}
		}
	}
}

// JsonMode JSON mode
func JsonMode(beforeConvertData *models.TransformBeforeConvert, runningCfg []config.TransformSchema, streamConfig *config.StreamConfig) *models.TransformAfterConvert {
	// init
	specialProcessKeys := make([]string, 0)
	resultData := make(map[string][]any)
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
