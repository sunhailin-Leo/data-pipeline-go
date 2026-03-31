```
Redis
===
## Redis Source Configuration Documentation
This configuration defines a Redis source, used to fetch data from a Redis database or subscribe to channel messages. Here is the configuration example for a Redis source:
```json
{
  "source": [
    {
      "type": "Redis",
      "source_name": "redis-1",
      "redis": {
        "db_num": <DB Number>,
        "key_or_channel_name": "<Redis Key or Channel Name>",
        "address": "<Redis Hosts>",
        "username": "<Redis Username>",
        "password": "<Redis Password>",
        "data_type": "<Data Type>"
      }
    }
  ]
}
```

## Configuration Parameter Description
* **type：**
Specifies the type of the source, here it is "Redis", indicating that Redis is used as the data source.

* **source_name：**
A custom name for the Redis source, which can be used to identify this source. In this example, it is set to "redis-1".

* **redis：**
Contains the core configuration information for the Redis source:

  * **db_num：**

    The Redis database number. Redis has 16 databases by default (numbered from 0 to 15). This parameter is used to specify which database to read data from. Example: 0 indicates the first Redis database.

  * **key_or_channel_name：**

    The Redis key or channel name. If using Redis as a publish/subscribe system, this parameter represents the channel name to subscribe to; otherwise, it represents the key name in Redis. Example: "my_channel".

  * **address：**

    The address of the Redis server, typically in "host:port" format. For example, "redis.example.com:6379" is the address of the Redis instance.

  * **username：**

    The Redis username (supported in version 6.0 and above), used for authentication. The default username is typically "default".

  * **password：**

    The password used to connect to Redis. This password is used to authenticate with Redis.

  * **data_type：**

    Specifies the data type, defining how to fetch data from Redis. Currently supported types include:

    * "lpop": Pop data from the left side of a Redis list.
    * "rpop": Pop data from the right side of a Redis list.
    * "subscribe": Subscribe to a Redis channel and receive real-time messages from the channel.

## Configuration Example
```json
{
  "source": [
    {
      "type": "Redis",
      "source_name": "redis-1",
      "redis": {
        "db_num": 0,
        "key_or_channel_name": "my_channel",
        "address": "redis.example.com:6379",
        "username": "default",
        "password": "mypassword",
        "data_type": "subscribe"
      }
    }
  ]
}
```

In this example:

* "redis-1" is the defined Redis source name.

* "db_num" is set to 0, indicating the use of the first Redis database.

* "key_or_channel_name" is "my_channel", specifying the name of the Redis channel.

* "address" specifies the Redis server address as "redis.example.com:6379".

* "username" uses the default username "default".

* "password" is set to "mypassword".

* "data_type" is set to "subscribe", indicating subscription to Redis channel messages.
```
