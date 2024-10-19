package convert

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
)

func TestProcessBySink(t *testing.T) {
	type testCaseUnit struct {
		inputData   map[string][]any
		config      *config.StreamConfig
		specialKeys []string
		expected    any
	}

	testUnits := []testCaseUnit{
		{
			inputData: map[string][]any{"sink1": {1, 2, 3}},
			config: &config.StreamConfig{
				Sink: []*config.SinkConfig{
					{
						Type:     "HTTP",
						SinkName: "sink1",
						HTTP:     config.HTTPSinkConfig{ContentType: "application/json"},
					},
				},
			},
			specialKeys: []string{"key1", "key2", "key3"},
			expected:    `{"key1":1,"key2":2,"key3":3}`,
		},
	}

	for _, testUnit := range testUnits {
		processBySink(testUnit.inputData, testUnit.config, testUnit.specialKeys)
		assert.Equal(t, testUnit.expected, string(testUnit.inputData["sink1"][0].([]byte)))
	}
}

func TestJsonMode(t *testing.T) {
	t.Helper()

	initLogger()

	type testCaseUnit struct {
		Name       string
		Config     config.TransformConfig
		InputData  *models.TransformBeforeConvert
		OutputData *models.TransformAfterConvert
	}

	testUnits := []testCaseUnit{
		{
			Name: "test-1",
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
		actual := JsonMode(unit.InputData, unit.Config.Schemas, &config.StreamConfig{})
		if actual == nil {
			t.Fatalf("TestJsonMode failed. TestName: %s, actual is nil", unit.Name)
		}
		if !reflect.DeepEqual(actual.AfterConvertData, unit.OutputData.AfterConvertData) {
			t.Fatalf("TestJsonMode failed. TestName: %s, Expected: %v, actual: %v", unit.Name, unit.OutputData.AfterConvertData, actual.AfterConvertData)
		}
	}
}
