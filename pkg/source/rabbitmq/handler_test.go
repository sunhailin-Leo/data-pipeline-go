package source

import (
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRabbitMQSource(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Source - Consumer
	testQueue := "alg-test.tunnel"
	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "rabbitmq-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceRabbitMQTagName,
			SourceName: "rabbitmq-1",
			RabbitMQ: config.RabbitMQSourceConfig{
				Address:    "<test address>",
				Username:   "<test username>",
				Password:   "<test password>",
				VHost:      "/",
				Queue:      testQueue,
				Exchange:   "alg-test",
				RoutingKey: "#",
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	rmq := NewRabbitMQSource(baseSource)
	// rmq.SetDebugMode(true)
	c := rmq.GetToTransformChan()

	// Consumer - Fetch Data
	go rmq.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from RabbitMQ failed")
	}

	rmq.CloseSource()
}
