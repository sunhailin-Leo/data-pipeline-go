Redis
====

## 概述

`RedisSinkConfig` 结构体用于配置 Redis Sink，包含与 Redis 数据传输相关的所有设置，以便将数据发送到 Redis 数据库中的指定键或频道。
## RedisSinkConfig 结构体

```golang
type RedisSinkConfig struct {
	DBNum            int    `json:"db_num"`              // db number（0-15）,cluster mode is invalid
	Expire           int    `json:"expire"`              // expire time, unit: seconds
	Address          string `json:"address"`             // redis hosts
	Username         string `json:"username"`            // Redis Username(6.0+)
	Password         string `json:"password"`            // Redis password
	DataType         string `json:"data_type"`           // data type: kv, hash, lpush, rpush, set, publish
	KeyOrChannelName string `json:"key_or_channel_name"` // Key or Channel name, if you want key is variable, please use "fromTransformHead"
}
```

## 字段说明

| 字段名                | 类型               | 描述                                                                                   |
|-----------------------|--------------------|----------------------------------------------------------------------------------------|
| **DBNum**             | `int`             | Redis 数据库编号（0-15），在集群模式下无效。                                             |
| **Expire**            | `int`             | 过期时间，单位为秒。                                                                    |
| **Address**           | `string`          | Redis 服务器地址。                                                                      |
| **Username**          | `string`          | Redis 用户名（适用于 6.0+ 版本）。                                                       |
| **Password**          | `string`          | Redis 密码。                                                                            |
| **DataType**          | `string`          | 数据类型，支持 `kv`、`hash`、`lpush`、`rpush`、`set`、`publish` 等。                     |
| **KeyOrChannelName**  | `string`          | 键或频道名称；如需动态生成键，请使用 `fromTransformHead`。                               |

## 示例
```json
{
  "sink": [
    {
      "type": "redis",
      "sink_name": "redis_data_sink",
      "redis": {
        "db_num": 0,
        "expire": 3600,
        "address": "localhost:6379",
        "username": "user",
        "password": "password",
        "data_type": "kv",
        "key_or_channel_name": "example_key"
      }
    }
  ]
}

```