package config

// ClickhouseTableColumn Clickhouse table column name
type ClickhouseTableColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Compress string `json:"compress"`
	Comment  string `json:"comment"`
}

// ClickhouseSinkConfig Clickhouse Sink config
type ClickhouseSinkConfig struct {
	Address           string                  `json:"address"`
	Username          string                  `json:"username"`
	Password          string                  `json:"password"`
	Database          string                  `json:"database"`
	TableName         string                  `json:"table_name"`
	BulkSize          int                     `json:"bulk_size"`
	IsAutoCreateTable bool                    `json:"is_auto_create_table"`
	Columns           []ClickhouseTableColumn `json:"columns"`
	Engine            string                  `json:"engine"`
	Partition         []string                `json:"partition"`
	PrimaryKey        []string                `json:"primary_key"`
	OrderBy           []string                `json:"order_by"`
	TTL               string                  `json:"ttl"`
	Comment           string                  `json:"comment"`
	Settings          []string                `json:"settings"`
}

// HTTPSinkConfig HTTP Sink config
type HTTPSinkConfig struct {
	URL                     string            `json:"url"`
	ContentType             string            `json:"content_type"`
	ReadTimeoutSecs         int               `json:"read_timeout_secs"`
	WriteTimeoutSecs        int               `json:"write_timeout_secs"`
	MaxIdleConnDurationSecs int               `json:"max_idle_conn_duration_secs"`
	MaxConnWaitTimeoutSecs  int               `json:"max_conn_wait_timeout_secs"`
	Headers                 map[string]string `json:"headers"`
}

// KafkaSinkConfig Kafka Sink config
type KafkaSinkConfig struct {
	Address     string `json:"address"`      // Kafka address
	Topic       string `json:"topic"`        // Kafka Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}

// RedisSinkConfig Redis Sink config
type RedisSinkConfig struct {
	DBNum            int    `json:"db_num"`              // db number（0-15）,cluster mode is invalid
	Expire           int    `json:"expire"`              // expire time, unit: seconds
	Address          string `json:"address"`             // redis hosts
	Username         string `json:"username"`            // Redis Username(6.0+)
	Password         string `json:"password"`            // Redis password
	DataType         string `json:"data_type"`           // data type: kv, hash, lpush, rpush, set, publish
	KeyOrChannelName string `json:"key_or_channel_name"` // Key or Channel name, if you want key is variable, please use "fromTransformHead"
}

// LocalFileSinkConfig LocalFile Sink config
type LocalFileSinkConfig struct {
	FileName       string `json:"file_name"`
	FileFormatType string `json:"file_format_type"` // file format type: text, csv
	RowDelimiter   string `json:"row_delimiter"`    // only file_format_type = text is affect
}

// PostgresSQLSinkConfig PostgresSQL Sink config
type PostgresSQLSinkConfig struct {
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	TableName string `json:"table_name"`
	BulkSize  int    `json:"bulk_size"`
}

// RocketMQSinkConfig RocketMQ Sink config
type RocketMQSinkConfig struct {
	Address     string `json:"address"`      // RocketMQ address
	Topic       string `json:"topic"`        // RocketMQ Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}

// RabbitMQSinkConfig RabbitMQ Sink config
type RabbitMQSinkConfig struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	VHost      string `json:"v_host"`
	Queue      string `json:"queue"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

// OracleSinkConfig Oracle Sink config
type OracleSinkConfig struct {
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	TableName string `json:"table_name"`
	BulkSize  int    `json:"bulk_size"`
}

// MySQLSinkConfig MySQL Sink config
type MySQLSinkConfig struct {
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	TableName string `json:"table_name"`
	BulkSize  int    `json:"bulk_size"`
}

// PulsarSinkConfig Pulsar Sink config
type PulsarSinkConfig struct {
	Address string `json:"address"` // Pulsar host
	Topic   string `json:"topic"`   // Pulsar Topic
}

// ElasticsearchSinkConfig Elasticsearch Sink config
type ElasticsearchSinkConfig struct {
	Address   string `json:"address"`     // Elasticsearch address
	Username  string `json:"username"`    // Elasticsearch username
	Password  string `json:"password"`    // Elasticsearch password
	IndexName string `json:"index_name"`  // Elasticsearch index name
	DocIdName string `json:"doc_id_name"` // Elasticsearch document id name (will take it from transform data)
	Version   string `json:"version"`     // Only use 7.X or 8.X
}
