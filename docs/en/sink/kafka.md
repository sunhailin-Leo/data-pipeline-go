```
Kafka
=======
## Overview
The `KafkaSinkConfig` struct is used to configure Kafka Sink, including all settings related to Kafka message transmission, to send data to Kafka topics using the specified message format.

## KafkaSinkConfig Struct

```golang
type KafkaSinkConfig struct {
	Address     string `json:"address"`      // Kafka address
	Topic       string `json:"topic"`        // Kafka Topic
	MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## Field Description

| Field Name       | Type       | Description                            |
|------------------|------------|----------------------------------------|
| **Address**      | `string`   | Kafka address, usually in the format `host:port`. |
| **Topic**        | `string`   | Target Kafka topic name.              |
| **MessageMode**  | `string`   | Message mode, supports `json` and `text` formats. |

## Example
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

```
