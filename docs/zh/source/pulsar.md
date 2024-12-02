Pulsar
===
## Pulsar Source 配置
此配置定义了一个 Pulsar source，用于从 Pulsar 主题中读取数据。以下是 Pulsar source 的配置示例：

```json
{
  "source": [
    {
      "type": "Pulsar",
      "source_name": "Pulsar-1",
      "pulsar": {
        "address": "<Pulsar Address>",                    // Pulsar集群地址，例如 "pulsar://pulsar.example.com:6650"
        "topic": "<Pulsar Topic>",                        // 要消费的主题名称，例如 "persistent://public/default/my-topic"
        "subscription_name": "<Pulsar Subscription Name>" // 消费者订阅名称，例如 "my-subscription"
      }
    }
  ]
}
```
配置参数说明

**type：**

指定 source 的类型，这里是 "Pulsar"，表示使用 Apache Pulsar 作为数据源。

**source_name：**

Pulsar source 的自定义名称，可以用来识别该 source。在本例中，它被设置为 "Pulsar-1"。

**pulsar：**

包含 Pulsar source 的核心配置信息：

* **address：**

    Pulsar broker 或集群的地址，格式通常为 "pulsar://host:port"。例如，"pulsar://pulsar.example.com:6650" 是 Pulsar broker 的地址。

* **topic：**

    要消费的 Pulsar 主题名称。Pulsar 使用命名空间和租户结构，主题的格式通常为 "persistent://tenant/namespace/topic"。

* **subscription_name：**

    消费者订阅的名称，用于管理消费者组，确保消息仅被相应组内的消费者读取。Pulsar 允许定义多种订阅模式（如独占、共享、失败回溯等），通过这个名称来区分不同的订阅。


## 配置示例

```json
{
  "source": [
    {
      "type": "Pulsar",
      "source_name": "Pulsar-1",
      "pulsar": {
        "address": "pulsar://pulsar.example.com:6650",
        "topic": "persistent://public/default/my-topic",
        "subscription_name": "my-subscription"
      }
    }
  ]
}
```
在该示例中：

* "Pulsar-1" 是定义的 Pulsar source 名称。

* "address" 指定了 Pulsar broker 的地址 "pulsar://pulsar.example.com:6650"。

* "topic" 指定了要消费的 Pulsar 主题 "persistent://public/default/my-topic"。

* "subscription_name" 定义了订阅名称 "my-subscription"，用于标识消费者订阅。