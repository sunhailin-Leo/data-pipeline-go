Redis
===
## Redis Source 配置文档
此配置定义了一个 Redis source，用于从 Redis 数据库中获取数据或订阅频道消息。以下是 Redis source 的配置示例：
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

## 配置参数说明
* **type：**
指定 source 的类型，这里是 "Redis"，表示使用 Redis 作为数据源。

* **source_name：**
Redis source 的自定义名称，可以用来识别该 source。在本例中，它被设置为 "redis-1"。

* **redis：**
包含 Redis source 的核心配置信息：

  * **db_num：**

    Redis 数据库编号。Redis 默认有 16 个数据库（编号从 0 到 15）。该参数用于指定从哪个数据库读取数据。 示例：0 表示 Redis 的第一个数据库。

  * **key_or_channel_name：**

    Redis 的键或频道名称。如果使用 Redis 作为发布/订阅系统，该参数表示要订阅的频道名称；否则，表示 Redis 中的键名称。 示例："my_channel"。

  * **address：**

    Redis 服务器的地址，通常为 "host:port" 格式。例如，"redis.example.com:6379" 是 Redis 实例的地址。

  * **username：**

    Redis 用户名（6.0 及以上版本支持），用于身份验证。默认用户名通常为 "default"。

  * **password：**

    用于连接 Redis 的密码。通过此密码可以对 Redis 进行身份验证。

  * **data_type：**

    指定数据类型，定义从 Redis 中获取数据的方式。目前支持的类型包括：

    * "lpop"：从 Redis 列表的左侧弹出数据。
    * "rpop"：从 Redis 列表的右侧弹出数据。
    * "subscribe"：订阅 Redis 频道，接收来自频道的实时消息。

## 配置示例
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

在该示例中：

* "redis-1" 是定义的 Redis source 名称。

* "db_num" 设置为 0，表示使用 Redis 的第一个数据库。

* "key_or_channel_name" 为 "my_channel"，指定了 Redis 频道的名称。

* "address" 指定了 Redis 服务器的地址为 "redis.example.com:6379"。

* "username" 使用了默认用户名 "default"。

* "password" 设置为 "mypassword"。

* "data_type" 设置为 "subscribe"，表示订阅 Redis 频道消息。