package transform

import (
	"reflect"
	"testing"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/stretchr/testify/assert"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wagslane/go-rabbitmq"
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

	// create a minimal stream config for json mode (JsonMode needs streamConfig.Sink for map preallocation)
	streamConfig := &config.StreamConfig{
		Name:   "test-stream",
		Enable: true,
		Source: []*config.SourceConfig{{Type: "Kafka", SourceName: "kafka-1"}},
		Sink: []*config.SinkConfig{
			{Type: "ClickHouse", SinkName: "Clickhouse-1"},
			{Type: "ClickHouse", SinkName: "Clickhouse-2"},
		},
	}

	for _, unit := range testUnits {
		handler.SetStreamConfig(streamConfig).InitTransform(&unit.Config, 0)
		actual := handler.ConvertModeSelector(unit.InputData)
		if !reflect.DeepEqual(actual, unit.OutputData) {
			t.Fatalf("assertion failed, unexpected: %v, expected: %v", actual, unit.OutputData)
		}
	}
}

// TestBaseTransform_SetMetricsHooks tests that SetMetricsHooks does not panic
func TestBaseTransform_SetMetricsHooks(t *testing.T) {
	baseTransform := &BaseTransform{}
	metrics := &middlewares.Metrics{}

	assert.NotPanics(t, func() {
		result := baseTransform.SetMetricsHooks(metrics)
		assert.NotNil(t, result)
		assert.Equal(t, metrics, baseTransform.metrics)
	})
}

// TestBaseTransform_Panics tests that BaseTransform placeholder methods panic
func TestBaseTransform_Panics(t *testing.T) {
	baseTransform := &BaseTransform{}

	t.Run("ConvertModeSelector should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			baseTransform.ConvertModeSelector(&models.TransformBeforeConvert{})
		})
	})

	t.Run("From should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			baseTransform.From()
		})
	})

	t.Run("Convert should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			baseTransform.Convert()
		})
	})

	t.Run("To should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			baseTransform.To()
		})
	})

	t.Run("CloseTransform should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			baseTransform.CloseTransform()
		})
	})
}

// TestNewTransformHandler tests that NewTransformHandler returns a non-nil handler
func TestNewTransformHandler(t *testing.T) {
	inputChan := make(chan *models.SourceOutput)
	outputChan := map[string]chan *models.TransformOutput{
		"sink1": make(chan *models.TransformOutput),
	}

	handler := NewTransformHandler(inputChan, outputChan)
	assert.NotNil(t, handler)

	// Use type assertion to access private fields
	handlerImpl, ok := handler.(*Handler)
	assert.True(t, ok)
	assert.Equal(t, inputChan, handlerImpl.inputChan)
	assert.Equal(t, outputChan, handlerImpl.outputChan)
}

// TestSourceDataSelector tests the sourceDataSelector method
func TestSourceDataSelector(t *testing.T) {
	initLogger()

	handler := NewTransformHandler(nil, nil)
	handlerImpl := handler.(*Handler)

	t.Run("Kafka source", func(t *testing.T) {
		kafkaRecord := &kgo.Record{Value: []byte("kafka-data")}
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Kafka"},
			SourceData: kafkaRecord,
		}
		result := handlerImpl.sourceDataSelector(sourceOutput)
		assert.Equal(t, []byte("kafka-data"), result)
	})

	t.Run("RocketMQ source", func(t *testing.T) {
		rocketMQRecord := &primitive.MessageExt{}
		rocketMQRecord.Body = []byte("rocketmq-data")
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "RocketMQ"},
			SourceData: rocketMQRecord,
		}
		result := handlerImpl.sourceDataSelector(sourceOutput)
		assert.Equal(t, []byte("rocketmq-data"), result)
	})

	t.Run("RabbitMQ source", func(t *testing.T) {
		rabbitMQRecord := &rabbitmq.Delivery{}
		rabbitMQRecord.Body = []byte("rabbitmq-data")
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "RabbitMQ"},
			SourceData: rabbitMQRecord,
		}
		result := handlerImpl.sourceDataSelector(sourceOutput)
		assert.Equal(t, []byte("rabbitmq-data"), result)
	})

	t.Run("PromMetrics source", func(t *testing.T) {
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "PromMetrics"},
			SourceData: []byte("prom-data"),
		}
		result := handlerImpl.sourceDataSelector(sourceOutput)
		assert.Equal(t, []byte("prom-data"), result)
	})

	t.Run("default source", func(t *testing.T) {
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: "raw-data",
		}
		result := handlerImpl.sourceDataSelector(sourceOutput)
		assert.Equal(t, "raw-data", result)
	})
}

// TestSourceTransformModeSelector tests the sourceTransformModeSelector method
func TestSourceTransformModeSelector(t *testing.T) {
	initLogger()

	handler := NewTransformHandler(nil, nil)
	handlerImpl := handler.(*Handler)

	t.Run("row mode", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{Mode: "row"}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: "test-data",
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.NotNil(t, result)
		assert.Equal(t, "test-data", result.BeforeConvertData)
	})

	t.Run("row mode with filters", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{
			Mode: "row",
			Filters: []config.TransformFilter{
				{Field: "key", Operator: "eq", Value: "value"},
			},
		}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: "test-data",
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.NotNil(t, result)
	})

	t.Run("json mode", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{Mode: "json"}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.NotNil(t, result)
		assert.IsType(t, map[string]any{}, result.BeforeConvertData)
	})

	t.Run("json mode with filter match", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{
			Mode: "json",
			Filters: []config.TransformFilter{
				{Field: "key", Operator: "eq", Value: "value"},
			},
		}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.NotNil(t, result)
	})

	t.Run("json mode with filter not match", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{
			Mode: "json",
			Filters: []config.TransformFilter{
				{Field: "key", Operator: "eq", Value: "other"},
			},
		}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.Nil(t, result)
	})

	t.Run("jsonPath mode", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{
			Mode: "jsonPath",
			Paths: []config.TransformJsonPath{
				{Path: "key", DestField: "key"},
			},
		}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.NotNil(t, result)
	})

	t.Run("jsonPath mode paths nil", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{
			Mode:  "jsonPath",
			Paths: nil,
		}, 1)
		_ = &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		// This test case is skipped because logger.Logger.Fatal will exit the program
		// The actual behavior is that the program will exit with a fatal error
		t.Skip("Skipping test that calls logger.Logger.Fatal")
	})

	t.Run("default mode", func(t *testing.T) {
		handlerImpl.InitTransform(&config.TransformConfig{Mode: "unknown"}, 1)
		sourceOutput := &models.SourceOutput{
			MetaData:   &models.MetaData{SourceTagName: "Unknown"},
			SourceData: []byte(`{"key":"value"}`),
		}
		result := handlerImpl.sourceTransformModeSelector(sourceOutput)
		assert.Nil(t, result)
	})
}

// TestHandler_InitTransform tests the InitTransform method
func TestHandler_InitTransform(t *testing.T) {
	initLogger()

	handler := NewTransformHandler(nil, nil)
	transformConfig := config.TransformConfig{
		Mode: "json",
		Schemas: []config.TransformSchema{
			{
				SourceKey:  "key",
				SinkKey:    "key",
				Converter:  "",
				SourceName: "kafka-1",
				SinkName:   "clickhouse-1",
			},
		},
	}

	result := handler.InitTransform(&transformConfig, 10)
	assert.NotNil(t, result)
}

// TestHandler_CloseTransform tests the CloseTransform method
func TestHandler_CloseTransform(t *testing.T) {
	initLogger()

	handler := NewTransformHandler(nil, nil)
	transformConfig := config.TransformConfig{Mode: "json"}
	handler.InitTransform(&transformConfig, 10)

	assert.NotPanics(t, func() {
		handler.CloseTransform()
	})
}

// TestHandler_Pipeline tests the complete From, Convert, To pipeline
func TestHandler_Pipeline(t *testing.T) {
	initLogger()

	inputChan := make(chan *models.SourceOutput, 1)
	sinkOutputChan := make(chan *models.TransformOutput, 1)
	outputChan := map[string]chan *models.TransformOutput{
		"Clickhouse-1": sinkOutputChan,
	}

	metrics := middlewares.NewMetrics("test")
	streamConfig := &config.StreamConfig{
		Name:   "test-stream",
		Enable: true,
		Source: []*config.SourceConfig{{Type: "Kafka", SourceName: "kafka-1"}},
		Sink:   []*config.SinkConfig{{Type: "ClickHouse", SinkName: "Clickhouse-1"}},
	}
	transformConfig := config.TransformConfig{
		Mode: "json",
		Schemas: []config.TransformSchema{
			{
				SourceKey:  "key1",
				SinkKey:    "key1",
				Converter:  "",
				SourceName: "kafka-1",
				SinkName:   "Clickhouse-1",
			},
		},
	}

	handler := NewTransformHandler(inputChan, outputChan)
	handler.SetStreamConfig(streamConfig).SetMetricsHooks(metrics).InitTransform(&transformConfig, 10)

	// Start pipeline goroutines
	go handler.(*Handler).From()
	go handler.(*Handler).Convert()
	go handler.(*Handler).To()

	// Send test data
	inputChan <- &models.SourceOutput{
		MetaData:   &models.MetaData{SourceTagName: "Unknown"},
		SourceData: []byte(`{"key1":"value1"}`),
	}

	// Receive output with timeout
	select {
	case output := <-sinkOutputChan:
		assert.NotNil(t, output)
		assert.Equal(t, "Clickhouse-1", output.SinkName)
		assert.Equal(t, []any{"value1"}, output.Data)
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for output")
	}

	// Close input to trigger exit
	close(inputChan)
}

// TestBaseTransform_Close tests the close method of BaseTransform
func TestBaseTransform_Close(t *testing.T) {
	bt := &BaseTransform{}
	bt.initChannel(1)
	assert.NotPanics(t, func() {
		bt.close()
	})
}
