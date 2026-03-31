Source
===
## Data Sources
Supported list:
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

## ACK Mode

In Stream configuration, the `ack_mode` field controls the message acknowledgment timing, only effective for MQ-type Sources (Kafka, RocketMQ, RabbitMQ, Pulsar):

| Value | Mode | Description |
|-------|------|-------------|
| `0` | Acknowledge after consumption | Message is acknowledged immediately after being consumed from Source (default) |
| `1` | Acknowledge after transformation | Message is acknowledged after being processed by Transform |
| `2` | Acknowledge after writing | Message is acknowledged after being successfully written to Sink (safest) |

## Field Description
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
Type: Array

Description: The source field is an array of data sources, supporting multiple different types of data sources.

### type
Type: String

Description: Indicates the type of data source, refer to `Supported list`. For Kafka data source, this should be "Kafka". This field is used to specify the data source type so that the system can correctly process and parse it.

### source_name
Type: String

Description: This is the name of this Kafka data source. source_name is used to identify the specific data source, making it convenient for developers to manage multiple data sources in the system. For example, the name can be "kafka-1" to represent an instance of the Kafka data source.

### kafka
Type: Object
```json
    "kafka": {
        "address": "<Kafka Hosts>",           // Kafka cluster address, e.g. "kafka1:9092,kafka2:9092"
        "group": "<Kafka Group>",             // Kafka consumer group name, e.g. "consumer-group-1"
        "topic": "<Kafka Topic>"              // Kafka topic name to consume, e.g. "my_topic"
    }
```
Contains configuration information related to Kafka, including the Kafka cluster address, consumer group, and the topic to consume. Below are the specific fields in the kafka object:

* address: Kafka cluster address, format is <host>:<port>. Can be a single broker address or multiple, separated by commas. For example: "kafka1:9092,kafka2:9092".

* group: Kafka consumer group name, used for multiple consumers to share load within the same group, ensuring each message is processed by only one consumer in the group. For example: "consumer-group-1".

* topic: Kafka topic name, representing message classification. Consumers will fetch messages from the specified topic. This name should match the topic in the Kafka cluster. For example: "my_topic".

### rocketmq
Type: Object
```json
    "rocketmq":{
        "address": "<RocketMQ Hosts>",             // RocketMQ cluster address, e.g. "rocketmq.example.com:9876"
        "group": "<RocketMQ Group>",               // RocketMQ consumer group name, e.g. "my_consumer_group"
        "topic": "<RocketMQ Topic>",               // RocketMQ topic name to consume, e.g. "my_topic"
        "access_key": "<RocketMQ Access Key>",     // Access key for accessing RocketMQ, e.g. "my_access_key"
        "secret_key": "<RocketMQ Secret Key>"      // Secret key for accessing RocketMQ, e.g. "my_secret_key"
}
```

### rabbitmq
Type: Object
```json
    "rabbitmq": {
        "address": "<RabbitMQ Address>",         // RabbitMQ cluster address, e.g. "rabbitmq.example.com:5672"
        "username": "<RabbitMQ Username>",       // Username for connecting to RabbitMQ, e.g. "guest"
        "password": "<RabbitMQ Password>",       // Password for connecting to RabbitMQ, e.g. "guest"
        "v_host": "<RabbitMQ Virtual Host>",     // RabbitMQ virtual host name, e.g. "/" 
        "queue": "<RabbitMQ Queue>",             // Queue name to consume, e.g. "my_queue"
        "exchange": "<RabbitMQ Exchange>",       // Exchange name, e.g. "my_exchange"
        "routing_key": "<RabbitMQ Routing Key>"  // Routing key for message routing, e.g. "my.routing.key"
    }
```

### prom_metrics
Type: Object
```json
    "prom_metrics": {
        "address": "<Prometheus Address>",  // Prometheus address, e.g. "http://prometheus.example.com:9090"
        "interval": <Interval in seconds>    // Query interval in seconds, e.g. 30
    }
```

### pulsar
Type: Object
```json
    "pulsar": {
        "address": "<Pulsar Address>",                    // Pulsar cluster address, e.g. "pulsar://pulsar.example.com:6650"
        "topic": "<Pulsar Topic>",                        // Topic name to consume, e.g. "persistent://public/default/my-topic"
        "subscription_name": "<Pulsar Subscription Name>" // Consumer subscription name, e.g. "my-subscription"
    }
```

### redis
Type: Object
```json
    "redis": {
        "db_num": <DB Number>,                             // Redis database number (0-15), e.g. 0
        "key_or_channel_name": "<Redis Key or Channel Name>", // Redis key or channel name, e.g. "my_channel"
        "address": "<Redis Hosts>",                         // Redis host address, e.g. "redis.example.com:6379"
        "username": "<Redis Username>",                     // Redis username (6.0+), e.g. "default"
        "password": "<Redis Password>",                     // Redis password, e.g. "mypassword"
        "data_type": "<Data Type>"                          // Data type, e.g. "lpop", "rpop" or "subscribe"
    }
```
