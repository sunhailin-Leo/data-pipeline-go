package sink

import (
	"sync"
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

func TestNewRocketMQSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	rocketmqAddr := testutil.GetEnvOrDefault(testutil.EnvRocketMQAddr, "127.0.0.1:9876")

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
			Address:     rocketmqAddr,
			Topic:       "integration-test-topic",
			MessageMode: "json",
		},
	}

	rmqClient := NewRocketMQSinkHandler(base, testSinkConfig.RocketMQ)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rmqClient.WriteData()
	}()

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

	// Wait for all data to be processed, then close channel and wait for goroutine exit
	time.Sleep(2 * time.Second)
	rmqClient.Close()
	wg.Wait()
}
