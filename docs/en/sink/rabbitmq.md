```
RabbitMQ
====

## Overview

The `RabbitMQSinkConfig` struct is used to configure RabbitMQ Sink, including all settings related to RabbitMQ connection and message transmission, to send data to the specified RabbitMQ queue.

## RabbitMQSinkConfig Struct

```golang
type RabbitMQSinkConfig struct {
	Address     string `json:"address"`      // RabbitMQ address
	Username    string `json:"username"`     // RabbitMQ username
	Password    string `json:"password"`     // RabbitMQ password
	VHost       string `json:"v_host"`       // RabbitMQ vhost
	Queue       string `json:"queue"`        // RabbitMQ queue
	Exchange    string `json:"exchange"`     // RabbitMQ exchange
	RoutingKey  string `json:"routing_key"`  // RabbitMQ routing key
	MessageMode string `json:"message_mode"` // message mode: json, text
}
```

## Field Description

| Field Name       | Type       | Description                                      |
|------------------|------------|--------------------------------------------------|
| **Address**      | `string`   | RabbitMQ address.                                |
| **Username**     | `string`   | RabbitMQ username.                               |
| **Password**     | `string`   | RabbitMQ password.                               |
| **VHost**        | `string`   | RabbitMQ virtual host.                           |
| **Queue**        | `string`   | RabbitMQ queue name.                             |
| **Exchange**     | `string`   | RabbitMQ exchange name.                          |
| **RoutingKey**   | `string`   | RabbitMQ routing key.                            |
| **MessageMode**  | `string`   | Message mode, supports `json` and `text` formats. |

## Example
```json
{
  "sink": [
    {
      "type": "rabbitmq",
      "sink_name": "rabbitmq_data_sink",
      "rabbitmq": {
        "address": "localhost:5672",
        "username": "user",
        "password": "password",
        "v_host": "/",
        "queue": "example_queue",
        "exchange": "example_exchange",
        "routing_key": "example_routing_key",
        "message_mode": "json"
      }
    }
  ]
}
```
```
