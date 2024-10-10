package sink

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// SetKV writes key-value data to Redis
func SetKV(client redis.UniversalClient, loggerName string, key string, value any, ex int) bool {
	keyTTL := time.Duration(-1)
	if ex > 0 {
		keyTTL = time.Duration(ex) * time.Second
	}

	_, setErr := client.Set(context.Background(), key, value, keyTTL).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}

// HashSet writes hash data to Redis
func HashSet(client redis.UniversalClient, loggerName string, key string, values []any) bool {
	_, setErr := client.HSet(context.Background(), key, values...).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}

// ListLPush writes list data to Redis
func ListLPush(client redis.UniversalClient, loggerName string, key string, values []any) bool {
	_, setErr := client.LPush(context.Background(), key, values...).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}

// ListRPush writes list data to Redis
func ListRPush(client redis.UniversalClient, loggerName string, key string, values []any) bool {
	_, setErr := client.RPush(context.Background(), key, values...).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}

// SetSAdd writes set data to Redis
func SetSAdd(client redis.UniversalClient, loggerName string, key string, values []any) bool {
	_, setErr := client.SAdd(context.Background(), key, values...).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}

// Publish writes message to Redis channel
func Publish(client redis.UniversalClient, loggerName string, channel string, message any) bool {
	_, setErr := client.Publish(context.Background(), channel, message).Result()
	if setErr != nil && !errors.Is(setErr, NilType) {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + loggerName + "]Failed to write data! Reason for error: " + setErr.Error())
		return false
	}
	return true
}
