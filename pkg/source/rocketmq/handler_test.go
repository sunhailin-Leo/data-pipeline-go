package source

import (
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const devRocketMQHosts string = "<test address>"

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRocketMQSource(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Source - Consumer
	testTopic := "alg_test_topic"
	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "rocketmq-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceRocketMQTagName,
			SourceName: "rocketmq-1",
			RocketMQ: config.RocketMQSourceConfig{
				Address:   devRocketMQHosts,
				Group:     utils.LogServiceName,
				Topic:     testTopic,
				AccessKey: "<test access key>",
				SecretKey: "<test secret key>",
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	rmq := NewRocketMQSource(baseSource)
	// rmq.SetDebugMode(true)
	c := rmq.GetToTransformChan()

	// Consumer - Fetch Data
	go rmq.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from RocketMQ failed")
	}

	rmq.CloseSource()
}
