HTTP
====
# HTTPSinkConfig 结构体

## 概述

`HTTPSinkConfig` 结构体用于配置 HTTP Sink，包括与 HTTP 请求相关的所有设置，以便将数据通过 HTTP 协议发送到指定的 URL。

## 结构体定义

```go
type HTTPSinkConfig struct {
	URL                     string            `json:"url"`
	ContentType             string            `json:"content_type"`
	ReadTimeoutSecs         int               `json:"read_timeout_secs"`
	WriteTimeoutSecs        int               `json:"write_timeout_secs"`
	MaxIdleConnDurationSecs int               `json:"max_idle_conn_duration_secs"`
	MaxConnWaitTimeoutSecs  int               `json:"max_conn_wait_timeout_secs"`
	Headers                 map[string]string `json:"headers"`
}
```

## 字段说明
| 字段名                      | 类型               | 描述                                        |
|----------------------------|--------------------|-------------------------------------------|
| **URL**                    | `string`           | 数据发送的目标 URL。                      |
| **ContentType**            | `string`           | HTTP 请求的内容类型，例如 `application/json`。 |
| **ReadTimeoutSecs**        | `int`              | 读取响应的超时时间（秒）。                |
| **WriteTimeoutSecs**       | `int`              | 写入请求的超时时间（秒）。                |
| **MaxIdleConnDurationSecs** | `int`             | 最大空闲连接持续时间（秒）。              |
| **MaxConnWaitTimeoutSecs**  | `int`             | 最大连接等待超时时间（秒）。              |
| **Headers**                | `map[string]string` | 自定义 HTTP 请求头的键值对。             |

## 示例
```json
{
  "sink": [
    {
      "type": "http",
      "sink_name": "api_data_sink",
      "http": {
        "url": "https://api.example.com/data",
        "content_type": "application/json",
        "read_timeout_secs": 10,
        "write_timeout_secs": 10,
        "max_idle_conn_duration_secs": 60,
        "max_conn_wait_timeout_secs": 5,
        "headers": {
          "Authorization": "Bearer token_value",
          "Custom-Header": "CustomValue"
        }
      }
    }
  ]
}

```