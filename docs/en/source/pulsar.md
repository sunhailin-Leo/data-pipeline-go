```
Pulsar
===
## Pulsar Source Configuration
This configuration defines a Pulsar source, used to read data from Pulsar topics. Here is the configuration example for a Pulsar source:

```json
{
  "source": [
    {
      "type": "Pulsar",
      "source_name": "Pulsar-1",
      "pulsar": {
        "address": "<Pulsar Address>",                    // Pulsar cluster address, for example "pulsar://pulsar.example.com:6650"
        "topic": "<Pulsar Topic>",                        // Topic name to consume, for example "persistent://public/default/my-topic"
        "subscription_name": "<Pulsar Subscription Name>" // Consumer subscription name, for example "my-subscription"
      }
    }
  ]
}
```
Configuration Parameter Description

**type：**

Specifies the type of the source, here it is "Pulsar", indicating that Apache Pulsar is used as the data source.

**source_name：**

A custom name for the Pulsar source, which can be used to identify this source. In this example, it is set to "Pulsar-1".

**pulsar：**

Contains the core configuration information for the Pulsar source:

* **address：**

    The address of the Pulsar broker or cluster, typically in the format "pulsar://host:port". For example, "pulsar://pulsar.example.com:6650" is the address of the Pulsar broker.

* **topic：**

    The name of the Pulsar topic to consume. Pulsar uses a namespace and tenant structure, and the topic format is typically "persistent://tenant/namespace/topic".

* **subscription_name：**

    The name of the consumer subscription, used to manage consumer groups and ensure messages are only read by consumers in the corresponding group. Pulsar allows defining multiple subscription modes (such as exclusive, shared, failover, etc.), distinguished by this name.

## Configuration Example

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
In this example:

* "Pulsar-1" is the defined Pulsar source name.

* "address" specifies the Pulsar broker address "pulsar://pulsar.example.com:6650".

* "topic" specifies the Pulsar topic to consume "persistent://public/default/my-topic".

* "subscription_name" defines the subscription name "my-subscription", used to identify the consumer subscription.
```
