```
RabbitMQ
===
## RabbitMQ Source Configuration
This configuration defines a RabbitMQ source, used to consume messages from RabbitMQ queues. Here is the configuration example for a RabbitMQ source:
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

## Configuration Parameter Description
**type：**

Specifies the type of the source, here it is "RabbitMQ", indicating that RabbitMQ is used as the data source.

**source_name：**

A custom name for the RabbitMQ source, which can be used to identify this source. In this example, it is set to "rabbitmq-1".

**rabbitmq：**

Contains the core configuration information for the RabbitMQ source:

* **address：**

    The address of the RabbitMQ broker, typically in the format "host:port". For example, "rabbitmq.example.com:5672" specifies the RabbitMQ connection address.

* **username：**

    The username used to connect to RabbitMQ. For example, "guest" is the default username.

* **password：**

    The password used to connect to RabbitMQ. For example, "guest" is the default password.

* **v_host：**
    
    The RabbitMQ virtual host name, used to isolate different applications. The default value is typically "/".

* **queue：**

    The name of the RabbitMQ queue to consume. For example, "my_queue" is the name of the queue from which to read messages.

* **exchange：**

    The name of the RabbitMQ exchange, used to publish messages to queues. For example, "my_exchange".

* **routing_key：**

    The RabbitMQ routing key, used to determine which queue messages should be routed to. The routing key, together with the exchange, determines the direction of message flow. For example, "my.routing.key".

## Configuration Example
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

In this example:

* "rabbitmq-1" is the defined RabbitMQ source name.

* "address" specifies the RabbitMQ broker address as "rabbitmq.example.com:5672".

* "username" and "password" use the default "guest" credentials.

* "v_host" is the default virtual host "/".

* "queue" specifies the queue to consume "my_queue".

* "exchange" defines the exchange "my_exchange".

* "routing_key" specifies the routing key "my.routing.key", used to route messages to the appropriate queue.
```
