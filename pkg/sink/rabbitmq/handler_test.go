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
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRabbitMQSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	rmqAddr := testutil.GetEnvOrDefault(testutil.EnvRabbitMQAddr, "localhost:5672")
	rmqUser := testutil.GetEnvOrDefault(testutil.EnvRabbitMQUser, "testuser")
	rmqPass := testutil.GetEnvOrDefault(testutil.EnvRabbitMQPass, "testpass")

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
			Address:    rmqAddr,
			Username:   rmqUser,
			Password:   rmqPass,
			VHost:      "/",
			Queue:      "integration-test-queue",
			Exchange:   "integration-test-exchange",
			RoutingKey: "#",
		},
	}

	rmqClient := NewRabbitMQSinkHandler(base, &testSinkConfig.RabbitMQ)
	go rmqClient.WriteData()

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

	time.Sleep(3 * time.Second)
	rmqClient.CloseSink()
}
