Config
===
在本项目中，配置文件是整个系统的核心部分。它不仅决定了工具如何运行，还允许用户根据自身的需求进行个性化设置，确保数据同步过程精确且高效。尽管用户可以灵活地自定义他们的同步需求，但要充分发挥工具的潜力，正确配置文件至关重要。接下来，我将详细说明如何设置配置文件。

## 环境变量

启动 data-pipeline-go 需要设置以下环境变量：

| 环境变量 | 必填 | 说明 | 示例 |
|---------|------|------|------|
| `CONFIG_SRC` | 是 | 配置来源类型 | `local` |
| `LOCAL_PATH` | 否 | 本地配置文件路径（`CONFIG_SRC=local` 时必填） | `example/kafka_to_http.json` |
| `APOLLO_APP_ID` | 否 | Apollo 应用 ID（`CONFIG_SRC=apollo` 时必填） | `my-app` |
| `APOLLO_CLUSTER` | 否 | Apollo 集群名称 | `default` |
| `APOLLO_IP` | 否 | Apollo 服务地址 | `http://apollo.example.com:8080` |
| `APOLLO_NAMESPACE` | 否 | Apollo 命名空间 | `application` |
| `APOLLO_SECRET` | 否 | Apollo 密钥 | - |
| `REDIS_ADDR` | 否 | Redis 地址（`CONFIG_SRC=redis` 时必填） | `redis.example.com:6379` |
| `REDIS_PASSWORD` | 否 | Redis 密码 | - |
| `REDIS_DB` | 否 | Redis 数据库编号 | `0` |
| `REDIS_KEY` | 否 | Redis 配置 Key | `data-pipeline-config` |
| `NACOS_IP` | 否 | Nacos 服务地址（`CONFIG_SRC=nacos` 时必填） | `nacos.example.com` |
| `NACOS_PORT` | 否 | Nacos 端口 | `8848` |
| `NACOS_GROUP` | 否 | Nacos 分组 | `DEFAULT_GROUP` |
| `NACOS_DATA_ID` | 否 | Nacos Data ID | `data-pipeline` |
| `NACOS_NAMESPACE_ID` | 否 | Nacos 命名空间 ID | - |
| `ZK_HOSTS` | 否 | Zookeeper 地址（`CONFIG_SRC=zookeeper` 时必填） | `zk1:2181,zk2:2181` |
| `ZK_PATH` | 否 | Zookeeper 配置路径 | `/data-pipeline/config` |
| `HTTP_URL` | 否 | HTTP 配置地址（`CONFIG_SRC=http-get` 时必填） | `http://config.example.com/api/config` |
| `HTTP_HEARTBEAT_URL` | 否 | HTTP 心跳地址 | `http://config.example.com/api/heartbeat` |
| `HTTP_HEARTBEAT_SECS` | 否 | HTTP 心跳间隔（秒） | `30` |

## 配置中心支持

| 配置来源 | `CONFIG_SRC` 值 | 说明 |
|---------|----------------|------|
| 本地文件 | `local` | 从本地 JSON 文件加载配置 |
| Apollo | `apollo` | 从携程 Apollo 配置中心加载 |
| Redis | `redis` | 从 Redis 中加载配置 |
| Nacos | `nacos` | 从阿里 Nacos 配置中心加载 |
| Zookeeper | `zookeeper` | 从 Zookeeper 加载配置 |
| HTTP-Get | `http-get` | 通过 HTTP GET 请求获取配置，支持心跳保活 |

## 配置文件结构

配置文件为 JSON 格式，顶层结构如下：

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

### streams 字段说明

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `name` | `string` | 是 | 数据流名称，必须唯一 |
| `enable` | `bool` | 否 | 是否启用该数据流，默认 `false` |
| `channel_size` | `int` | 否 | Go Channel 缓冲区大小，默认 `0`（无缓冲） |
| `ack_mode` | `int` | 否 | 消息确认模式，仅对 MQ 类 Source 生效（见下表） |
| `source` | `array` | 是 | 数据源配置数组，详见 [Source](source) |
| `transform` | `object` | 是 | 数据转换配置，详见 [Transform](transform) |
| `sink` | `array` | 是 | 数据输出配置数组，详见 [Sink](sink) |

### ACK 模式说明

| 值 | 模式 | 说明 |
|----|------|------|
| `0` | 消费后确认 | 消息从 Source 消费后立即确认（默认） |
| `1` | 转换后确认 | 消息经过 Transform 处理后确认 |
| `2` | 写入后确认 | 消息成功写入 Sink 后确认（最安全，但吞吐量最低） |

> **建议**：对于数据一致性要求高的场景，推荐使用 `ack_mode: 2`；对于允许少量数据丢失但追求高吞吐的场景，使用 `ack_mode: 0`。


