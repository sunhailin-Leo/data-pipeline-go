package config

// TunnelCfg data tunnel config.
var TunnelCfg *TunnelConfig

// SourceConfig Source config
type SourceConfig struct {
	Type        string                  `json:"type"`
	SourceName  string                  `json:"source_name"`
	Kafka       KafkaSourceConfig       `json:"kafka"`
	RocketMQ    RocketMQSourceConfig    `json:"rocketmq"`
	RabbitMQ    RabbitMQSourceConfig    `json:"rabbitmq"`
	PromMetrics PromMetricsSourceConfig `json:"prom_metrics"`
	Pulsar      PulsarSourceConfig      `json:"pulsar"`
	Redis       RedisSourceConfig       `json:"redis"`
}

// TransformSchema transform unit config
type TransformSchema struct {
	SourceKey    string `json:"source_key"`
	SinkKey      string `json:"sink_key"`
	Converter    string `json:"converter"`      // Converter, Like: toInt, toFloat32, toString, etc.
	IsIgnore     bool   `json:"is_ignore"`      // is ignored key
	IsStrictMode bool   `json:"is_strict_mode"` // is strict mode
	IsKeepKeys   bool   `json:"is_keep_keys"`   // key is keep origin key
	IsExpand     bool   `json:"is_expand"`      // is expanded col
	ExpandValue  any    `json:"expand_value"`   // expand value
	SourceName   string `json:"source_name"`
	SinkName     string `json:"sink_name"`
}

// TransformJsonPath transform json path config
type TransformJsonPath struct {
	SrcField  string `json:"src_field"`  // source field
	Path      string `json:"path"`       // json path
	DestField string `json:"dest_field"` // destination field
}

// TransformConfig transform config
type TransformConfig struct {
	Mode         string              `json:"mode"`          // Mode: row, json
	Schemas      []TransformSchema   `json:"schemas"`       // Schema
	RowSeparator string              `json:"row_separator"` // only mode is row will affect, and also only row mode will use strings.Split to split
	Paths        []TransformJsonPath `json:"paths"`         // json paths
}

// SinkConfig Sink config
type SinkConfig struct {
	Type          string                  `json:"type"`
	SinkName      string                  `json:"sink_name"`
	Clickhouse    ClickhouseSinkConfig    `json:"clickhouse"`
	HTTP          HTTPSinkConfig          `json:"http"`
	Kafka         KafkaSinkConfig         `json:"kafka"`
	Redis         RedisSinkConfig         `json:"redis"`
	LocalFile     LocalFileSinkConfig     `json:"local_file"`
	PostgresSQL   PostgresSQLSinkConfig   `json:"postgres_sql"`
	RocketMQ      RocketMQSinkConfig      `json:"rocketmq"`
	RabbitMQ      RabbitMQSinkConfig      `json:"rabbitmq"`
	Oracle        OracleSinkConfig        `json:"oracle"`
	MySQL         MySQLSinkConfig         `json:"mysql"`
	Pulsar        PulsarSinkConfig        `json:"pulsar"`
	Elasticsearch ElasticsearchSinkConfig `json:"elasticsearch"`
}

// StreamConfig stream config
type StreamConfig struct {
	Name          string          `json:"name"`         // stream name
	Enable        bool            `json:"enable"`       // is enabled stream
	ChannelSize   int             `json:"channel_size"` // 通道缓冲区大小（用于削峰，不给默认 0）
	Source        []*SourceConfig `json:"source"`       // 数据源配置项
	SourceAckMode int             `json:"ack_mode"`     // 仅当 Source 时 MQ 组件才有效; 0: 消费后提交; 1: 转换成功后提交; 2: 发送后提交
	Transform     TransformConfig `json:"transform"`    // 数据转换配置项
	Sink          []*SinkConfig   `json:"sink"`         // 数据目标配置项
}

// GetSourceBySourceName based on SourceName to get Source
func (s *StreamConfig) GetSourceBySourceName(name string) *SourceConfig {
	for _, source := range s.Source {
		if source.SourceName == name {
			return source
		}
	}
	return nil
}

// GetSourceTagBySourceName based on SourceName to get Source Tag
func (s *StreamConfig) GetSourceTagBySourceName(name string) string {
	for _, source := range s.Source {
		if source.SourceName == name {
			return source.Type
		}
	}
	return ""
}

// GetSinkBySinkName based on SinkName to get Sink
func (s *StreamConfig) GetSinkBySinkName(name string) *SinkConfig {
	for _, sink := range s.Sink {
		if sink.SinkName == name {
			return sink
		}
	}
	return nil
}

// GetSinkTagBySinkName based on SinkName to get Sink Tag
func (s *StreamConfig) GetSinkTagBySinkName(name string) string {
	for _, sink := range s.Sink {
		if sink.SinkName == name {
			return sink.Type
		}
	}
	return ""
}

// TunnelConfig data tunnel config struct.
type TunnelConfig struct {
	Streams []*StreamConfig `json:"streams"`
}
