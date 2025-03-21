package config

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/bytedance/sonic"
	"github.com/go-zookeeper/zk"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/mmap"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type TunnelConfigLoader struct {
	config *TunnelConfig
}

func (c *TunnelConfigLoader) bindAllConfigEnv() {
	// Load from source
	_ = viper.BindEnv(utils.ConfigFromSourceName)
	// Local
	_ = viper.BindEnv(utils.ConfigFromLocalPathEnvName)
	// Apollo
	_ = viper.BindEnv(utils.ConfigFromApolloEnvHost)
	_ = viper.BindEnv(utils.ConfigFromApolloEnvAppId)
	_ = viper.BindEnv(utils.ConfigFromApolloEnvNamespace)
	_ = viper.BindEnv(utils.ConfigFromApolloEnvClusterKey)
	_ = viper.BindEnv(utils.ConfigFromApolloEnvConfigKey)
	// Redis
	_ = viper.BindEnv(utils.ConfigFromRedisEnvHost)
	_ = viper.BindEnv(utils.ConfigFromRedisEnvUsername)
	_ = viper.BindEnv(utils.ConfigFromRedisEnvPassword)
	_ = viper.BindEnv(utils.ConfigFromRedisEnvDBNum)
	_ = viper.BindEnv(utils.ConfigFromRedisEnvConfigKey)
	// Nacos
	_ = viper.BindEnv(utils.ConfigFromNacosEnvServerIP)
	_ = viper.BindEnv(utils.ConfigFromNacosEnvServerPort)
	_ = viper.BindEnv(utils.ConfigFromNacosEnvNamespaceId)
	_ = viper.BindEnv(utils.ConfigFromNacosEnvDataId)
	_ = viper.BindEnv(utils.ConfigFromNacosEnvGroup)
	// Zookeeper
	_ = viper.BindEnv(utils.ConfigFromZookeeperEnvHosts)
	_ = viper.BindEnv(utils.ConfigFromZookeeperEnvConfigPath)
	// HTTP
	_ = viper.BindEnv(utils.ConfigFromHTTPEnvHosts)
	_ = viper.BindEnv(utils.ConfigFromHTTPEnvConfigURI)
	_ = viper.BindEnv(utils.ConfigFromHTTPEnvHeartBeatURI)
	_ = viper.BindEnv(utils.ConfigFromHTTPEnvHeartBeatIntervalSecs)
}

func (c *TunnelConfigLoader) IsConfigLoaded() bool {
	return c.config != nil
}

func (c *TunnelConfigLoader) validateConfig(cfg *TunnelConfig) error {
	// vd check
	if validateErr := vd.Validate(cfg); validateErr != nil {
		return validateErr
	}

	// duplicate check
	streamNameMap := make(map[string]struct{})
	for _, stream := range cfg.Streams {
		_, ok := streamNameMap[stream.Name]
		if ok {
			return errors.New("stream name " + stream.Name + " is duplicated")
		}
	}
	return nil
}

func (c *TunnelConfigLoader) loadFile(data []byte) *TunnelConfig {
	defaultConfig := TunnelConfig{}
	if err := sonic.Unmarshal(data, &defaultConfig); err != nil {
		logger.Logger.Fatal("Project configuration file parsing failed! Error reason: " + err.Error())
		os.Exit(1)
	}

	if validateErr := c.validateConfig(&defaultConfig); validateErr != nil {
		logger.Logger.Fatal("Project configuration file validate failed! Error reason: " + validateErr.Error())
		os.Exit(1)
	}

	return &defaultConfig
}

func (c *TunnelConfigLoader) loadFromLocal(path string) {
	// use mmap to load file
	reader, readerErr := mmap.Open(path)
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
	// TODO live reload
}

func (c *TunnelConfigLoader) loadFromApollo() {
	// Host
	apolloHostObj := viper.Get(utils.ConfigFromApolloEnvHost)
	if apolloHostObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Apollo config error, reason: APOLLO_HOST is empty")
		os.Exit(1)
	}
	apolloHost := cast.ToString(apolloHostObj)

	// AppID
	var apolloAppId string
	apolloAppIdObj := viper.Get(utils.ConfigFromApolloEnvAppId)
	if apolloAppIdObj == nil {
		apolloAppId = utils.ServiceName
		logger.Logger.Warn(utils.LogServiceName + "Load Apollo config error, reason: APOLLO_APP_ID is empty, use default AppID to instead")
	} else {
		apolloAppId = cast.ToString(apolloAppIdObj)
	}

	// Namespace
	var apolloNamespace string
	apolloNamespaceObj := viper.Get(utils.ConfigFromApolloEnvNamespace)
	if apolloNamespaceObj == nil {
		apolloNamespace = "application"
		logger.Logger.Warn(utils.LogServiceName + "Load Apollo config error, reason: APOLLO_NAMESPACE is empty, use default NamespaceName to instead")
	} else {
		apolloNamespace = cast.ToString(apolloNamespaceObj)
	}

	// Cluster key
	var apolloClusterKey string
	apolloClusterKeyObj := viper.Get(utils.ConfigFromApolloEnvClusterKey)
	if apolloClusterKeyObj == nil {
		apolloClusterKey = "default"
		logger.Logger.Warn(utils.LogServiceName + "Load Apollo config error, reason: APOLLO_CLUSTER_KEY is empty, use default NamespaceName to instead")
	} else {
		apolloClusterKey = cast.ToString(apolloClusterKeyObj)
	}

	// Config Key
	var apolloConfigKey string
	apolloConfigKeyObj := viper.Get(utils.ConfigFromApolloEnvConfigKey)
	if apolloConfigKeyObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Apollo config error, reason: APOLLO_CONFIG_KEY is empty")
		os.Exit(1)
	}
	apolloConfigKey = cast.ToString(apolloConfigKeyObj)

	// apollo config
	cfg := &config.AppConfig{
		AppID:             apolloAppId,
		Cluster:           apolloClusterKey,
		IP:                apolloHost,
		NamespaceName:     apolloNamespace,
		SyncServerTimeout: 3, // hard code 3 secs
		MustStart:         true,
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) { return cfg, nil })
	if err != nil {
		logger.Logger.Fatal(utils.LogServiceName + "load Apollo config error,reason: " + err.Error())
		return
	}

	// get config from apollo configuration
	cache := client.GetConfigCache(apolloNamespace)
	count := 0
	cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count < 1 {
		logger.Logger.Warn(utils.LogServiceName + "Apollo config is empty! downgrade to using default config")
	} else {
		value, loadCacheErr := cache.Get(apolloConfigKey)
		if loadCacheErr != nil {
			logger.Logger.Error(utils.LogServiceName + "Apollo local config is empty! Error reason: " + loadCacheErr.Error())
		}
		TunnelCfg = c.loadFile([]byte(value.(string)))
		logger.Logger.Info(utils.LogServiceName + "Load Apollo config is successful!")
	}
	// TODO live reload
}

func (c *TunnelConfigLoader) loadFromRedis() {
	// Redis Host
	redisHostObj := viper.Get(utils.ConfigFromRedisEnvHost)
	if redisHostObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Redis config error, reason: REDIS_HOST is empty")
		os.Exit(1)
	}
	redisHosts := cast.ToString(redisHostObj)

	// Redis Username (Redis 6.0+)
	var redisUsername string
	redisUsernameObj := viper.Get(utils.ConfigFromRedisEnvUsername)
	if redisUsernameObj != nil {
		redisUsername = cast.ToString(redisUsernameObj)
	}

	// Redis Password
	var redisPassword string
	redisPasswordObj := viper.Get(utils.ConfigFromRedisEnvPassword)
	if redisPasswordObj != nil {
		redisPassword = cast.ToString(redisPasswordObj)
	}

	// Redis DB Num
	var redisDBNum int
	redisDBNumObj := viper.Get(utils.ConfigFromRedisEnvDBNum)
	if redisDBNumObj != nil {
		redisDBNum = cast.ToInt(redisDBNumObj)
	} else {
		// maybe is cluster mode
		redisDBNum = 0
	}

	// Redis Config Key
	var redisConfigKey string
	redisConfigKeyObj := viper.Get(utils.ConfigFromRedisEnvConfigKey)
	if redisConfigKeyObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Redis config error, reason: REDIS_CONFIG_KEY is empty")
		os.Exit(1)
	}
	redisConfigKey = cast.ToString(redisConfigKeyObj)

	// Redis Client
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(redisHosts, ","),
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisDBNum,
	})
	if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
		logger.Logger.Error(utils.LogServiceName + "Failed to connect Redis! Reason for exception: " + pingErr.Error())
		return
	}

	rConfig, getErr := client.Get(context.Background(), redisConfigKey).Result()
	if getErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "load Redis config error,reason: " + getErr.Error())
		return
	}
	TunnelCfg = c.loadFile([]byte(rConfig))
	logger.Logger.Info(utils.LogServiceName + "Load Redis config is successful!")
	// TODO live reload
}

func (c *TunnelConfigLoader) loadFromNacos() {
	// Nacos host and port
	nacosHostObj := viper.Get(utils.ConfigFromNacosEnvServerIP)
	nacosPortObj := viper.Get(utils.ConfigFromNacosEnvServerPort)
	if nacosHostObj == nil || nacosPortObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Nacos config error, reason: NACOS_IP or NACOS_PORT is empty")
		os.Exit(1)
	}

	// Nacos namespace id
	var nacosNamespaceId string
	nacosNamespaceIdObj := viper.Get(utils.ConfigFromNacosEnvNamespaceId)
	if nacosNamespaceIdObj != nil {
		nacosNamespaceId = cast.ToString(nacosNamespaceIdObj)
	}

	// Nacos server config (not support multi cluster mode)
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: cast.ToString(nacosHostObj),
			Port:   cast.ToUint64(nacosPortObj),
		},
	}

	// Nacos client config
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosNamespaceId, // 如果不需要命名空间，可以留空
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// create config client
	configClient, createClientErr := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if createClientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "failed to create Nacos client, Reason for exception: " + createClientErr.Error())
		return
	}

	nacosDataIdObj := viper.Get(utils.ConfigFromNacosEnvDataId)
	nacosGroupObj := viper.Get(utils.ConfigFromNacosEnvGroup)
	if nacosDataIdObj == nil || nacosGroupObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Nacos config error, reason: NACOS_DATA_ID or NACOS_GROUP is empty")
		os.Exit(1)
	}

	// get config
	content, getConfigErr := configClient.GetConfig(vo.ConfigParam{
		DataId: cast.ToString(nacosDataIdObj),
		Group:  cast.ToString(nacosGroupObj),
	})
	if getConfigErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "failed to get Nacos config, Reason for exception: " + getConfigErr.Error())
		return
	}
	TunnelCfg = c.loadFile([]byte(content))
	logger.Logger.Info(utils.LogServiceName + "Load Nacos config is successful!")
	// TODO live reload
}

func (c *TunnelConfigLoader) loadFromZookeeper() {
	// Zookeeper Hosts
	zkHostsObj := viper.Get(utils.ConfigFromZookeeperEnvHosts)
	if zkHostsObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Zookeeper config error, reason: ZOOKEEPER_HOSTS is empty")
		os.Exit(1)
	}

	conn, _, connErr := zk.Connect(strings.Split(cast.ToString(zkHostsObj.(string)), ","), time.Second)
	if connErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "failed to create Zookeeper client, Reason for exception: " + connErr.Error())
		return
	}
	defer conn.Close()

	// ConfigPath
	zkConfigPathObj := viper.Get(utils.ConfigFromZookeeperEnvConfigPath)
	if zkConfigPathObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load Zookeeper config error, reason: ZOOKEEPER_CONFIG_PATH is empty")
		return
	}

	configData, _, getErr := conn.Get(cast.ToString(zkConfigPathObj))
	if getErr != nil {
		logger.Logger.Fatal(utils.LogServiceName + "failed to get Zookeeper config, Reason for exception: " + getErr.Error())
		return
	}
	TunnelCfg = c.loadFile(configData)
	logger.Logger.Info(utils.LogServiceName + "Load Zookeeper config is successful!")
	// TODO live reload
}

func (c *TunnelConfigLoader) loadFromHTTP() {
	httpHostsObj := viper.Get(utils.ConfigFromHTTPEnvHosts)
	if httpHostsObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load HTTP config error, reason: HTTP_HOSTS is empty")
		os.Exit(1)
	}

	httpUriObj := viper.Get(utils.ConfigFromHTTPEnvConfigURI)
	if httpUriObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load HTTP config error, reason: HTTP_CONFIG_URI is empty")
		os.Exit(1)
	}

	// init client, request and response
	client := &fasthttp.Client{}
	clientReq := fasthttp.AcquireRequest()
	clientReq.SetRequestURI(cast.ToString(httpHostsObj) + cast.ToString(httpUriObj))
	clientReq.Header.SetContentType("application/json")
	clientReq.Header.SetMethod(fasthttp.MethodGet)
	clientResp := fasthttp.AcquireResponse()
	// do request
	if err := client.Do(clientReq, clientResp); err != nil {
		logger.Logger.Fatal(utils.LogServiceName + "failed to load HTTP config, Reason for exception: " + err.Error())
		os.Exit(1)
	}

	// set config
	TunnelCfg = c.loadFile(clientResp.Body())
	logger.Logger.Info(utils.LogServiceName + "Load HTTP config is successful!")

	// release
	fasthttp.ReleaseRequest(clientReq)
	fasthttp.ReleaseResponse(clientResp)

	// heart beat
	httpHeartBeatUriObj := viper.Get(utils.ConfigFromHTTPEnvHeartBeatURI)
	if httpHeartBeatUriObj != nil {
		logger.Logger.Fatal(utils.LogServiceName + "Load HTTP config error, reason: HTTP_HEARTBEAT_URI is empty")
		os.Exit(1)
	} else {
		httpHeartBeatIntervalSecsObj := viper.Get(utils.ConfigFromHTTPEnvHeartBeatIntervalSecs)
		if httpHeartBeatIntervalSecsObj == nil {
			// default interval 15 secs
			httpHeartBeatIntervalSecsObj = 15
		}

		// heart beat
		go func(c *fasthttp.Client) {
			for {
				heartBeatReq := fasthttp.AcquireRequest()
				heartBeatReq.SetRequestURI(cast.ToString(httpHostsObj) + cast.ToString(httpHeartBeatUriObj))
				heartBeatReq.Header.SetMethod(fasthttp.MethodGet)

				if heartBeatErr := c.Do(heartBeatReq, nil); heartBeatErr != nil {
					logger.Logger.Error(utils.LogServiceName + "failed to keep heartbeat, Reason for exception: " + heartBeatErr.Error())
					break
				}

				time.Sleep(time.Duration(cast.ToInt(httpHeartBeatIntervalSecsObj)) * time.Second)
			}
		}(client)
	}

	// TODO live reload
}

// Load load config api
func (c *TunnelConfigLoader) Load() {
	// bind viper setting
	c.bindAllConfigEnv()
	// Check source env
	loadSrcObj := viper.Get(utils.ConfigFromSourceName)
	if loadSrcObj == nil {
		logger.Logger.Fatal(utils.LogServiceName + "config load error! CONFIG_SRC env is not set")
		os.Exit(1)
	}

	switch cast.ToString(loadSrcObj) {
	case utils.ConfigFromLocalTagName:
		localPathObj := viper.Get(utils.ConfigFromLocalPathEnvName)
		if localPathObj == nil {
			logger.Logger.Fatal(utils.LogServiceName + "config load error! LOCAL_PATH env is not set!")
			os.Exit(1)
		}
		c.loadFromLocal(cast.ToString(localPathObj))
	case utils.ConfigFromApolloTagName:
		c.loadFromApollo()
	case utils.ConfigFromRedisTagName:
		c.loadFromRedis()
	case utils.ConfigFromNacosTagName:
		c.loadFromNacos()
	case utils.ConfigFromZookeeperTagName:
		c.loadFromZookeeper()
	case utils.ConfigFromHTTPTagName:
		c.loadFromHTTP()
	default:
		logger.Logger.Fatal(utils.LogServiceName + "config load error! CONFIG_SRC env is unsupported!")
		os.Exit(1)
	}

	if TunnelCfg != nil {
		c.config = TunnelCfg
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
