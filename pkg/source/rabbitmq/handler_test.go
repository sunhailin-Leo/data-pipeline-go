package source

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/wagslane/go-rabbitmq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRabbitMQSource(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	rmqAddr := testutil.GetEnvOrDefault(testutil.EnvRabbitMQAddr, "localhost:5672")
	rmqUser := testutil.GetEnvOrDefault(testutil.EnvRabbitMQUser, "testuser")
	rmqPass := testutil.GetEnvOrDefault(testutil.EnvRabbitMQPass, "testpass")

	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "rabbitmq-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceRabbitMQTagName,
			SourceName: "rabbitmq-1",
			RabbitMQ: config.RabbitMQSourceConfig{
				Address:    rmqAddr,
				Username:   rmqUser,
				Password:   rmqPass,
				VHost:      "/",
				Queue:      "integration-test-source-queue",
				Exchange:   "integration-test-source-exchange",
				RoutingKey: "test.routing.key",
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	rmq := NewRabbitMQSource(baseSource)
	c := rmq.GetToTransformChan()

	// Create publisher connection first to ensure exchange and queue are properly set up
	publisherConn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s/",
			url.QueryEscape(rmqUser),
			url.QueryEscape(rmqPass),
			rmqAddr),
		rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		t.Fatalf("Failed to create publisher connection: %v", err)
	}
	defer publisherConn.Close()

	// Declare exchange and queue binding to ensure proper routing
	publisher, err := rabbitmq.NewPublisher(
		publisherConn,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		t.Fatalf("Failed to create publisher: %v", err)
	}
	defer publisher.Close()

	// Start consumer
	go rmq.FetchData()

	// Wait for consumer to be ready
	time.Sleep(2 * time.Second)

	// Publish message to the exchange
	publishErr := publisher.Publish(
		[]byte("test-message"),
		[]string{"test.routing.key"},
		rabbitmq.WithPublishOptionsExchange("integration-test-source-exchange"),
	)
	if publishErr != nil {
		t.Fatalf("Failed to publish message: %v", publishErr)
	}

	// Wait for message with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case fetchData, ok := <-c:
		if !ok || fetchData == nil {
			t.Fatalf("Fetch data from RabbitMQ failed")
		}
	case <-ctx.Done():
		t.Fatalf("Test timed out waiting for RabbitMQ message")
	}

	rmq.CloseSource()
}
