Sink
===
## Sink 模块简介
在 ETL（Extract, Transform, Load）流程中，Sink 模块用于将处理后的数据输出到指定的存储目标。Sink 是数据流的最后一环，负责将经过清洗、转换的数据持久化到数据库、消息队列或文件系统中，以支持后续的查询或分析。

## Sink 结构
当前的 Sink 结构包含以下两个字段：

Type：表示 Sink 的类型，例如数据库、消息队列、文件存储等。

SinkName：标识 Sink 的名称，用于区分不同的 Sink 实例，便于管理和监控。

## Sink 类型支持
* [ClickHouse](./sink/clickhouse)
  * 用于将数据写入 ClickHouse 数据库，适合高效的数据分析。

* [HTTP](./sink/http) 
  * 通过 HTTP 请求发送数据到指定的 REST API 端点。

* [Kafka](./sink/kafka) 
  * 发布数据到 Kafka 消息队列，支持高吞吐量的数据流处理。

* [Redis](./sink/redis) 
  * 将数据存储到 Redis，适合快速的数据访问和缓存。

* [LocalFile](./sink/local_file)
  * 将数据写入本地文件系统，便于日志记录和备份。

* [PostgresSQL](./sink/postgressql) 
  * 将数据存储到 PostgresSQL 数据库，支持复杂查询和事务处理。

* [RocketMQ](./sink/rocketmq)
  * 发布数据到 RocketMQ 队列，适用于高可用的分布式消息传递。

* [RabbitMQ](./sink/rabbitmq)
  * 将数据发布到 RabbitMQ 队列，适合消息驱动的架构。

* [Oracle](./sink/oracle)
  * 将数据写入 Oracle 数据库，适用于企业级应用。

* [MySQL](./sink/mysql)
  * 将数据存储到 MySQL 数据库，适合关系型数据存储。

* [Pulsar](./sink/pulsar)
  * 发布数据到 Apache Pulsar，支持多租户和高吞吐量。

* [Elasticsearch](./sink/elasticsearch)
  * 将数据写入 Elasticsearch，适合实时搜索和分析。