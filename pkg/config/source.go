package config

type KafkaSourceConfig struct {
	Address  string `json:"address"`
	Group    string `json:"group"`
	Topic    string `json:"topic"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RocketMQSourceConfig struct {
	Address   string `json:"address"`
	Group     string `json:"group"`
	Topic     string `json:"topic"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type RabbitMQSourceConfig struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	VHost      string `json:"v_host"`
	Queue      string `json:"queue"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

type PromMetricsSourceConfig struct {
	Address  string `json:"address"`
	Interval int64  `json:"interval"`
}

type PulsarSourceConfig struct {
	Address          string `json:"address"`
	Topic            string `json:"topic"`
	SubscriptionName string `json:"subscription_name"`
}

type RedisSourceConfig struct {
	DBNum            int    `json:"db_num"`              // db number（0-15）,cluster mode is invalid
	KeyOrChannelName string `json:"key_or_channel_name"` // Key or Channel name
	Address          string `json:"address"`             // redis hosts
	Username         string `json:"username"`            // Redis Username(6.0+)
	Password         string `json:"password"`            // Redis password
	DataType         string `json:"data_type"`           // data type: lpop, rpop, subscribe
}
