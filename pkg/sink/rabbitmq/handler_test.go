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

func TestNewRabbitMQSinkHandler(t *testing.T) {
	t.Helper()
	// init logger
	initLogger()
	// Sink config
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "rabbitmq-1",
		ChanSize:      100,
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkRabbitMQTagName,
		SinkName: "rabbitmq-1",
		RabbitMQ: config.RabbitMQSinkConfig{
			Address:    "<test address>",
			Username:   "<test username>",
			Password:   "<test password>",
			VHost:      "/",
			Queue:      "<test queue>",
			Exchange:   "<test exchange>",
			RoutingKey: "#",
		},
	}
	// init RabbitMQSinkHandler
	rmqClient := NewRabbitMQSinkHandler(base, testSinkConfig.RabbitMQ)
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
			SinkName: "rabbitmq-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	rmqClient.CloseSink()

}
