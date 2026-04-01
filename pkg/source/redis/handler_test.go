package source

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

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

func TestNewRedisSourceHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	redisAddr := testutil.GetEnvOrDefault(testutil.EnvRedisAddr, "localhost:6379")
	channelName := "integration-test-channel"

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
				KeyOrChannelName: channelName,
				Address:          redisAddr,
				DataType:         utils.RedisDataTypeSubscribe,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	r := NewRedisSourceHandler(baseSource)
	c := r.GetToTransformChan()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.FetchData()
	}()

	// Wait for subscription to be established
	time.Sleep(2 * time.Second)

	// Publisher
	publisher := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: strings.Split(redisAddr, ","),
	})
	defer publisher.Close()

	for i := 0; i < 10; i++ {
		result := publisher.Publish(context.Background(), channelName, "test-msg-"+strconv.Itoa(i))
		if result.Err() != nil {
			t.Fatalf("Failed to publish message: %v", result.Err())
		}
	}

	// Wait a bit for message to be delivered
	time.Sleep(500 * time.Millisecond)

	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from Redis failed")
	}

	fetchMessage := fetchData.SourceData.(string)
	t.Logf("Received message: %s", fetchMessage)

	// Drain remaining messages from channel to prevent FetchData from blocking on send
	done := make(chan struct{})
	go func() {
		for range c {
		}
		close(done)
	}()

	// CloseSource closes pubsub (which exits FetchData's for-range loop) then closes channel
	r.CloseSource()
	// Wait for drain goroutine to finish
	<-done
	// Wait for FetchData goroutine to fully exit
	wg.Wait()
}
