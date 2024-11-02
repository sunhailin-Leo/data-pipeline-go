package config

// TunnelCfg data tunnel config.
var TunnelCfg *TunnelConfig

// SourceConfig Source config
type SourceConfig struct {
	Type        string                  `json:"type" vd:"len($)>0"`
	SourceName  string                  `json:"source_name" vd:"len($)>0"`
	Kafka       KafkaSourceConfig       `json:"kafka"`
	RocketMQ    RocketMQSourceConfig    `json:"rocketmq"`
	RabbitMQ    RabbitMQSourceConfig    `json:"rabbitmq"`
	PromMetrics PromMetricsSourceConfig `json:"prom_metrics"`
	Pulsar      PulsarSourceConfig      `json:"pulsar"`
	Redis       RedisSourceConfig       `json:"redis"`
}

// TransformSchema transform unit config
type TransformSchema struct {
	SourceKey    string `json:"source_key" vd:"len($)>0"`  // source key
	SinkKey      string `json:"sink_key" vd:"len($)>0"`    // sink key
	Converter    string `json:"converter"`                 // Converter, Like: toInt, toFloat32, toString, etc.
	IsIgnore     bool   `json:"is_ignore"`                 // is ignored key
	IsStrictMode bool   `json:"is_strict_mode"`            // is strict mode
	IsKeepKeys   bool   `json:"is_keep_keys"`              // key is keep origin key
	IsExpand     bool   `json:"is_expand"`                 // is expanded col
	ExpandValue  any    `json:"expand_value"`              // expand value
	SourceName   string `json:"source_name" vd:"len($)>0"` // source alias name
	SinkName     string `json:"sink_name" vd:"len($)>0"`   // sink alias name
}

// TransformJsonPath transform json path config
type TransformJsonPath struct {
	SrcField  string `json:"src_field"`  // source field
	Path      string `json:"path"`       // json path
	DestField string `json:"dest_field"` // destination field
}

// TransformConfig transform config
type TransformConfig struct {
	Mode         string              `json:"mode" vd:"len($)>0"` // Mode: row, json
	Schemas      []TransformSchema   `json:"schemas"`            // Schema
	RowSeparator string              `json:"row_separator"`      // only mode is row will affect, and also only row mode will use strings.Split to split
	Paths        []TransformJsonPath `json:"paths"`              // json paths
}

// SinkConfig Sink config
type SinkConfig struct {
	Type          string                  `json:"type" vd:"len($)>0"`
	SinkName      string                  `json:"sink_name" vd:"len($)>0"`
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
	Name          string          `json:"name" vd:"len($)>0"`   // stream name
	Enable        bool            `json:"enable"`               // is enabled stream
	ChannelSize   int             `json:"channel_size"`         // channel buffer size（default: 0）
	Source        []*SourceConfig `json:"source" vd:"len($)>0"` // source config slice
	SourceAckMode int             `json:"ack_mode"`             // It worked only Source is MQ component; 0: commit after consume; 1: commit after transform; 2: commit after sink
	Transform     TransformConfig `json:"transform"`            // transform config
	Sink          []*SinkConfig   `json:"sink" vd:"len($)>0"`   // sink config slice
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
	Streams []*StreamConfig `json:"streams" vd:"len($)>0"`
}
