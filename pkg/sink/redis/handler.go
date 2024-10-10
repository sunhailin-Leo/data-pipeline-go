package sink

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	NilType           = redis.Nil
	RedisPoolSize int = 200 // pool size
	RedisMinIdles int = 100 // min idles connection
)

type RedisSinkHandler struct {
	sink.BaseSink

	sinkRedisCfg config.RedisSinkConfig
	redisClient  redis.UniversalClient
}

// SinkName returns the name of the Redis sink
func (r *RedisSinkHandler) SinkName() string {
	return utils.SinkRedisTagName
}

// WriteData sink to Redis
func (r *RedisSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Redis-Sink][Current config: " + r.SinkAliasName + "]Start waiting for data to be written...")
	for {
		if r.redisClient == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Redis-Sink][Current config: " + r.SinkAliasName + "]Redis client connection closed or not connected!!")
			return
		}

		data, ok := <-r.GetFromTransformChan()
		r.Metrics.OnSinkInput(r.StreamName, r.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName +
				"[Redis-Sink][Current config: " + r.SinkAliasName + "]Sink is already closed!")
			return
		}
		r.Metrics.OnSinkInputSuccess(r.StreamName, r.SinkAliasName)

		var isSuccess bool
		key := r.sinkRedisCfg.KeyOrChannelName
		if r.sinkRedisCfg.KeyOrChannelName == utils.RedisDynamicKeyFlagName {
			var popKey any
			data.Data, popKey = utils.PopSlice(data.Data)
			if len(data.Data) == 0 {
				logger.Logger.Error(utils.LogServiceName + "[Redis-Sink][Current config: " + r.SinkAliasName + "]Data is empty!")
				continue
			}
			key = popKey.(string)
		}

		switch r.sinkRedisCfg.DataType {
		case utils.RedisDataTypeKV:
			isSuccess = SetKV(r.redisClient, r.SinkAliasName, key, data.Data[0], r.sinkRedisCfg.Expire)
		case utils.RedisDataTypeHash:
			isSuccess = HashSet(r.redisClient, r.SinkAliasName, key, data.Data)
		case utils.RedisDataTypeLPush:
			isSuccess = ListLPush(r.redisClient, r.SinkAliasName, key, data.Data)
		case utils.RedisDataTypeRPush:
			isSuccess = ListRPush(r.redisClient, r.SinkAliasName, key, data.Data)
		case utils.RedisDataTypeSet:
			isSuccess = SetSAdd(r.redisClient, r.SinkAliasName, key, data.Data)
		case utils.RedisDataTypePublish:
			isSuccess = Publish(r.redisClient, r.SinkAliasName, key, data.Data[0])
		default:
			logger.Logger.Error(utils.LogServiceName + "[Redis-Sink][Current config: " + r.SinkAliasName + "]Redis operation type error!")
		}
		r.Metrics.OnSinkOutput(r.StreamName, r.SinkAliasName)
		if isSuccess {
			r.Metrics.OnSinkOutputSuccess(r.StreamName, r.SinkAliasName)
			r.MessageCommit(data.SourceObj, data.SourceData, r.SinkAliasName)
		}
	}
}

// InitSink initializes the Redis sink
func (r *RedisSinkHandler) InitSink() {
	// single/cluster/proxy
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        strings.Split(r.sinkRedisCfg.Address, ","),
		Username:     r.sinkRedisCfg.Username,
		Password:     r.sinkRedisCfg.Password,
		DB:           r.sinkRedisCfg.DBNum,
		PoolSize:     RedisPoolSize,
		MinIdleConns: RedisMinIdles,
	})
	if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Sink][Current config: " + r.SinkAliasName + "]Failed to connect Redis! Reason for exception: " + pingErr.Error())
		return
	}
	r.redisClient = client

	logger.Logger.Info(utils.LogServiceName + "[Redis-Sink][Current config: " + r.SinkAliasName + "]Init Redis Successful!")
}

// CloseSink closes the Redis sink
func (r *RedisSinkHandler) CloseSink() {
	if r.redisClient != nil {
		_ = r.redisClient.Close()
	}
	r.Close()
}

// NewRedisSinkHandler creates a new Redis sink handler
func NewRedisSinkHandler(baseSink sink.BaseSink, sinkRedisCfg config.RedisSinkConfig) *RedisSinkHandler {
	handler := &RedisSinkHandler{BaseSink: baseSink, sinkRedisCfg: sinkRedisCfg}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
