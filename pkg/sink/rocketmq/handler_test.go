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

func TestNewRocketMQSinkHandler(t *testing.T) {
	t.Helper()
	// init logger
	initLogger()
	// Sink config
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "rocketmq-1",
		ChanSize:      100,
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkRocketMQTagName,
		SinkName: "rocketmq-1",
		RocketMQ: config.RocketMQSinkConfig{
			Address:     "<test address>",
			Topic:       "<test topic>",
			MessageMode: "json",
		},
	}
	// init RocketMQSinkHandler
	rmqClient := NewRocketMQSinkHandler(base, testSinkConfig.RocketMQ)
	// Sink Write
	go rmqClient.WriteData()
	// Channel
	r := rmqClient.GetFromTransformChan()
	for i := 0; i < 5; i++ {
		r <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "rocketmq-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	rmqClient.CloseSink()
}
