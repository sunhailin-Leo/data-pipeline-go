package sink

import (
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewConsoleSinkHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink Console Test
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "console-1",
		ChanSize:      100,
	}
	consoleClient := NewConsoleSinkHandler(base)
	// Sink Write
	go consoleClient.WriteData()

	// channel
	c := consoleClient.GetFromTransformChan()
	for i := 1; i < 100; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				i,
				cast.ToFloat32(i),
				cast.ToString(i),
				"test",
			},
			SinkName: "console-1",
		}
	}

	// for waiting data insert
	time.Sleep(20 * time.Second)

	consoleClient.CloseSink()
}

func TestConsoleSinkName(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink Console Test
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "console-1",
		ChanSize:      100,
	}
	consoleClient := NewConsoleSinkHandler(base)

	assert.Equal(t, "Console", consoleClient.SinkName(), "SinkName should return 'Console'")
}

func TestConsoleSinkInitAndClose(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink Console Test
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "console-1",
		ChanSize:      100,
	}
	consoleClient := NewConsoleSinkHandler(base)

	// Verify InitSink does not panic
	assert.NotPanics(t, func() {
		consoleClient.InitSink()
	}, "InitSink should not panic")

	// Verify CloseSink does not panic
	assert.NotPanics(t, func() {
		consoleClient.CloseSink()
	}, "CloseSink should not panic")
}
