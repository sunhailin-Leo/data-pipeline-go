NSQ
===

## NSQ Source 配置

NSQ Source 用于从 NSQ topic 中消费消息，并将消息体传入 Transform。

```json
{
  "source": [
    {
      "type": "NSQ",
      "source_name": "nsq-1",
      "nsq": {
        "address": "127.0.0.1:4150",
        "lookupd_address": "",
        "topic": "data-pipeline-topic",
        "channel": "data-pipeline-channel",
        "max_in_flight": 100
      }
    }
  ]
}
```

## 字段说明

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `type` | `string` | 是 | 固定为 `NSQ`。 |
| `source_name` | `string` | 是 | Source 别名，在 Transform schema 中通过 `source_name` 引用。 |
| `nsq.address` | `string` | 否 | nsqd TCP 地址，支持逗号分隔；未配置 `lookupd_address` 时使用。 |
| `nsq.lookupd_address` | `string` | 否 | nsqlookupd HTTP 地址，支持逗号分隔；非空时优先使用。 |
| `nsq.topic` | `string` | 是 | 要消费的 NSQ topic。 |
| `nsq.channel` | `string` | 是 | NSQ channel。 |
| `nsq.max_in_flight` | `int` | 否 | 最大 in-flight 消息数，留空时使用 go-nsq 默认值。 |

## ACK

NSQ 消息通过 `Finish()` 确认。`ack_mode` 为 `0` 时在 Source 消费后确认；为 `1` 或 `2` 时分别延后到 Transform 或 Sink 阶段确认。
