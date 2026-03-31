<h1 align="center">data-pipeline-go</h1>

<p align="center">
  <em>基于 Go 实现的高性能数据同步工具，类似 SeaTunnel，简便易用</em>
</p>

<p align="center">
  <a href="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/test.yml"><img src="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/test.yml/badge.svg" alt="Test & Coverage"></a>
  <a href="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/lint.yml"><img src="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/lint.yml/badge.svg" alt="Lint Check"></a>
  <a href="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/build.yml"><img src="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/build.yml/badge.svg" alt="Build"></a>
  <a href="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/codeql.yml"><img src="https://github.com/sunhailin-Leo/data-pipeline-go/actions/workflows/codeql.yml/badge.svg" alt="CodeQL"></a>
  <a href="https://codecov.io/gh/sunhailin-Leo/data-pipeline-go"><img src="https://codecov.io/gh/sunhailin-Leo/data-pipeline-go/branch/main/graph/badge.svg" alt="codecov"></a>
  <a href="https://goreportcard.com/report/github.com/sunhailin-Leo/data-pipeline-go"><img src="https://goreportcard.com/badge/github.com/sunhailin-Leo/data-pipeline-go" alt="Go Report Card"></a>
  <a href="https://pkg.go.dev/github.com/sunhailin-Leo/data-pipeline-go"><img src="https://pkg.go.dev/badge/github.com/sunhailin-Leo/data-pipeline-go.svg" alt="Go Reference"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/sunhailin-Leo/data-pipeline-go" alt="License"></a>
  <a href="https://github.com/sunhailin-Leo/data-pipeline-go/releases"><img src="https://img.shields.io/github/v/release/sunhailin-Leo/data-pipeline-go" alt="Release"></a>
</p>

<p align="center">
  <a href="README.md">中文</a> | <a href="docs/en/README.md">English</a>
</p>

---

## 项目介绍

基于 Golang 实现一个类似 SeaTunnel 的数据同步工具，主要是为了**简便易用**：

- **数据源多样**：兼容基本常用的数据源
- **管理和维护简单**：基于容器化部署或二进制部署，部署维护简便
- **资源利用率高/高性能**：Golang 天然资源利用率高 + Channel 实现的高性能同步数据流
- **多源输入**：支持 Fan-in 模式，多个 Source 合并到同一个 Stream
- **行级数据过滤**：支持 11 种操作符的条件过滤

### 项目架构

![framework.jpg](docs/static/framework.jpg)

## 核心架构说明

- **Source（数据源）**：Kafka、RocketMQ、RabbitMQ、Pulsar、Redis、Prometheus Metrics
- **Transform（数据转换）**：Row、JSON、JsonPath 三种模式，类型转换、字段过滤、字段扩展、行级数据过滤
- **Sink（数据输出）**：ClickHouse、Console、HTTP、Kafka、Redis、PostgreSQL、MySQL、Oracle、Elasticsearch 7/8、LocalFile、Pulsar、RocketMQ、RabbitMQ
- **Stream（流管理）**：Channel 连接各模块，ants 协程池管理并发，支持多源 Fan-in 合并
- **Committer（消息确认）**：三种 ACK 模式 - 消费后确认(0)、转换后确认(1)、写入后确认(2)

## 配置中心支持

| 配置源 | CONFIG_SRC | 必需环境变量 | 可选环境变量 |
|--------|-----------|-------------|-------------|
| Local | local | LOCAL_PATH | - |
| Apollo | apollo | APOLLO_HOST, APOLLO_CONFIG_KEY | APOLLO_APP_ID, APOLLO_NAMESPACE, APOLLO_CLUSTER_KEY |
| Redis | redis | REDIS_HOST, REDIS_CONFIG_KEY | REDIS_USERNAME, REDIS_PASSWORD, REDIS_DB |
| Nacos | nacos | NACOS_IP, NACOS_PORT, NACOS_DATA_ID, NACOS_GROUP | NACOS_NAMESPACE_ID |
| Zookeeper | zookeeper | ZOOKEEPER_HOSTS, ZOOKEEPER_CONFIG_PATH | - |
| HTTP-Get | http | HTTP_HOSTS, HTTP_CONFIG_URI | HTTP_HEARTBEAT_URI, HTTP_HEARTBEAT_INTERVAL_SECS |

## Transform 模式说明

- **Row 模式**：按分隔符拆分原始文本
- **JSON 模式**：字段映射、类型转换、is_ignore、is_strict_mode、is_keep_keys、is_expand
- **JsonPath 模式**：JsonPath 表达式提取嵌套数据

### 数据过滤

在 JSON / JsonPath 模式下支持行级数据过滤，支持以下操作符：

| 操作符 | 说明 | 示例 |
|--------|------|------|
| `eq` | 等于 | `{"field": "status", "operator": "eq", "value": "active"}` |
| `neq` | 不等于 | `{"field": "status", "operator": "neq", "value": "deleted"}` |
| `gt` / `gte` | 大于 / 大于等于 | `{"field": "age", "operator": "gt", "value": 18}` |
| `lt` / `lte` | 小于 / 小于等于 | `{"field": "score", "operator": "lt", "value": 60}` |
| `contains` | 包含子串 | `{"field": "name", "operator": "contains", "value": "test"}` |
| `not_contains` | 不包含子串 | `{"field": "name", "operator": "not_contains", "value": "temp"}` |
| `regex` | 正则匹配 | `{"field": "email", "operator": "regex", "value": "^.*@gmail\\.com$"}` |
| `in` / `not_in` | 在/不在列表中 | `{"field": "type", "operator": "in", "value": ["A","B"]}` |

## 监控与告警

- **Prometheus Metrics**：端口 8080，路径 `/metrics`
- **pprof**：通过 `net/http/pprof` 自动注册
- **Grafana Dashboard**：预置 13 个监控面板（`deploy/grafana/dashboard.json`）
- **告警规则**：预置 7 条 Prometheus 告警规则（`deploy/prometheus/alert_rules.yml`）
- **监控栈部署**：`docker compose -f deploy/docker-compose.monitoring.yml up -d`

## 快速启动

### 前置要求

- Go >= 1.24.0
- 配置文件（JSON 格式）

### 添加作业配置文件来定义作业

配置文件示例：[example/kafka_to_http.json](example/kafka_to_http.json)

```json
{
    "streams": [
        {
            "name": "stream-1",
            "enable": true,
            "channel_size": 1000,
            "source": [
                {
                    "type": "Kafka",
                    "source_name": "kafka-1",
                    "kafka": {
                        "address": "kfk-01.com:9092,kfk-01.com:9092,kfk-01.com:9092",
                        "group": "test-default",
                        "topic": "data-pipeline-events"
                    }
                }
            ],
            "transform": {
                "mode": "json",
                "schemas": [
                    {
                        "source_key": "key",
                        "sink_key":  "key",
                        "converter": "toString",
                        "is_ignore": false,
                        "is_strict_mode": true,
                        "is_keep_keys": true,
                        "source_name": "kafka-1",
                        "sink_name": "http-1"
                    }
                ]
            },
            "sink": [
                {
                    "type": "HTTP",
                    "sink_name": "http-1",
                    "http": {
                        "url": "http://0.0.0.0:8000/api/v1",
                        "content_type": "application/json",
                        "headers": {
                            "key": "value"
                        }
                    }
                }
            ]
        }
    ]
}
```

配置文件说明：
- 配置文件的格式为 JSON
- `streams` 为作业数组，每个元素为一个独立的数据流
- `source` 为输入源（支持多个，自动 Fan-in 合并）
- `transform` 为数据转换规则
- `sink` 为输出目标

### 直接运行

```shell
export CONFIG_SRC=local
export LOCAL_PATH=../example/kafka_to_http.json
cd data-pipeline-go/cmd && go run main.go
```

### 从二进制文件启动

```shell
cd data-pipeline-go
make build
export CONFIG_SRC=local
export LOCAL_PATH=example/kafka_to_http.json
./cmd/data-pipeline-go
```

### Docker 部署

```shell
make docker-build
docker run --rm \
  -e CONFIG_SRC=local \
  -e LOCAL_PATH=/app/config.json \
  -v $(pwd)/example/kafka_to_http.json:/app/config.json \
  data-pipeline-go
```

### 运行效果

#### 随机写入 10 条数据到 kafka 中
![write_kafka.png](docs/static/write_kafka.png)

#### data-pipeline-go 运行结果
![dpg_result.png](docs/static/dpg_result.png)

#### HTTP 接口打印请求数据
![http_resp.png](docs/static/http_resp.png)

## 开发指南

### Makefile 命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示帮助信息 |
| `make lint` | golangci-lint 静态检查 |
| `make nilaway` | nilaway nil 检查 |
| `make test` | 运行单元测试 |
| `make coverage` | 生成覆盖率报告（HTML） |
| `make benchmark` | 运行基准测试 |
| `make build` | 构建二进制包 |
| `make fmt` | 格式化代码 |
| `make docker-build` | 构建 Docker 镜像 |
| `make clean` | 清理构建文件 |

### 静态检查工具

- **golangci-lint**：`curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6` 或 `brew install golangci-lint`
- **nilaway**：`go install go.uber.org/nilaway/cmd/nilaway@latest`

## 实现模块

[ROADMAP](ROADMAP.md)

## 版本日志

[CHANGELOG](CHANGELOG.md)

## 贡献

欢迎提交 Issue 和 Pull Request！请确保：

1. 代码通过 `make lint` 检查
2. 新功能附带单元测试
3. 更新相关文档

## 许可证

[MIT License](LICENSE)
## 版本日志

[CHANGELOG](CHANGELOG.md)
