Kafka
=======
## 概述
`KafkaSinkConfig` 结构体用于配置 Kafka Sink，包括与 Kafka 消息传输相关的所有设置，以便将数据通过指定的消息格式发送到 Kafka 主题。
## KafkaSinkConfig 结构体

```golang
type KafkaSinkConfig struct {
	Address     string `json:"address"`      // Kafka address
	Topic       string `json:"topic"`        // Kafka Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}
```
## 字段说明

| 字段名             | 类型       | 描述                            |
|-----------------|----------|-------------------------------|
| **Address**     | `string` | Kafka 的地址，通常格式为 `host:port`。  |
| **Topic**       | `string` | 目标 Kafka 主题名称。                |
| **MessageMode** | `string` | 消息模式，支持 `json` 和 `text` 两种格式。 |

## 示例
```json
{
  "sink": [
    {
      "type": "kafka",
      "sink_name": "kafka_data_sink",
      "kafka": {
        "address": "localhost:9092",
        "topic": "data_topic",
        "message_mode": "json"
      }
    }
  ]
}

```
