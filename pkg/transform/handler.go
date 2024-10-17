package transform

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wagslane/go-rabbitmq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/commiter"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/transform/convert"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type Handler struct {
	BaseTransform

	inputChan  chan *models.SourceOutput
	outputChan map[string]chan *models.TransformOutput
}

// sourceDataSelector source selector
func (t *Handler) sourceDataSelector(sourceData *models.SourceOutput) any {
	switch sourceData.SourceTagName {
	case utils.SourceKafkaTagName:
		kafkaRecord := sourceData.SourceData.(*kgo.Record)
		return kafkaRecord.Value
	case utils.SourceRocketMQTagName:
		rocketMQRecord := sourceData.SourceData.(*primitive.MessageExt)
		return rocketMQRecord.Body
	case utils.SinkRabbitMQTagName:
		rabbitMQRecord := sourceData.SourceData.(*rabbitmq.Delivery)
		return rabbitMQRecord.Body
	case utils.SourcePromMetricsTagName:
		return sourceData.SourceData.([]byte)
	default:
		return sourceData.SourceData
	}
}

// sourceTransformModeSelector source transform mode selector
func (t *Handler) sourceTransformModeSelector(sourceData *models.SourceOutput) *models.TransformBeforeConvert {
	switch t.configs.Mode {
	case utils.TransformRowMode:
		return &models.TransformBeforeConvert{
			SourceOutput:      sourceData,
			BeforeConvertData: t.sourceDataSelector(sourceData),
		}
	case utils.TransformJsonMode:
		return &models.TransformBeforeConvert{
			SourceOutput:      sourceData,
			BeforeConvertData: convert.JsonToMap(t.sourceDataSelector(sourceData).([]byte)),
		}
	default:
		logger.Logger.Error(utils.LogServiceName + "[Transform-From]unknown transform mode!")
		return nil
	}
}

// From read data from source
func (t *Handler) From() {
	for {
		// from source
		data, ok := <-t.inputChan
		if !ok || data == nil {
			logger.Logger.Error(utils.LogServiceName + "[Transform-From]Source data is nil or empty!")
			break
		}
		t.metrics.OnTransformInput(t.streamConfig.Name)
		// select transform mode
		beforeTransformData := t.sourceTransformModeSelector(data)
		if beforeTransformData == nil {
			continue
		}
		// transform data to Convert module
		t.beforeConvertChan <- beforeTransformData
		t.metrics.OnTransformInputSuccess(t.streamConfig.Name)
	}
}

// ConvertModeSelector convert mode selector
func (t *Handler) ConvertModeSelector(beforeConvertData *models.TransformBeforeConvert) *models.TransformAfterConvert {
	switch t.configs.Mode {
	case utils.TransformRowMode:
		return convert.RowMode(beforeConvertData, t.configs)
	case utils.TransformJsonMode:
		/*
			Json Mode:
				根据 ==> t.runningCfg:
					[
						{
							"sourceKey": "uniqueSequence",
							"sinkKey": "uniqueSequence",
							"converter": "",
							"isIgnore": false,
							"sourceName": "Kafka-1",
							"sinkName": "Clickhouse-1,Clickhouse-2"
						},
						{
							"sourceKey": "data",
							"sinkKey": "data",
							"converter": "toInt",
							"isIgnore": false,
							"sourceName": "Kafka-1",
							"sinkName": "Clickhouse-1,Clickhouse-2"
						}
					]
				输入 ==> beforeConvertData: {"uniqueSequence": "123456", "data": "123"}
				输出 ==> resultData: {
					"Clickhouse-1": []any{"123456", 123},
					"Clickhouse-2": []any{"123456", 123},
				}

			Row Mode:
				根据 ==> t.runningCfg:
					[
						{
							"converter": "toInt",
							"sourceName": "Kafka-1",
							"sinkName": "Clickhouse-1,Clickhouse-2"
						}
					]
				输入 ==> beforeConvertData: "123456"
				输出 ==> resultData: {
					"Clickhouse-1": []any{123},
					"Clickhouse-2": []any{123},
				}
		*/
		return convert.JsonMode(beforeConvertData, t.configs.Schemas, t.streamConfig)
	default:
		logger.Logger.Error(utils.LogServiceName + "[Transform-Convert]unknown transform mode!")
		return nil
	}
}

// Convert convert data
func (t *Handler) Convert() {
	for {
		// read data from [From] module
		data, ok := <-t.beforeConvertChan
		if !ok || data == nil {
			logger.Logger.Error(utils.LogServiceName + "[Transform-Convert]From data is nil or empty!")
			break
		}
		t.metrics.OnTransformConvert(t.streamConfig.Name)
		// select Convert type
		afterConvertData := t.ConvertModeSelector(data)
		if afterConvertData == nil {
			logger.Logger.Error(utils.LogServiceName + "[Transform-Convert]After convert data is nil!")
			continue
		}
		// transform data to [To] module
		t.afterConvertChan <- afterConvertData
		t.metrics.OnTransformConvertSuccess(t.streamConfig.Name)
	}
}

func (t *Handler) To() {
	for {
		// read data from [Convert] module
		data, ok := <-t.afterConvertChan
		if !ok || data == nil {
			logger.Logger.Error(utils.LogServiceName + "[Transform-To]Convert data is nil or empty!")
			break
		}
		t.metrics.OnTransformOutput(t.streamConfig.Name)
		// based on sink mapping to transform data to sink
		for sinkName, sinkData := range data.AfterConvertData {
			t.outputChan[sinkName] <- &models.TransformOutput{
				SourceOutput: data.SourceOutput,
				Data:         sinkData,
				SinkName:     sinkName,
			}
		}

		// commit
		if t.streamConfig.SourceAckMode == utils.AckModeInTransform {
			// Hard code for configName
			commiter.MessageCommit(data.SourceObj, data.SourceOutput, "transform")
		}

		t.metrics.OnTransformOutputSuccess(t.streamConfig.Name)
	}
}

// InitTransform init transform module
func (t *Handler) InitTransform(transformConfig config.TransformConfig, chanSize int) Transform {
	t.BaseTransform.InitTransform(transformConfig, chanSize)
	logger.Logger.Info(utils.LogServiceName + "Transform initialize successful!")
	return t
}

// CloseTransform close transform
func (t *Handler) CloseTransform() {
	t.close()
}

// NewTransformHandler create new transform handler
func NewTransformHandler(inputChan chan *models.SourceOutput, outputChan map[string]chan *models.TransformOutput) Transform {
	return &Handler{inputChan: inputChan, outputChan: outputChan}
}
