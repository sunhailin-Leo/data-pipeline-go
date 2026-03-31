```
Redis
====

## Overview

The `RedisSinkConfig` struct is used to configure Redis Sink, containing all settings related to Redis data transmission, to send data to the specified key or channel in the Redis database.

## RedisSinkConfig Struct

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

## Field Description

| Field Name               | Type               | Description                                                                                   |
|--------------------------|--------------------|-----------------------------------------------------------------------------------------------|
| **DBNum**                | `int`             | Redis database number (0-15), invalid in cluster mode.                                        |
| **Expire**               | `int`             | Expiration time, unit is seconds.                                                             |
| **Address**              | `string`          | Redis server address.                                                                         |
| **Username**             | `string`          | Redis username (applicable for version 6.0+).                                                |
| **Password**             | `string`          | Redis password.                                                                               |
| **DataType**             | `string`          | Data type, supports `kv`, `hash`, `lpush`, `rpush`, `set`, `publish`, etc.                   |
| **KeyOrChannelName**     | `string`          | Key or channel name; if you want the key to be variable, please use `fromTransformHead`.      |

## Example
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
```
