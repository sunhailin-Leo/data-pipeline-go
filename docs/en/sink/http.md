```
HTTP
====
## Overview

The `HTTPSinkConfig` struct is used to configure HTTP Sink, including all settings related to HTTP requests, to send data to the specified URL via HTTP protocol.

## HTTPSinkConfig Struct

```golang
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

## Field Description
| Field Name                      | Type               | Description                                        |
|---------------------------------|--------------------|---------------------------------------------------|
| **URL**                         | `string`           | Target URL for data sending.                      |
| **ContentType**                 | `string`           | HTTP request content type, such as `application/json`. |
| **ReadTimeoutSecs**             | `int`              | Timeout for reading response (in seconds).        |
| **WriteTimeoutSecs**            | `int`              | Timeout for writing request (in seconds).         |
| **MaxIdleConnDurationSecs**     | `int`              | Maximum idle connection duration (in seconds).    |
| **MaxConnWaitTimeoutSecs**      | `int`              | Maximum connection wait timeout (in seconds).     |
| **Headers**                     | `map[string]string` | Custom HTTP request header key-value pairs.      |

## Example
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
```
