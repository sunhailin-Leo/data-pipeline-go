```
RocketMQ
====

## Overview

The `RocketMQSinkConfig` struct is used to configure RocketMQ Sink, including all settings related to RocketMQ connection and message transmission, to send data to the specified RocketMQ Topic.

## RocketMQSinkConfig Struct

```golang   
type RocketMQSinkConfig struct {
	Address     string `json:"address"`      // RocketMQ address
	Topic       string `json:"topic"`        // RocketMQ Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## Field Description

| Field Name       | Type       | Description                                      |
|------------------|------------|--------------------------------------------------|
| **Address**      | `string`   | RocketMQ address.                                |
| **Topic**        | `string`   | RocketMQ topic name.                             |
| **MessageMode**  | `string`   | Message mode, supports `json` and `text` formats. |

## Example

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

```
