package source

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	NilType           = redis.Nil
	RedisPoolSize int = 200 // pool size
	RedisMinIdles int = 100 // min pool idle connection
)

type RedisSourceHandler struct {
	source.BaseSource

	sourceRedisCfg config.RedisSourceConfig
	redisClient    redis.UniversalClient
	redisPubSub    *redis.PubSub
}

// SourceName returns the name of the Redis source
func (r *RedisSourceHandler) SourceName() string {
	return utils.SourceRedisTagName
}

// SourceTopic returns the topic of the Redis source
func (r *RedisSourceHandler) SourceTopic() string {
	return "redis"
}

// FetchData fetches data from Redis
func (r *RedisSourceHandler) FetchData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Redis-Source][Current config: " + r.SourceAliasName + "]Start waiting for data to be written...")

	switch r.sourceRedisCfg.DataType {
	case utils.RedisDataTypeSubscribe:
		for msg := range r.redisPubSub.Channel() {
			if r.redisClient == nil {
				logger.Logger.Fatal(utils.LogServiceName +
					"[Redis-Source][Current config: " + r.SourceAliasName + "]Redis client is already closed or not connected!")
				return
			}

			logger.Logger.Debug(utils.LogServiceName +
				"[Redis-Source][Current config: " + r.SourceAliasName + "]Redis consume data: " + msg.Payload)
			r.GetToTransformChan() <- &models.SourceOutput{MetaData: r.MetaData, SourceData: msg.Payload}

			r.Metrics.OnSourceOutput(r.StreamName, r.SourceAliasName)
			r.Metrics.OnSourceOutputSuccess(r.StreamName, r.SourceAliasName)
		}
	case utils.RedisDataTypeLPop, utils.RedisDataTypeRPop:
		for {
			if r.redisClient == nil {
				logger.Logger.Fatal(utils.LogServiceName +
					"[Redis-Source][Current config: " + r.SourceAliasName + "]Redis client is already closed or not connected!")
				return
			}

			if r.sourceRedisCfg.DataType == utils.RedisDataTypeLPop {
				lPopResult, lPopErr := r.redisClient.LPop(context.Background(), r.sourceRedisCfg.KeyOrChannelName).Result()
				if lPopErr != nil && !errors.Is(lPopErr, redis.Nil) {
					logger.Logger.Error(utils.LogServiceName +
						"[Redis-Source][Current config: " + r.SourceAliasName + "]Failed to read by lpop, Reason for exception: " + lPopErr.Error())
					return
				}

				logger.Logger.Debug(utils.LogServiceName +
					"[Redis-Source][Current config: " + r.SourceAliasName + "]Redis lpop consume data: " + lPopResult)
				// 往 Transform 管道写数据
				r.GetToTransformChan() <- &models.SourceOutput{MetaData: r.MetaData, SourceData: lPopResult}
			}

			if r.sourceRedisCfg.DataType == utils.RedisDataTypeRPop {
				rPopResult, rPopErr := r.redisClient.RPop(context.Background(), r.sourceRedisCfg.KeyOrChannelName).Result()
				if rPopErr != nil && !errors.Is(rPopErr, redis.Nil) {
					logger.Logger.Error(utils.LogServiceName +
						"[Redis-Source][Current config: " + r.SourceAliasName + "]Failed to read by rpop, Reason for exception: " + rPopErr.Error())
					return
				}

				logger.Logger.Debug(utils.LogServiceName +
					"[Redis-Source][Current config: " + r.SourceAliasName + "]Redis rpop consume data: " + rPopResult)
				r.GetToTransformChan() <- &models.SourceOutput{MetaData: r.MetaData, SourceData: rPopResult}
			}

			r.Metrics.OnSourceOutput(r.StreamName, r.SourceAliasName)
			r.Metrics.OnSourceOutputSuccess(r.StreamName, r.SourceAliasName)
		}
	}
}

// InitSource initializes the Redis source
func (r *RedisSourceHandler) InitSource() {
	// single/cluster/proxy
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        strings.Split(r.sourceRedisCfg.Address, ","),
		Username:     r.sourceRedisCfg.Username,
		Password:     r.sourceRedisCfg.Password,
		DB:           r.sourceRedisCfg.DBNum,
		PoolSize:     RedisPoolSize,
		MinIdleConns: RedisMinIdles,
	})
	if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[Redis-Source][Current config: " + r.SourceAliasName + "]Failed to connect Redis! Reason for exception: " + pingErr.Error())
		return
	}
	r.redisClient = client

	if r.sourceRedisCfg.DataType == utils.RedisDataTypeSubscribe {
		r.redisPubSub = r.redisClient.Subscribe(context.Background(), r.sourceRedisCfg.KeyOrChannelName)
	}
	r.MetaData = &models.MetaData{
		StreamName:    r.StreamName,
		SourceTagName: r.SourceName(),
		AliasName:     r.SourceAliasName,
	}

	logger.Logger.Info(utils.LogServiceName +
		"[Redis-Source][Current config: " + r.SourceAliasName + "]Init Redis Successful!")
}

// CloseSource closes the Redis source
func (r *RedisSourceHandler) CloseSource() {
	if r.redisClient != nil {
		_ = r.redisClient.Close()
	}
	r.Close()
}

// NewRedisSourceHandler initializes a new Redis source handler
func NewRedisSourceHandler(baseSource source.BaseSource) *RedisSourceHandler {
	handler := &RedisSourceHandler{
		BaseSource:     baseSource,
		sourceRedisCfg: baseSource.SourceConfig.Redis,
	}
	handler.InitSource()
	handler.SetToTransformChan()
	return handler
}
