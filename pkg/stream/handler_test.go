package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewSink(t *testing.T) {
	initLogger()

	metrics := middlewares.NewMetrics("data")

	t.Run("Console sink should return non-nil", func(t *testing.T) {
		baseSink := sink.BaseSink{
			Metrics:       metrics,
			StreamName:    "test-stream",
			SinkAliasName: "console-1",
			ChanSize:      100,
			StreamConfig:  &config.StreamConfig{},
		}
		sinkConfig := &config.SinkConfig{}

		sinker := NewSink(utils.SinkConsoleTagName, baseSink, sinkConfig)

		assert.NotNil(t, sinker, "Console sink should not be nil")
	})

	t.Run("Unknown sink type should return nil", func(t *testing.T) {
		baseSink := sink.BaseSink{
			Metrics:       metrics,
			StreamName:    "test-stream",
			SinkAliasName: "unknown-1",
			ChanSize:      100,
			StreamConfig:  &config.StreamConfig{},
		}
		sinkConfig := &config.SinkConfig{}

		sinker := NewSink("UnknownType", baseSink, sinkConfig)

		assert.Nil(t, sinker, "Unknown sink type should return nil")
	})
}

func TestNewSource(t *testing.T) {
	initLogger()

	metrics := middlewares.NewMetrics("data")

	t.Run("Unknown source type should return nil", func(t *testing.T) {
		baseSource := source.BaseSource{
			ChanSize:        100,
			AckMode:         utils.AckModeInSource,
			StreamName:      "test-stream",
			SourceAliasName: "unknown-1",
			SourceConfig:    &config.SourceConfig{},
			Metrics:         metrics,
		}

		src := NewSource("UnknownType", &baseSource)

		assert.Nil(t, src, "Unknown source type should return nil")
	})
}

func TestBaseStreamSetupWorkPool(t *testing.T) {
	initLogger()

	metrics := middlewares.NewMetrics("data")

	t.Run("setupWorkPool should initialize streamWorkPool", func(t *testing.T) {
		streamConfig := &config.StreamConfig{
			Name:        "test-stream",
			Enable:      true,
			ChannelSize: 100,
			Source:      []*config.SourceConfig{},
			Sink:        []*config.SinkConfig{},
		}

		baseStream := BaseStream{
			metrics:       metrics,
			streamsConfig: []*config.StreamConfig{streamConfig},
		}

		baseStream.setupWorkPool()

		assert.NotNil(t, baseStream.streamWorkPool, "streamWorkPool should not be nil after setupWorkPool")

		baseStream.releaseWorkPool()
	})
}
