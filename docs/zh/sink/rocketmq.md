RocketMQ
====

# RocketMQSinkConfig 结构体

## 概述

`RocketMQSinkConfig` 结构体用于配置 RocketMQ Sink，包括与 RocketMQ 连接和消息传输相关的所有设置，以便将数据发送到指定的 RocketMQ Topic 中。

## 结构体定义

```go
type RocketMQSinkConfig struct {
	Address     string `json:"address"`      // RocketMQ address
	Topic       string `json:"topic"`        // RocketMQ Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## 字段说明

| 字段名             | 类型       | 描述                            |
|-----------------|----------|-------------------------------|
| **Address**     | `string` | RocketMQ 地址。                  |
| **Topic**       | `string` | RocketMQ 的主题名称。               |
| **MessageMode** | `string` | 消息模式，支持 `json` 和 `text` 两种格式。 |

## 示例

```json
{
  "sink": [
    {
      "type": "rocketmq",
      "sink_name": "rocketmq_data_sink",
      "rocketmq": {
        "address": "localhost:9876",
        "topic": "data_topic",
        "message_mode": "json"
      }
    }
  ]
}
```
