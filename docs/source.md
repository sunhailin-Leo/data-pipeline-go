Source
===
## 数据输入源
支持列表：
- [X] Kafka
- [X] RocketMQ
- [X] Redis
    - [X] List
        - [X] LPop
        - [X] RPop
    - [X] Pub/Sub
        - [X] Subscribe(default using PSubscribe)
- [X] Pulsar
- [X] RabbitMQ
- [X] Prometheus metrics exporter(etc. get /metrics)

## 字段说明
**example:**
```json
{
  ...
  "source": [
    {
      "type": "Kafka",
      "source_name": "kafka-1",
      "kafka": {
        "address": "<Kafka Hosts>",
        "group": "<Kafka Group>",
        "topic": "<Kafka Topic>"
      }
    }
  ],
  ...
}
```

### source
类型: Array

说明: source 字段是一个数据源的数组，支持有多个不同类型的数据源。

### type
类型: String

说明: 表示数据源的类型,参考 `支持列表`。对于 Kafka 数据源，这里应填写 "Kafka"。该字段用于指明所使用的数据源类型，以便系统能够正确处理和解析。

### source_name
类型: String

说明: 这是该 Kafka 数据源的名称。source_name 用于标识具体的数据源，方便开发者在系统中管理多个数据源。例如，名称可以是 "kafka-1"，用来表示 Kafka 数据源的一个实例

### kafka
类型: Object
```json
    "kafka": {
        "address": "<Kafka Hosts>",           // Kafka集群地址，例如 "kafka1:9092,kafka2:9092"
        "group": "<Kafka Group>",             // Kafka消费者组名称，例如 "consumer-group-1"
        "topic": "<Kafka Topic>"              // 要消费的Kafka主题名称，例如 "my_topic"
    }
```
包含与 Kafka 相关的配置信息，包括 Kafka 集群的地址、消费者组、以及要消费的主题。以下是 kafka 对象中的具体字段：

* address: Kafka 集群地址，格式为 <host>:<port>。可以是单个 broker 地址，也可以是多个，用逗号分隔。例如："kafka1:9092,kafka2:9092"。

* group: Kafka 消费者组名称，用于多个消费者在同一个组内共享负载，确保每条消息只会被组内一个消费者处理。例如："consumer-group-1"。

* topic: Kafka 主题名称，代表消息分类。消费者将从指定的主题中获取消息。该名称应与 Kafka 集群中的主题一致。例如："my_topic"。

### rocketmq
类型: Object
```json
    "rocketmq":{
        "address": "<RocketMQ Hosts>",             // RocketMQ集群地址，例如 "rocketmq.example.com:9876"
        "group": "<RocketMQ Group>",               // RocketMQ消费者组名称，例如 "my_consumer_group"
        "topic": "<RocketMQ Topic>",               // 要消费的RocketMQ主题名称，例如 "my_topic"
        "access_key": "<RocketMQ Access Key>",     // 用于访问RocketMQ的访问密钥，例如 "my_access_key"
        "secret_key": "<RocketMQ Secret Key>"      // 用于访问RocketMQ的安全密钥，例如 "my_secret_key"
}
```

### rabbitmq
类型: Object
```json
    "rabbitmq": {
        "address": "<RabbitMQ Address>",         // RabbitMQ集群地址，例如 "rabbitmq.example.com:5672"
        "username": "<RabbitMQ Username>",       // 用于连接RabbitMQ的用户名，例如 "guest"
        "password": "<RabbitMQ Password>",       // 用于连接RabbitMQ的密码，例如 "guest"
        "v_host": "<RabbitMQ Virtual Host>",     // RabbitMQ的虚拟主机名称，例如 "/" 
        "queue": "<RabbitMQ Queue>",             // 要消费的队列名称，例如 "my_queue"
        "exchange": "<RabbitMQ Exchange>",       // 交换机名称，例如 "my_exchange"
        "routing_key": "<RabbitMQ Routing Key>"  // 路由键，用于消息路由，例如 "my.routing.key"
    }
```

### prom_metrics
类型: Object
```json
    "prom_metrics": {
        "address": "<Prometheus Address>",  // Prometheus的地址，例如 "http://prometheus.example.com:9090"
        "interval": <Interval in seconds>    // 查询的时间间隔（秒），例如 30
    }
```

### pulsar
类型: Object
```json
    "pulsar": {
        "address": "<Pulsar Address>",                    // Pulsar集群地址，例如 "pulsar://pulsar.example.com:6650"
        "topic": "<Pulsar Topic>",                        // 要消费的主题名称，例如 "persistent://public/default/my-topic"
        "subscription_name": "<Pulsar Subscription Name>" // 消费者订阅名称，例如 "my-subscription"
    }
```

### redis
类型: Object
```json
    "redis": {
        "db_num": <DB Number>,                             // Redis数据库编号（0-15），例如 0
        "key_or_channel_name": "<Redis Key or Channel Name>", // Redis的键或频道名称，例如 "my_channel"
        "address": "<Redis Hosts>",                         // Redis主机地址，例如 "redis.example.com:6379"
        "username": "<Redis Username>",                     // Redis用户名（6.0+），例如 "default"
        "password": "<Redis Password>",                     // Redis密码，例如 "mypassword"
        "data_type": "<Data Type>"                          // 数据类型，例如 "lpop"、"rpop"或"subscribe"
    }
```