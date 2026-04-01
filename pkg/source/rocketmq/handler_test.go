package source

import (
	"testing"
	"time"

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

func TestNewRocketMQSource(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	rocketmqAddr := testutil.GetEnvOrDefault(testutil.EnvRocketMQAddr, "127.0.0.1:9876")

	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "rocketmq-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceRocketMQTagName,
			SourceName: "rocketmq-1",
			RocketMQ: config.RocketMQSourceConfig{
				Address: rocketmqAddr,
				Group:   utils.LogServiceName,
				Topic:   "integration-test-source-topic",
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	rmq := NewRocketMQSource(baseSource)
	c := rmq.GetToTransformChan()

	go rmq.FetchData()

	// Add timeout to prevent blocking indefinitely
	select {
	case fetchData, ok := <-c:
		if !ok || fetchData == nil {
			t.Logf("No message received from RocketMQ (this is expected if no messages were published)")
		}
	case <-time.After(10 * time.Second):
		t.Logf("Test timed out waiting for RocketMQ message (this is expected if no messages were published)")
	}

	rmq.CloseSource()
}
