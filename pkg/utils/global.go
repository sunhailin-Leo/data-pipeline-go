package utils

const (
	ServiceName string = "data-pipeline-go"

	LogMaxAge           int    = 7
	LogMaxSize          int    = 128
	LogMaxBackups       int    = 30
	LogConsoleSeparator string = "|-|"
	LogServiceName             = ServiceName + LogConsoleSeparator

	PromHTTPServerPort string = "8080"
	PromHTTPRoute      string = "/metrics"

	TransformRowMode      string = "row"
	TransformJsonMode     string = "json"
	TransformJsonPathMode string = "jsonPath"

	AckModeInSource    int = 0
	AckModeInTransform int = 1
	AckModeInSink      int = 2

	ConfigFromLocalTagName     string = "local"
	ConfigFromApolloTagName    string = "apollo"
	ConfigFromRedisTagName     string = "redis"
	ConfigFromNacosTagName     string = "nacos"
	ConfigFromZookeeperTagName string = "zookeeper"

	ConfigFromSourceName             string = "CONFIG_SRC"            // Required
	ConfigFromLocalPathEnvName       string = "LOCAL_PATH"            // If CONFIG_SRC = local, then required
	ConfigFromApolloEnvHost          string = "APOLLO_HOST"           // If CONFIG_SRC = apollo, then required
	ConfigFromApolloEnvAppId         string = "APOLLO_APP_ID"         // optional
	ConfigFromApolloEnvNamespace     string = "APOLLO_NAMESPACE"      // optional
	ConfigFromApolloEnvClusterKey    string = "APOLLO_CLUSTER_KEY"    // optional
	ConfigFromApolloEnvConfigKey     string = "APOLLO_CONFIG_KEY"     // If CONFIG_SRC = apollo, then required
	ConfigFromRedisEnvHost           string = "REDIS_HOST"            // If CONFIG_SRC = redis, then required
	ConfigFromRedisEnvUsername       string = "REDIS_USERNAME"        // optional
	ConfigFromRedisEnvPassword       string = "REDIS_PASSWORD"        // optional
	ConfigFromRedisEnvDBNum          string = "REDIS_DB"              // optional
	ConfigFromRedisEnvConfigKey      string = "REDIS_CONFIG_KEY"      // If CONFIG_SRC = redis, then required
	ConfigFromNacosEnvServerIP       string = "NACOS_IP"              // If CONFIG_SRC = nacos, then required
	ConfigFromNacosEnvServerPort     string = "NACOS_PORT"            // If CONFIG_SRC = nacos, then required
	ConfigFromNacosEnvNamespaceId    string = "NACOS_NAMESPACE_ID"    // optional
	ConfigFromNacosEnvDataId         string = "NACOS_DATA_ID"         // If CONFIG_SRC = nacos, then required
	ConfigFromNacosEnvGroup          string = "NACOS_GROUP"           // If CONFIG_SRC = nacos, then required
	ConfigFromZookeeperEnvHosts      string = "ZOOKEEPER_HOSTS"       // If CONFIG_SRC = zookeeper, then required
	ConfigFromZookeeperEnvConfigPath string = "ZOOKEEPER_CONFIG_PATH" // If CONFIG_SRC = zookeeper, then required

	SourceKafkaTagName       string = "Kafka"
	SourceRocketMQTagName    string = "RocketMQ"
	SourceRabbitMQTagName    string = "RabbitMQ"
	SourcePromMetricsTagName string = "PromMetrics"
	SourcePulsarTagName      string = "Pulsar"
	SourceRedisTagName       string = "Redis"

	SinkClickhouseTagName    string = "ClickHouse"
	SinkConsoleTagName       string = "Console"
	SinkHTTPTagName          string = "HTTP"
	SinkKafkaTagName         string = "Kafka"
	SinkRedisTagName         string = "Redis"
	SinkLocalFileTagName     string = "LocalFile"
	SinkPostgresSQLTagName   string = "PostgresSQL"
	SinkRocketMQTagName      string = "RocketMQ"
	SinkRabbitMQTagName      string = "RabbitMQ"
	SinkOracleTagName        string = "Oracle"
	SinkMySQLTagName         string = "MySQL"
	SinkPulsarTagName        string = "Pulsar"
	SinkElasticsearchTagName string = "Elasticsearch"

	HTTPContentTypeJSON     string = "application/json"
	RedisDynamicKeyFlagName string = "fromTransformHead"
	RedisDataTypeKV         string = "kv"
	RedisDataTypeHash       string = "hash"
	RedisDataTypeLPush      string = "lpush"
	RedisDataTypeLPop       string = "lpop"
	RedisDataTypeRPush      string = "rpush"
	RedisDataTypeRPop       string = "rpop"
	RedisDataTypeSet        string = "set"
	RedisDataTypePublish    string = "publish"
	RedisDataTypeSubscribe  string = "subscribe" // 默认用 PSUBSCRIBE
	LocalFileCSVFormatType  string = "csv"
	LocalFileTextFormatType string = "text"
)
