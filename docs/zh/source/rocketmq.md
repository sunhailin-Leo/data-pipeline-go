RocketMQ
===
# RocketMQ Source 配置文档
此配置定义了一个 RocketMQ source，用于从 RocketMQ 主题中消费消息。以下是 RocketMQ source 的配置示例：
```json
{
  "source": [
    {
      "type": "RocketMQ",
      "source_name": "rocketmq-1",
      "rocketmq": {
        "address": "<RocketMQ Hosts>",
        "group": "<RocketMQ Group>",
        "topic": "<RocketMQ Topic>",
        "access_key": "<RocketMQ Access Key>",
        "secret_key": "<RocketMQ Secret Key>"
      }
    }
  ]
}
```

## 配置参数说明
* **type：**
指定 source 的类型，这里是 "RocketMQ"，表示使用 RocketMQ 作为数据源。

* **source_name：**
RocketMQ source 的自定义名称，可以用来识别该 source。在本例中，它被设置为 "rocketmq-1"。

* **rocketmq：**
包含 RocketMQ source 的核心配置信息：

  * **address：**

    RocketMQ NameServer 或 broker 集群的地址，通常格式为 "host:port"。例如，"rocketmq.example.com:9876" 指定了 RocketMQ 集群的地址。

  * **group：**

    RocketMQ 消费者组名称。不同的消费者组可以独立地消费同一个主题，组内的消费者共享消息消费的负载。 示例："my_consumer_group"。

  * **topic：**

    RocketMQ 主题名称。定义了要消费的消息所在的主题。 示例："my_topic"。

  * **access_key：**

    用于访问 RocketMQ 的访问密钥，通常用于身份验证和授权。此密钥可以由 RocketMQ 管理平台生成。

  * **secret_key：**

    用于与 access_key 配合使用的安全密钥，确保消息传输的安全性和权限控制。

## 配置示例
```json
{
  "source": [
    {
      "type": "RocketMQ",
      "source_name": "rocketmq-1",
      "rocketmq": {
        "address": "rocketmq.example.com:9876",
        "group": "my_consumer_group",
        "topic": "my_topic",
        "access_key": "my_access_key",
        "secret_key": "my_secret_key"
      }
    }
  ]
}
```

在该示例中：
* "rocketmq-1" 是定义的 RocketMQ source 名称。
* "address" 指定了 RocketMQ 集群的地址 "rocketmq.example.com:9876"。
* "group" 定义了消费者组 "my_consumer_group"，负责处理从指定主题接收到的消息。
* "topic" 指定了要消费的 RocketMQ 主题 "my_topic"。
* "access_key" 和 "secret_key" 是用于访问和保护 RocketMQ 服务的凭证。