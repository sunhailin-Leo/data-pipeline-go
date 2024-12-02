RabbitMQ
====

## 概述

`RabbitMQSinkConfig` 结构体用于配置 RabbitMQ Sink，包括与 RabbitMQ 连接和消息传输相关的所有设置，以便将数据发送到指定的 RabbitMQ 队列。
## RabbitMQSinkConfig 结构体

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
## 字段说明

| 字段名             | 类型       | 描述                            |
|-----------------|----------|-------------------------------|
| **Address**     | `string` | RabbitMQ 地址。                  |
| **Username**    | `string` | RabbitMQ 用户名。                 |
| **Password**    | `string` | RabbitMQ 密码。                  |
| **VHost**       | `string` | RabbitMQ 虚拟主机。                |
| **Queue**       | `string` | RabbitMQ 队列名称。                |
| **Exchange**    | `string` | RabbitMQ 交换机名称。               |
| **RoutingKey**  | `string` | RabbitMQ 路由键。                 |
| **MessageMode** | `string` | 消息模式，支持 `json` 和 `text` 两种格式。 |


## 示例
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