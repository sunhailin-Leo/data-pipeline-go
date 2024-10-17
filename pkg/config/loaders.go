package config

import (
	"context"
	_ "embed"
	"os"
	"strings"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/mmap"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

//go:embed config.json
var defaultConfigFile []byte

type TunnelConfigLoader struct {
	config *TunnelConfig

	apolloClient agollo.Client
	redisClient  redis.UniversalClient
}

func (c *TunnelConfigLoader) IsConfigLoaded() bool {
	return c.config != nil
}

func (c *TunnelConfigLoader) loadFile(data []byte) *TunnelConfig {
	defaultConfig := TunnelConfig{}
	if err := sonic.Unmarshal(data, &defaultConfig); err != nil {
		logger.Logger.Fatal("Project configuration file parsing failed! Error reason: " + err.Error())
		os.Exit(1)
	}
	return &defaultConfig
}

// LoadDefaultConfig load default config
func (c *TunnelConfigLoader) loadDefaultConfig() {
	TunnelCfg = c.loadFile(defaultConfigFile)
	c.config = TunnelCfg
}

func (c *TunnelConfigLoader) loadFromLocal() {
	// use mmap to load file
	reader, readerErr := mmap.Open(c.config.Config.Local.Path)
	if readerErr != nil {
		logger.Logger.Fatal("Project configuration file read failed! Error reason: " + readerErr.Error())
		os.Exit(1)
	}
	// create cache buffer to get file content
	data := make([]byte, reader.Len())
	if _, err := reader.ReadAt(data, 0); err != nil {
		logger.Logger.Fatal("Project configuration file parsing failed! Error reason: " + err.Error())
		os.Exit(1)
	}
	TunnelCfg = c.loadFile(data)
}

func (c *TunnelConfigLoader) loadFromApollo() {
	apolloHost := c.config.Config.Apollo.Host
	if apolloHost == "" {
		logger.Logger.Fatal(utils.LogServiceName + "Load Apollo config error, reason: host is empty")
		os.Exit(1)
	}
	apolloAppId := c.config.Config.Apollo.AppID
	if apolloAppId == "" {
		apolloAppId = utils.ServiceName
		logger.Logger.Warn(utils.LogServiceName + "Load Apollo config error, reason: appId is empty, use default appId instead")
	}
	logger.Logger.Info(utils.LogServiceName + "Current using Apollo host: " + apolloHost)

	cfg := &config.AppConfig{
		AppID:             apolloAppId,
		Cluster:           c.config.Config.Apollo.ClusterKey,
		IP:                apolloHost,
		NamespaceName:     c.config.Config.Apollo.Namespace,
		SyncServerTimeout: 3, // hard code 3 secs
		MustStart:         true,
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) { return cfg, nil })
	if err != nil {
		logger.Logger.Fatal(utils.LogServiceName + "load Apollo config error,reason: " + err.Error())
		return
	}
	c.apolloClient = client

	// get config from apollo configuration
	cache := c.apolloClient.GetConfigCache(c.config.Config.Apollo.Namespace)
	count := 0
	cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count < 1 {
		logger.Logger.Warn(utils.LogServiceName + "Apollo config is empty! downgrade to using default config")
	} else {
		value, loadCacheErr := cache.Get(c.config.Config.Apollo.ConfigKey)
		if loadCacheErr != nil {
			logger.Logger.Error(utils.LogServiceName + "Apollo local config is empty! Error reason: " + loadCacheErr.Error())
		}
		TunnelCfg = c.loadFile([]byte(value.(string)))
		logger.Logger.Info(utils.LogServiceName + "Load Apollo config is successful!")
	}
	// TODO 后续再实现实时更新的操作
}

func (c *TunnelConfigLoader) loadFromRedis() {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(c.config.Config.Redis.Host, ","),
		Username: c.config.Config.Redis.Username,
		Password: c.config.Config.Redis.Password,
		DB:       c.config.Config.Redis.DB,
	})
	if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
		logger.Logger.Error(utils.LogServiceName + "Failed to connect Redis! Reason for exception: " + pingErr.Error())
		return
	}
	c.redisClient = client

	rConfig, getErr := c.redisClient.Get(context.Background(), c.config.Config.Redis.ConfigKey).Result()
	if getErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "load Redis config error,reason: " + getErr.Error())
		return
	}
	TunnelCfg = c.loadFile([]byte(rConfig))
	logger.Logger.Info(utils.LogServiceName + "Load Redis config is successful!")
	// TODO 后续再实现实时更新的操作
}

// Load load config api
func (c *TunnelConfigLoader) Load() {
	c.loadDefaultConfig()

	if c.config.Config != nil {
		switch c.config.Config.From {
		case utils.ConfigFromLocalTagName:
			c.loadFromLocal()
		case utils.ConfigFromApolloTagName:
			c.loadFromApollo()
		case utils.ConfigFromRedisTagName:
			c.loadFromRedis()
		}
	}
}

func NewTunnelConfigLoader() {
	loader := &TunnelConfigLoader{}
	loader.Load()
	if loader.IsConfigLoaded() {
		logger.Logger.Info(utils.LogServiceName + "config load successful!")
	} else {
		logger.Logger.Fatal(utils.LogServiceName + "config load error!")
		os.Exit(1)
	}
}
