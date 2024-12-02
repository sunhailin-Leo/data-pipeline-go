Pulsar
====

## 概述

`PulsarSinkConfig` 结构体用于配置 Pulsar Sink，包括与 Pulsar 连接和数据传输相关的所有设置，以便将数据发送到指定的 Pulsar Topic。
## PulsarSinkConfig 结构体

```golang
type PulsarSinkConfig struct {
Address     string `json:"address"`      // Pulsar host
Topic       string `json:"topic"`        // Pulsar Topic
MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## 字段说明

| 字段名             | 类型       | 描述                          |
|-----------------|----------|-----------------------------|
| **Address**     | `string` | Pulsar 主机地址。                |
| **Topic**       | `string` | Pulsar 主题。                  |
| **MessageMode** | `string` | 消息模式，支持 `json` 和 `text` 格式。 |

## 示例
```json
{
  "sink": [
    {
      "type": "pulsar",
      "sink_name": "pulsar_data_sink",
      "pulsar": {
        "address": "pulsar://localhost:6650",
        "topic": "persistent://public/default/example-topic",
        "message_mode": "json"
      }
    }
  ]
}

```