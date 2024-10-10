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
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRedisSinkHandler(t *testing.T) {
	t.Helper()
	// init logger
	initLogger()
	// Sink 配置
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "redis-1",
		ChanSize:      100,
	}
	testSinkConfig := &config.SinkConfig{
		Type: "Redis",
		Redis: config.RedisSinkConfig{
			DBNum:            0,
			Address:          "localhost:6379",
			DataType:         "kv",
			KeyOrChannelName: "test-key",
		},
	}
	// init RedisSinkHandler
	redisClient := NewRedisSinkHandler(base, testSinkConfig.Redis)
	// Sink Write
	go redisClient.WriteData()
	// Channel
	r := redisClient.GetFromTransformChan()
	for i := 1; i < 5; i++ {
		r <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "redis-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	redisClient.CloseSink()
}
