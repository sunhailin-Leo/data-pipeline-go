RabbitMQ
===
## RabbitMQ Source 配置
此配置定义了一个 RabbitMQ source，用于从 RabbitMQ 队列中消费消息。以下是 RabbitMQ source 的配置示例：
```json
{
  "source": [
    {
      "type": "RabbitMQ",
      "source_name": "rabbitmq-1",
      "rabbitmq": {
        "address": "<RabbitMQ Address>",
        "username": "<RabbitMQ Username>",
        "password": "<RabbitMQ Password>",
        "v_host": "<RabbitMQ Virtual Host>",
        "queue": "<RabbitMQ Queue>",
        "exchange": "<RabbitMQ Exchange>",
        "routing_key": "<RabbitMQ Routing Key>"
      }
    }
  ]
}
```

## 配置参数说明
**type：**

指定 source 的类型，这里是 "RabbitMQ"，表示使用 RabbitMQ 作为数据源。

**source_name：**

RabbitMQ source 的自定义名称，可以用来识别该 source。在本例中，它被设置为 "rabbitmq-1"。

**rabbitmq：**

包含 RabbitMQ source 的核心配置信息：

* **address：**

    RabbitMQ broker 的地址，通常格式为 "host:port"。例如，"rabbitmq.example.com:5672" 指定了 RabbitMQ 的连接地址。

* **username：**

    用于连接 RabbitMQ 的用户名。例如，"guest" 是默认的用户名。

* **password：**

    用于连接 RabbitMQ 的密码。例如，"guest" 是默认的密码。

* **v_host：**
    
    RabbitMQ 的虚拟主机名称，用于隔离不同的应用程序。默认值通常为 "/"。

* **queue：**

    要消费的 RabbitMQ 队列名称。例如，"my_queue" 是要从中读取消息的队列名称。

* **exchange：**

    RabbitMQ 交换机的名称，用于将消息发布到队列。例如，"my_exchange"。

* **routing_key：**

    RabbitMQ 路由键，用于决定消息应该被路由到哪个队列。路由键与交换机一起决定消息流动的方向。例如，"my.routing.key"。

## 配置示例
```json
{
  "source": [
    {
      "type": "RabbitMQ",
      "source_name": "rabbitmq-1",
      "rabbitmq": {
        "address": "rabbitmq.example.com:5672",
        "username": "guest",
        "password": "guest",
        "v_host": "/",
        "queue": "my_queue",
        "exchange": "my_exchange",
        "routing_key": "my.routing.key"
      }
    }
  ]
}
```

在该示例中：

* "rabbitmq-1" 是定义的 RabbitMQ source 名称。

* "address" 指定了 RabbitMQ broker 的地址为 "rabbitmq.example.com:5672"。

* "username" 和 "password" 使用了默认的 "guest" 凭证。

* "v_host" 为默认虚拟主机 "/"。

* "queue" 指定了要消费的队列 "my_queue"。

* "exchange" 定义了交换机 "my_exchange"。

* "routing_key" 指定了路由键 "my.routing.key"，用于将消息路由到适当的队列。