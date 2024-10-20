Kafka
====

# Kafka Source 配置

此配置定义了一个 Kafka source，用于从 Kafka 主题中读取数据。

```json
{
"source": [
    {
    "type": "Kafka",
    "source_name": "kafka-1",
    "kafka": {
        "address": "<Kafka Hosts>",
        "group": "<Kafka Group>",
        "topic": "<Kafka Topic>"
      }
    }
  ]
}
```
## 配置参数说明：

**type：**

指定 source 的类型，这里是 "Kafka"，表示使用 Kafka 作为数据源。

**source_name：**

Kafka source 的自定义名称，可以用来识别该 source。在本例中，它被设置为 "kafka-1"。

**kafka：**

包含 Kafka source 的核心配置信息。

* **address：**

    Kafka broker 的地址列表，格式通常为 "host1:port1,host2:port2"，多个地址用逗号分隔。

* **group：**

    Kafka 消费者组名称。同一个组内的消费者将共享读取该主题的负载。

* **topic：**

    要消费的 Kafka 主题名称。可以指定一个或多个主题。

## 配置示例

```json
{
"source": [
    {
        "type": "Kafka",
        "source_name": "kafka-1",
        "kafka": {
            "address": "kafka-broker-1:9092,kafka-broker-2:9092",
            "group": "consumer-group-1",
            "topic": "topic-1"
        }
    }
  ]
}
```
在该示例中：
* "kafka-1" 是定义的 Kafka source 名称。

* "address" 指定了 Kafka 集群的地址 "kafka-broker-1:9092,kafka-broker-2:9092"。

* "group" 定义了消费者组 "my_consumer_group"，用于共享 Kafka 主题的消费负载。

* "topic" 指定了要消费的 Kafka 主题 "my_topic"。