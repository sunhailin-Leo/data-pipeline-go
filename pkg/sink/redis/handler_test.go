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
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewRedisSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	redisAddr := testutil.GetEnvOrDefault(testutil.EnvRedisAddr, "localhost:6379")

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
			Address:          redisAddr,
			DataType:         "kv",
			KeyOrChannelName: "test-key",
		},
	}

	redisClient := NewRedisSinkHandler(base, testSinkConfig.Redis)
	go redisClient.WriteData()

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

	time.Sleep(3 * time.Second)
	redisClient.CloseSink()
}
