package convert

import (
	"encoding/json"
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
		// compare as parsed JSON to avoid map iteration order issues
		var actualMap, expectedMap map[string]any
		_ = json.Unmarshal(testUnit.inputData["sink1"][0].([]byte), &actualMap)
		_ = json.Unmarshal([]byte(testUnit.expected.(string)), &expectedMap)
		assert.Equal(t, expectedMap, actualMap)
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
		ExpectNil  bool
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
			ExpectNil: false,
		},
		{
			Name: "test-strict-mode-missing-field",
			Config: config.TransformConfig{
				Mode: "json",
				Schemas: []config.TransformSchema{
					{
						SourceKey:    "requiredField",
						SinkKey:      "requiredField",
						Converter:    "",
						IsIgnore:     false,
						IsStrictMode: true,
						SourceName:   "Kafka-1",
						SinkName:     "Clickhouse-1",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput: &models.SourceOutput{},
				BeforeConvertData: map[string]any{
					"otherField": "value",
				},
			},
			OutputData: &models.TransformAfterConvert{},
			ExpectNil:  true,
		},
		{
			Name: "test-ignore-field",
			Config: config.TransformConfig{
				Mode: "json",
				Schemas: []config.TransformSchema{
					{
						SourceKey:  "field1",
						SinkKey:    "field1",
						Converter:  "",
						IsIgnore:   false,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1",
					},
					{
						SourceKey:  "field2",
						SinkKey:    "field2",
						Converter:  "",
						IsIgnore:   true,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput: &models.SourceOutput{},
				BeforeConvertData: map[string]any{
					"field1": "value1",
					"field2": "value2",
				},
			},
			OutputData: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {"value1"},
				},
			},
			ExpectNil: false,
		},
		{
			Name: "test-expand-field-with-uuid",
			Config: config.TransformConfig{
				Mode: "json",
				Schemas: []config.TransformSchema{
					{
						SourceKey:  "field1",
						SinkKey:    "field1",
						Converter:  "",
						IsIgnore:   false,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1",
					},
					{
						SourceKey:   "field2",
						SinkKey:     "field2",
						Converter:   "",
						IsIgnore:    false,
						IsExpand:    true,
						ExpandValue: "$.UUID()",
						SourceName:  "Kafka-1",
						SinkName:    "Clickhouse-1",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput: &models.SourceOutput{},
				BeforeConvertData: map[string]any{
					"field1": "value1",
				},
			},
			OutputData: nil, // UUID is dynamically generated, verified separately below
			ExpectNil:  false,
		},
		{
			Name: "test-expand-field-with-existing-value",
			Config: config.TransformConfig{
				Mode: "json",
				Schemas: []config.TransformSchema{
					{
						SourceKey:  "field1",
						SinkKey:    "field1",
						Converter:  "",
						IsIgnore:   false,
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1",
					},
					{
						SourceKey:   "field2",
						SinkKey:     "field2",
						Converter:   "",
						IsIgnore:    false,
						IsExpand:    true,
						ExpandValue: "$.UUID()",
						SourceName:  "Kafka-1",
						SinkName:    "Clickhouse-1",
					},
				},
			},
			InputData: &models.TransformBeforeConvert{
				SourceOutput: &models.SourceOutput{},
				BeforeConvertData: map[string]any{
					"field1": "value1",
					"field2": "existingValue",
				},
			},
			OutputData: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {"value1", "existingValue"},
				},
			},
			ExpectNil: false,
		},
	}

	for _, unit := range testUnits {
		actual := JsonMode(unit.InputData, unit.Config.Schemas, &config.StreamConfig{})
		if unit.ExpectNil {
			if actual != nil {
				t.Fatalf("TestJsonMode failed. TestName: %s, expected nil but got %v", unit.Name, actual)
			}
		} else {
			if actual == nil {
				t.Fatalf("TestJsonMode failed. TestName: %s, actual is nil", unit.Name)
			}
			if unit.OutputData == nil {
				// dynamic output (e.g. UUID expand), just verify result is not empty
				for sinkName, data := range actual.AfterConvertData {
					if len(data) == 0 {
						t.Fatalf("TestJsonMode failed. TestName: %s, sink %s has empty data", unit.Name, sinkName)
					}
				}
			} else if !reflect.DeepEqual(actual.AfterConvertData, unit.OutputData.AfterConvertData) {
				t.Fatalf("TestJsonMode failed. TestName: %s, Expected: %v, actual: %v", unit.Name, unit.OutputData.AfterConvertData, actual.AfterConvertData)
			}
		}
	}
}

// Benchmark functions
func BenchmarkJsonMode(b *testing.B) {
	initLogger()
	b.ReportAllocs()

	inputData := &models.TransformBeforeConvert{
		SourceOutput: &models.SourceOutput{},
		BeforeConvertData: map[string]any{
			"uniqueSequence": "123456",
			"data":           "123",
		},
	}

	schemas := []config.TransformSchema{
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
	}

	streamConfig := &config.StreamConfig{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JsonMode(inputData, schemas, streamConfig)
	}
}

func TestConvertToRowBytes(t *testing.T) {
	initLogger()
	result := []any{"hello", 123, "world"}
	actual := convertToRowBytes(result)
	assert.Len(t, actual, 1)
	assert.Equal(t, []byte("hello123world"), actual[0])
}

func TestProcessByMessageMode(t *testing.T) {
	initLogger()

	// Test JSON mode
	resultData := make(map[string][]any)
	result := []any{"value1", "value2"}
	specialProcessKeys := []string{"key1", "key2"}
	processByMessageMode("json", resultData, "sink1", specialProcessKeys, result)
	assert.Len(t, resultData["sink1"], 1)
	// Verify it's JSON bytes
	var jsonMap map[string]any
	err := json.Unmarshal(resultData["sink1"][0].([]byte), &jsonMap)
	assert.NoError(t, err)
	assert.Equal(t, "value1", jsonMap["key1"])
	assert.Equal(t, "value2", jsonMap["key2"])

	// Test row mode
	resultData = make(map[string][]any)
	processByMessageMode("row", resultData, "sink1", specialProcessKeys, result)
	assert.Len(t, resultData["sink1"], 1)
	assert.Equal(t, []byte("value1value2"), resultData["sink1"][0])
}

func TestGetMessageMode(t *testing.T) {
	// Test Kafka sink
	sinkCfg := &config.SinkConfig{
		Type:     "Kafka",
		SinkName: "kafka-sink",
		Kafka:    config.KafkaSinkConfig{MessageMode: "json"},
	}
	result := getMessageMode(sinkCfg)
	assert.Equal(t, "json", result)

	// Test RocketMQ sink
	sinkCfg = &config.SinkConfig{
		Type:     "RocketMQ",
		SinkName: "rocketmq-sink",
		RocketMQ: config.RocketMQSinkConfig{MessageMode: "row"},
	}
	result = getMessageMode(sinkCfg)
	assert.Equal(t, "row", result)

	// Test RabbitMQ sink
	sinkCfg = &config.SinkConfig{
		Type:     "RabbitMQ",
		SinkName: "rabbitmq-sink",
		RabbitMQ: config.RabbitMQSinkConfig{MessageMode: "json"},
	}
	result = getMessageMode(sinkCfg)
	assert.Equal(t, "json", result)

	// Test Pulsar sink
	sinkCfg = &config.SinkConfig{
		Type:     "Pulsar",
		SinkName: "pulsar-sink",
		Pulsar:   config.PulsarSinkConfig{MessageMode: "row"},
	}
	result = getMessageMode(sinkCfg)
	assert.Equal(t, "row", result)

	// Test default (non-MQ sink)
	sinkCfg = &config.SinkConfig{
		Type:     "ClickHouse",
		SinkName: "clickhouse-sink",
	}
	result = getMessageMode(sinkCfg)
	assert.Equal(t, "", result)
}

func TestProcessBySink_MQ(t *testing.T) {
	initLogger()

	resultData := map[string][]any{
		"kafka-sink": {"value1", "value2"},
	}

	streamConfig := &config.StreamConfig{
		Sink: []*config.SinkConfig{
			{
				Type:     "Kafka",
				SinkName: "kafka-sink",
				Kafka:    config.KafkaSinkConfig{MessageMode: "json"},
			},
		},
	}

	specialProcessKeys := []string{"key1", "key2"}
	processBySink(resultData, streamConfig, specialProcessKeys)

	// Verify result is converted to JSON bytes
	assert.Len(t, resultData["kafka-sink"], 1)
	var jsonMap map[string]any
	err := json.Unmarshal(resultData["kafka-sink"][0].([]byte), &jsonMap)
	assert.NoError(t, err)
	assert.Equal(t, "value1", jsonMap["key1"])
	assert.Equal(t, "value2", jsonMap["key2"])
}

func TestJsonMode_IsKeepKeys(t *testing.T) {
	initLogger()

	inputData := &models.TransformBeforeConvert{
		SourceOutput: &models.SourceOutput{},
		BeforeConvertData: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	schemas := []config.TransformSchema{
		{
			SourceKey:  "key1",
			SinkKey:    "key1",
			Converter:  "",
			IsIgnore:   false,
			IsKeepKeys: true,
			SourceName: "Kafka-1",
			SinkName:   "kafka-sink",
		},
		{
			SourceKey:  "key2",
			SinkKey:    "key2",
			Converter:  "",
			IsIgnore:   false,
			IsKeepKeys: true,
			SourceName: "Kafka-1",
			SinkName:   "kafka-sink",
		},
	}

	streamConfig := &config.StreamConfig{
		Sink: []*config.SinkConfig{
			{
				Type:     "Kafka",
				SinkName: "kafka-sink",
				Kafka:    config.KafkaSinkConfig{MessageMode: "json"},
			},
		},
	}

	result := JsonMode(inputData, schemas, streamConfig)
	assert.NotNil(t, result)
	assert.Len(t, result.AfterConvertData["kafka-sink"], 1)

	// Verify it's JSON bytes with the kept keys
	var jsonMap map[string]any
	err := json.Unmarshal(result.AfterConvertData["kafka-sink"][0].([]byte), &jsonMap)
	assert.NoError(t, err)
	assert.Equal(t, "value1", jsonMap["key1"])
	assert.Equal(t, "value2", jsonMap["key2"])
}

func TestJsonMode_IsStrictMode_Missing(t *testing.T) {
	initLogger()

	inputData := &models.TransformBeforeConvert{
		SourceOutput: &models.SourceOutput{},
		BeforeConvertData: map[string]any{
			"existingField": "value",
		},
	}

	schemas := []config.TransformSchema{
		{
			SourceKey:    "requiredField",
			SinkKey:      "requiredField",
			Converter:    "",
			IsIgnore:     false,
			IsStrictMode: true,
			SourceName:   "Kafka-1",
			SinkName:     "clickhouse-sink",
		},
	}

	streamConfig := &config.StreamConfig{}

	result := JsonMode(inputData, schemas, streamConfig)
	assert.Nil(t, result)
}
