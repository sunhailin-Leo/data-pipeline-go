package sink

import (
	"testing"
	"time"

	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewPulsarSinkHandler(t *testing.T) {
	t.Helper()
	// 初始化日志
	initLogger()
	// Sink 配置
	base := sink.BaseSink{
		DebugMode:     false,
		StreamName:    "",
		SinkAliasName: "pulsar-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamConfig:  &config.StreamConfig{},
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkPulsarTagName,
		SinkName: "pulsar-1",
		Pulsar: config.PulsarSinkConfig{
			Address: "172.20.49.19:16650",
			Topic:   "alg_test",
		},
	}
	// 初始化
	pulsarClient := NewPulsarSinkHandler(base, testSinkConfig.Pulsar)
	// Sink Write
	go pulsarClient.WriteData()
	// Channel
	p := pulsarClient.GetFromTransformChan()
	for i := 0; i < 10; i++ {
		p <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "pulsar-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	pulsarClient.CloseSink()
}
