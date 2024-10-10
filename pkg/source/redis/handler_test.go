package source

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func testRedisPublisherMock(channelName string) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split("<Redis address>", ","),
		Password: "<Redis password>",
	})

	for i := 0; i < 10; i++ {
		result := client.Publish(context.Background(), channelName, "test-msg-"+strconv.Itoa(i))
		if result.Err() != nil {
			panic(result.Err())
		} else {
			println(result.String())
		}
	}
}

func TestNewRedisSourceHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Source - Subscribe
	baseSource := source.BaseSource{
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "redis-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceRedisTagName,
			SourceName: "redis-1",
			Redis: config.RedisSourceConfig{
				DBNum:            0,
				KeyOrChannelName: "alg-test-redis",
				Address:          "<Redis address>",
				Password:         "<Redis password>",
				DataType:         utils.RedisDataTypeSubscribe,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	r := NewRedisSourceHandler(baseSource)
	// redis.SetDebugMode(true)
	c := r.GetToTransformChan()

	// Consumer - FetchData
	go r.FetchData()
	// Publisher
	testRedisPublisherMock("alg-test-redis")

	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from Redis failed")
	}

	// Parse Data
	fetchMessage := fetchData.SourceData.(string)
	println(fetchMessage)

	r.CloseSource()
}
