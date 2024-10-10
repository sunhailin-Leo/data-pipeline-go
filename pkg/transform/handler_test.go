package transform

import (
	"reflect"
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestConvertModeSelector(t *testing.T) {
	t.Helper()

	initLogger()
	handler := NewTransformHandler(nil, nil)

	type testCaseUnit struct {
		Config     config.TransformConfig
		InputData  *models.TransformBeforeConvert
		OutputData *models.TransformAfterConvert
	}

	testUnits := []testCaseUnit{
		{
			Config: config.TransformConfig{
				Mode: "row",
				Schemas: []config.TransformSchema{
					{
						Converter:  "toInt",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput:      &models.SourceOutput{},
				BeforeConvertData: "123",
			},
			OutputData: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {123},
					"Clickhouse-2": {123},
				},
			},
		},
		{
			Config: config.TransformConfig{
				Mode: "row",
				Schemas: []config.TransformSchema{
					{
						Converter:  "toInt",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
					{
						Converter:  "toString",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
				},
				RowSeparator: ",",
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput:      &models.SourceOutput{},
				BeforeConvertData: "123,abc",
			},
			OutputData: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {123, "abc"},
					"Clickhouse-2": {123, "abc"},
				},
			},
		},
		{
			Config: config.TransformConfig{
				Mode: "json",
				Schemas: []config.TransformSchema{
					{
						SourceKey:  "uniqueSequence",
						SinkKey:    "uniqueSequence",
						Converter:  "",
						IsIgnore:   false,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
					{
						SourceKey:  "data",
						SinkKey:    "data",
						Converter:  "toInt",
						IsIgnore:   false,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput: &models.SourceOutput{},
				BeforeConvertData: map[string]any{
					"uniqueSequence": "123456",
					"data":           "123",
				},
			},
			OutputData: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {"123456", 123},
					"Clickhouse-2": {"123456", 123},
				},
			},
		},
	}

	for _, unit := range testUnits {
		handler.InitTransform(unit.Config, 0)
		actual := handler.ConvertModeSelector(unit.InputData)
		if !reflect.DeepEqual(actual, unit.OutputData) {
			t.Fatalf("assertion failed, unexpected: %v, expected: %v", actual, unit.OutputData)
		}
	}
}
