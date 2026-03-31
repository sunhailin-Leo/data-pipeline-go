Sink
===
## Sink Module Introduction
In the ETL (Extract, Transform, Load) process, the Sink module is used to output processed data to specified storage targets. Sink is the final stage of the data flow, responsible for persisting cleaned and transformed data to databases, message queues, or file systems to support subsequent queries or analysis.

## Sink Structure
The current Sink structure contains the following two fields:

Type: Represents the type of Sink, such as database, message queue, file storage, etc.

SinkName: Identifies the name of the Sink, used to distinguish different Sink instances for easy management and monitoring.

## Supported Sink Types
* [ClickHouse](./sink/clickhouse)
  * Used to write data to ClickHouse database, suitable for efficient data analysis.

* [HTTP](./sink/http) 
  * Sends data to specified REST API endpoints via HTTP requests.

* [Kafka](./sink/kafka) 
  * Publishes data to Kafka message queue, supporting high-throughput data stream processing.

* [Redis](./sink/redis) 
  * Stores data in Redis, suitable for fast data access and caching.

* [LocalFile](./sink/local_file)
  * Writes data to local file system, convenient for logging and backup.

* [PostgresSQL](./sink/postgressql) 
  * Stores data in PostgresSQL database, supporting complex queries and transaction processing.

* [RocketMQ](./sink/rocketmq)
  * Publishes data to RocketMQ queue, suitable for high-availability distributed messaging.

* [RabbitMQ](./sink/rabbitmq)
  * Publishes data to RabbitMQ queue, suitable for message-driven architectures.

* [Oracle](./sink/oracle)
  * Writes data to Oracle database, suitable for enterprise applications.

* [MySQL](./sink/mysql)
  * Stores data in MySQL database, suitable for relational data storage.

* [Pulsar](./sink/pulsar)
  * Publishes data to Apache Pulsar, supporting multi-tenancy and high throughput.

* [Elasticsearch](./sink/elasticsearch)
  * Writes data to Elasticsearch, suitable for real-time search and analysis.

* Console
  * Outputs data to console (standard output), suitable for debugging and development phases.
