```
Kafka
====

## Kafka Source Configuration

This configuration defines a Kafka source, used to read data from Kafka topics.

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
## Configuration Parameter Description:

**type：**

Specifies the type of the source, here it is "Kafka", indicating that Kafka is used as the data source.

**source_name：**

A custom name for the Kafka source, which can be used to identify this source. In this example, it is set to "kafka-1".

**kafka：**

Contains the core configuration information for the Kafka source.

* **address：**

    The list of Kafka broker addresses, typically in the format "host1:port1,host2:port2", with multiple addresses separated by commas.

* **group：**

    The Kafka consumer group name. Consumers within the same group will share the load of reading from this topic.

* **topic：**

    The name of the Kafka topic to consume. One or more topics can be specified.

## Configuration Example

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
In this example:
* "kafka-1" is the defined Kafka source name.

* "address" specifies the Kafka cluster address "kafka-broker-1:9092,kafka-broker-2:9092".

* "group" defines the consumer group "my_consumer_group", used to share the consumption load of the Kafka topic.

* "topic" specifies the Kafka topic "my_topic" to be consumed.
```
