NSQ
===

## NSQ Source Configuration

NSQ Source consumes messages from an NSQ topic and passes each message body into Transform.

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

## Fields

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `type` | `string` | Yes | Must be `NSQ`. |
| `source_name` | `string` | Yes | Source alias referenced by `source_name` in Transform schemas. |
| `nsq.address` | `string` | No | nsqd TCP address, comma-separated; used when `lookupd_address` is empty. |
| `nsq.lookupd_address` | `string` | No | nsqlookupd HTTP address, comma-separated; takes precedence when non-empty. |
| `nsq.topic` | `string` | Yes | NSQ topic to consume. |
| `nsq.channel` | `string` | Yes | NSQ channel. |
| `nsq.max_in_flight` | `int` | No | Max in-flight messages; empty means go-nsq default. |

## ACK

NSQ messages are acknowledged with `Finish()`. With `ack_mode: 0`, the source confirms after consumption. With `ack_mode: 1` or `2`, confirmation is delayed to Transform or Sink.
