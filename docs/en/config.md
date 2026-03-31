Config
===
In this project, the configuration file is the core component of the entire system. It not only determines how the tool operates but also allows users to customize settings according to their needs, ensuring precise and efficient data synchronization. While users can flexibly customize their synchronization requirements, proper configuration is crucial to fully leverage the tool's potential. Below, I will provide a detailed explanation of how to set up the configuration file.

## Environment Variables

The following environment variables need to be set to start data-pipeline-go:

| Environment Variable | Required | Description | Example |
|---------------------|----------|-------------|---------|
| `CONFIG_SRC` | Yes | Configuration source type | `local` |
| `LOCAL_PATH` | No | Local configuration file path (required when `CONFIG_SRC=local`) | `example/kafka_to_http.json` |
| `APOLLO_APP_ID` | No | Apollo application ID (required when `CONFIG_SRC=apollo`) | `my-app` |
| `APOLLO_CLUSTER` | No | Apollo cluster name | `default` |
| `APOLLO_IP` | No | Apollo service address | `http://apollo.example.com:8080` |
| `APOLLO_NAMESPACE` | No | Apollo namespace | `application` |
| `APOLLO_SECRET` | No | Apollo secret key | - |
| `REDIS_ADDR` | No | Redis address (required when `CONFIG_SRC=redis`) | `redis.example.com:6379` |
| `REDIS_PASSWORD` | No | Redis password | - |
| `REDIS_DB` | No | Redis database number | `0` |
| `REDIS_KEY` | No | Redis configuration key | `data-pipeline-config` |
| `NACOS_IP` | No | Nacos service address (required when `CONFIG_SRC=nacos`) | `nacos.example.com` |
| `NACOS_PORT` | No | Nacos port | `8848` |
| `NACOS_GROUP` | No | Nacos group | `DEFAULT_GROUP` |
| `NACOS_DATA_ID` | No | Nacos Data ID | `data-pipeline` |
| `NACOS_NAMESPACE_ID` | No | Nacos namespace ID | - |
| `ZK_HOSTS` | No | Zookeeper address (required when `CONFIG_SRC=zookeeper`) | `zk1:2181,zk2:2181` |
| `ZK_PATH` | No | Zookeeper configuration path | `/data-pipeline/config` |
| `HTTP_URL` | No | HTTP configuration address (required when `CONFIG_SRC=http-get`) | `http://config.example.com/api/config` |
| `HTTP_HEARTBEAT_URL` | No | HTTP heartbeat address | `http://config.example.com/api/heartbeat` |
| `HTTP_HEARTBEAT_SECS` | No | HTTP heartbeat interval (seconds) | `30` |

## Configuration Center Support

| Configuration Source | `CONFIG_SRC` Value | Description |
|---------------------|-------------------|-------------|
| Local File | `local` | Load configuration from local JSON file |
| Apollo | `apollo` | Load from Ctrip Apollo configuration center |
| Redis | `redis` | Load configuration from Redis |
| Nacos | `nacos` | Load from Alibaba Nacos configuration center |
| Zookeeper | `zookeeper` | Load configuration from Zookeeper |
| HTTP-Get | `http-get` | Get configuration via HTTP GET request, supports heartbeat keepalive |

## Configuration File Structure

The configuration file is in JSON format with the following top-level structure:

```json
{
  "streams": [
    {
      "name": "stream-1",
      "enable": true,
      "channel_size": 1000,
      "ack_mode": 0,
      "source": [...],
      "transform": {...},
      "sink": [...]
    }
  ]
}
```

### streams Field Description

| Field Name | Type | Required | Description |
|------------|------|----------|-------------|
| `name` | `string` | Yes | Data stream name, must be unique |
| `enable` | `bool` | No | Whether to enable this data stream, default `false` |
| `channel_size` | `int` | No | Go Channel buffer size, default `0` (unbuffered) |
| `ack_mode` | `int` | No | Message acknowledgment mode, only effective for MQ-type Sources (see table below) |
| `source` | `array` | Yes | Data source configuration array, see [Source](source) |
| `transform` | `object` | Yes | Data transformation configuration, see [Transform](transform) |
| `sink` | `array` | Yes | Data output configuration array, see [Sink](sink) |

### ACK Mode Description

| Value | Mode | Description |
|-------|------|-------------|
| `0` | Acknowledge after consumption | Message is acknowledged immediately after being consumed from Source (default) |
| `1` | Acknowledge after transformation | Message is acknowledged after being processed by Transform |
| `2` | Acknowledge after writing | Message is acknowledged after being successfully written to Sink (safest, but lowest throughput) |

> **Recommendation**: For scenarios with high data consistency requirements, it is recommended to use `ack_mode: 2`; for scenarios that allow minor data loss but pursue high throughput, use `ack_mode: 0`.
