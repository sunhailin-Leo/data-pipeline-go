```
Pulsar
====

## Overview

The `PulsarSinkConfig` struct is used to configure Pulsar Sink, including all settings related to Pulsar connection and data transmission, to send data to the specified Pulsar Topic.

## PulsarSinkConfig Struct

```golang
type PulsarSinkConfig struct {
Address     string `json:"address"`      // Pulsar host
Topic       string `json:"topic"`        // Pulsar Topic
MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## Field Description

| Field Name       | Type       | Description                                      |
|------------------|------------|--------------------------------------------------|
| **Address**      | `string`   | Pulsar host address.                             |
| **Topic**        | `string`   | Pulsar topic.                                    |
| **MessageMode**  | `string`   | Message mode, supports `json` and `text` formats. |

## Example
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
```
