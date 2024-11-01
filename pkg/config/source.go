package config

type KafkaSourceConfig struct {
	Address  string `json:"address"`  // kafka address
	Group    string `json:"group"`    // kafka group
	Topic    string `json:"topic"`    // kafka topic
	User     string `json:"user"`     // kafka user
	Password string `json:"password"` // kafka password
}

type RocketMQSourceConfig struct {
	Address   string `json:"address"`    // rocketmq address
	Group     string `json:"group"`      // rocketmq group
	Topic     string `json:"topic"`      // rocketmq topic
	AccessKey string `json:"access_key"` // rocketmq access key
	SecretKey string `json:"secret_key"` // rocketmq secret key
}

type RabbitMQSourceConfig struct {
	Address    string `json:"address"`     // rabbitmq address
	Username   string `json:"username"`    // rabbitmq username
	Password   string `json:"password"`    // rabbitmq password
	VHost      string `json:"v_host"`      // rabbitmq vhost
	Queue      string `json:"queue"`       // rabbitmq queue
	Exchange   string `json:"exchange"`    // rabbitmq exchange
	RoutingKey string `json:"routing_key"` // rabbitmq routing key
}

type PromMetricsSourceConfig struct {
	Address  string `json:"address"`  // prometheus address
	Interval int64  `json:"interval"` // prometheus interval
}

type PulsarSourceConfig struct {
	Address          string `json:"address"`           // pulsar address
	Topic            string `json:"topic"`             // pulsar topic
	SubscriptionName string `json:"subscription_name"` // pulsar subscription name
}

type RedisSourceConfig struct {
	DBNum            int    `json:"db_num"`              // db number（0-15）,cluster mode is invalid
	KeyOrChannelName string `json:"key_or_channel_name"` // Key or Channel name
	Address          string `json:"address"`             // redis hosts
	Username         string `json:"username"`            // Redis Username(6.0+)
	Password         string `json:"password"`            // Redis password
	DataType         string `json:"data_type"`           // data type: lpop, rpop, subscribe
}
