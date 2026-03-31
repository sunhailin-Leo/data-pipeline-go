package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestLoadConfigFromLocal(t *testing.T) {
	initLogger()

	// Reset global variable before test
	TunnelCfg = nil

	// Create temporary config file
	testConfigFile := "test_config.json"
	testConfigContent := `{
		"streams": [
			{
				"name": "stream-1",
				"enable": true,
				"channel_size": 1000,
				"source": [
					{
						"type": "Kafka",
						"source_name": "kafka-1",
						"kafka": {
							"address": "localhost:9092",
							"group": "test-group",
							"topic": "test-topic"
						}
					}
				],
				"transform": {
					"mode": "json",
					"schemas": [
						{
							"source_key": "col1",
							"sink_key": "col1",
							"converter": "toInt32",
							"is_ignore": false,
							"is_strict_mode": true,
							"source_name": "kafka-1",
							"sink_name": "clickhouse-1"
						}
					]
				},
				"sink": [
					{
						"type": "ClickHouse",
						"sink_name": "clickhouse-1",
						"clickhouse": {
							"address": "localhost:9000",
							"username": "default",
							"password": "",
							"database": "test_db",
							"bulk_size": 100,
							"table_name": "test_table",
							"columns": [
								{"name": "col1", "type": "Int32"}
							]
						}
					}
				]
			},
			{
				"name": "stream-2",
				"enable": true,
				"channel_size": 500,
				"source": [
					{
						"type": "Kafka",
						"source_name": "kafka-2",
						"kafka": {
							"address": "localhost:9092",
							"group": "test-group-2",
							"topic": "test-topic-2"
						}
					}
				],
				"transform": {
					"mode": "json",
					"schemas": []
				},
				"sink": [
					{
						"type": "ClickHouse",
						"sink_name": "clickhouse-2",
						"clickhouse": {
							"address": "localhost:9000",
							"username": "default",
							"password": "",
							"database": "test_db",
							"bulk_size": 50,
							"table_name": "test_table_2",
							"columns": []
						}
					}
				]
			}
		]
	}`

	// Write temporary config file
	err := os.WriteFile(testConfigFile, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigFile)

	// Set environment variables
	_ = os.Setenv(utils.ConfigFromSourceName, "local")
	_ = os.Setenv(utils.ConfigFromLocalPathEnvName, testConfigFile)

	// Load config
	loader := &TunnelConfigLoader{}
	loader.loadFromLocal(testConfigFile)

	// Verify TunnelCfg is not nil
	assert.NotNil(t, TunnelCfg)

	// Verify streams count
	assert.Equal(t, 2, len(TunnelCfg.Streams))

	// Verify stream names
	assert.Equal(t, "stream-1", TunnelCfg.Streams[0].Name)
	assert.Equal(t, "stream-2", TunnelCfg.Streams[1].Name)

	// Reset global variable after test
	TunnelCfg = nil
}

func TestValidateConfig(t *testing.T) {
	initLogger()

	// Test normal config - should pass validation
	t.Run("Normal config should pass validation", func(t *testing.T) {
		loader := &TunnelConfigLoader{}
		normalConfig := &TunnelConfig{
			Streams: []*StreamConfig{
				{
					Name:   "stream-1",
					Enable: true,
					Source: []*SourceConfig{
						{
							Type:       "Kafka",
							SourceName: "kafka-1",
						},
					},
					Transform: TransformConfig{Mode: "json"},
					Sink: []*SinkConfig{
						{
							Type:     "ClickHouse",
							SinkName: "clickhouse-1",
						},
					},
				},
				{
					Name:   "stream-2",
					Enable: true,
					Source: []*SourceConfig{
						{
							Type:       "Kafka",
							SourceName: "kafka-2",
						},
					},
					Transform: TransformConfig{Mode: "json"},
					Sink: []*SinkConfig{
						{
							Type:     "ClickHouse",
							SinkName: "clickhouse-2",
						},
					},
				},
			},
		}

		err := loader.validateConfig(normalConfig)
		assert.NoError(t, err)
	})

	// Test duplicate stream name - should return error
	t.Run("Duplicate stream name should return error", func(t *testing.T) {
		loader := &TunnelConfigLoader{}
		duplicateConfig := &TunnelConfig{
			Streams: []*StreamConfig{
				{
					Name:   "stream-1",
					Enable: true,
					Source: []*SourceConfig{
						{
							Type:       "Kafka",
							SourceName: "kafka-1",
						},
					},
					Transform: TransformConfig{Mode: "json"},
					Sink: []*SinkConfig{
						{
							Type:     "ClickHouse",
							SinkName: "clickhouse-1",
						},
					},
				},
				{
					Name:   "stream-1",
					Enable: true,
					Source: []*SourceConfig{
						{
							Type:       "Kafka",
							SourceName: "kafka-2",
						},
					},
					Transform: TransformConfig{Mode: "json"},
					Sink: []*SinkConfig{
						{
							Type:     "ClickHouse",
							SinkName: "clickhouse-2",
						},
					},
				},
			},
		}

		err := loader.validateConfig(duplicateConfig)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stream name stream-1 is duplicated")
	})
}

func TestLoadFile(t *testing.T) {
	initLogger()

	// Reset global variable before test
	TunnelCfg = nil

	loader := &TunnelConfigLoader{}
	validJSON := `{
		"streams": [
			{
				"name": "test-stream",
				"enable": true,
				"source": [
					{
						"type": "Kafka",
						"source_name": "kafka-1",
						"kafka": {
							"address": "localhost:9092",
							"group": "test-group",
							"topic": "test-topic"
						}
					}
				],
				"transform": {
					"mode": "json",
					"schemas": []
				},
				"sink": [
					{
						"type": "ClickHouse",
						"sink_name": "clickhouse-1",
						"clickhouse": {
							"address": "localhost:9000",
							"username": "default",
							"password": "",
							"database": "test_db",
							"bulk_size": 100,
							"table_name": "test_table",
							"columns": []
						}
					}
				]
			}
		]
	}`

	// Note: loadFile calls logger.Fatal on error, which will exit the process
	// In a real test environment, you might want to mock the logger or use a different approach
	// For this test, we'll just verify the method exists and can be called
	// If the JSON is valid, it should parse correctly
	result := loader.loadFile([]byte(validJSON))
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Streams))
	assert.Equal(t, "test-stream", result.Streams[0].Name)

	// Reset global variable after test
	TunnelCfg = nil
}

func TestIsConfigLoaded(t *testing.T) {
	initLogger()

	// Reset global variable before test
	TunnelCfg = nil

	t.Run("Config not loaded should return false", func(t *testing.T) {
		loader := &TunnelConfigLoader{}
		assert.False(t, loader.IsConfigLoaded())
	})

	t.Run("Config loaded should return true", func(t *testing.T) {
		loader := &TunnelConfigLoader{}
		loader.config = &TunnelConfig{
			Streams: []*StreamConfig{
				{
					Name:   "test-stream",
					Enable: true,
					Source: []*SourceConfig{
						{
							Type:       "Kafka",
							SourceName: "kafka-1",
						},
					},
					Sink: []*SinkConfig{
						{
							Type:     "ClickHouse",
							SinkName: "clickhouse-1",
						},
					},
				},
			},
		}
		assert.True(t, loader.IsConfigLoaded())
	})
}

// TestBindAllConfigEnv tests that bindAllConfigEnv does not panic
func TestBindAllConfigEnv(t *testing.T) {
	loader := &TunnelConfigLoader{}
	assert.NotPanics(t, func() {
		loader.bindAllConfigEnv()
	})
}

func BenchmarkValidateConfig(b *testing.B) {
	loader := &TunnelConfigLoader{}
	normalConfig := &TunnelConfig{
		Streams: []*StreamConfig{
			{
				Name:   "stream-1",
				Enable: true,
				Source: []*SourceConfig{
					{
						Type:       "Kafka",
						SourceName: "kafka-1",
					},
				},
				Transform: TransformConfig{Mode: "json"},
				Sink: []*SinkConfig{
					{
						Type:     "ClickHouse",
						SinkName: "clickhouse-1",
					},
				},
			},
			{
				Name:   "stream-2",
				Enable: true,
				Source: []*SourceConfig{
					{
						Type:       "Kafka",
						SourceName: "kafka-2",
					},
				},
				Transform: TransformConfig{Mode: "json"},
				Sink: []*SinkConfig{
					{
						Type:     "ClickHouse",
						SinkName: "clickhouse-2",
					},
				},
			},
		},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loader.validateConfig(normalConfig)
	}
}

// TestLoad_LocalPath tests loading config from local file via environment variables
func TestLoad_LocalPath(t *testing.T) {
	initLogger()
	TunnelCfg = nil

	// Create temp config file
	testConfigFile := "test_load_config.json"
	testConfigContent := `{"streams":[{"name":"s1","enable":true,"source":[{"type":"Kafka","source_name":"k1"}],"transform":{"mode":"json","schemas":[]},"sink":[{"type":"ClickHouse","sink_name":"c1"}]}]}`
	err := os.WriteFile(testConfigFile, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigFile)

	// Set env vars
	os.Setenv("CONFIG_SRC", "local")
	os.Setenv("LOCAL_PATH", testConfigFile)
	defer os.Unsetenv("CONFIG_SRC")
	defer os.Unsetenv("LOCAL_PATH")

	loader := &TunnelConfigLoader{}
	loader.Load()

	assert.NotNil(t, TunnelCfg)
	assert.True(t, loader.IsConfigLoaded())

	TunnelCfg = nil
}

// TestValidateConfig_VdError tests validation with vd error
func TestValidateConfig_VdError(t *testing.T) {
	loader := &TunnelConfigLoader{}
	// Config with missing required fields (empty stream name, etc.)
	invalidConfig := &TunnelConfig{
		Streams: []*StreamConfig{
			{
				Name:   "", // empty name should fail vd validation
				Enable: true,
			},
		},
	}
	err := loader.validateConfig(invalidConfig)
	assert.Error(t, err)
}
