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

	TransformRowMode  string = "row"
	TransformJsonMode string = "json"

	AckModeInSource    int = 0
	AckModeInTransform int = 1
	AckModeInSink      int = 2

	ConfigFromLocalTagName  string = "local"
	ConfigFromApolloTagName string = "apollo"
	ConfigFromRedisTagName  string = "redis"

	SourceKafkaTagName       string = "Kafka"
	SourceRocketMQTagName    string = "RocketMQ"
	SourceRabbitMQTagName    string = "RabbitMQ"
	SourcePromMetricsTagName string = "PromMetrics"
	SourcePulsarTagName      string = "Pulsar"
	SourceRedisTagName       string = "Redis"

	SinkClickhouseTagName  string = "ClickHouse"
	SinkConsoleTagName     string = "Console"
	SinkHTTPTagName        string = "HTTP"
	SinkKafkaTagName       string = "Kafka"
	SinkRedisTagName       string = "Redis"
	SinkLocalFileTagName   string = "LocalFile"
	SinkPostgresSQLTagName string = "PostgresSQL"
	SinkRocketMQTagName    string = "RocketMQ"
	SinkRabbitMQTagName    string = "RabbitMQ"
	SinkOracleTagName      string = "Oracle"
	SinkMySQLTagName       string = "MySQL"
	SinkPulsarTagName      string = "Pulsar"

	HTTPContentTypeJSON        string = "application/json"
	RedisDynamicKeyFlagName    string = "fromTransformHead"
	RedisDataTypeKV            string = "kv"
	RedisDataTypeHash          string = "hash"
	RedisDataTypeLPush         string = "lpush"
	RedisDataTypeLPop          string = "lpop"
	RedisDataTypeRPush         string = "rpush"
	RedisDataTypeRPop          string = "rpop"
	RedisDataTypeSet           string = "set"
	RedisDataTypePublish       string = "publish"
	RedisDataTypeSubscribe     string = "subscribe" // 默认用 PSUBSCRIBE
	LocalFileCSVFormatType     string = "csv"
	LocalFileParquetFormatType string = "parquet"
	LocalFileTextFormatType    string = "text"
)
