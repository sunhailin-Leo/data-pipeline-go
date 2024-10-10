* 2024-10-10 v0.1.8
  * Adjust some code structures to be more consistent with design patterns

* 2024-09-22 v0.1.7
  * Adjust config and config loader code structure
  * Update core go.mod dependencies

* 2024-09-05 v0.1.6
  * Added Redis Source mode

* 2024-09-03 v0.1.5
  * Added Pulsar's Sink and Source modes

* 2024-09-02 v0.1.4
  * Added Source mode for Prometheus and RabbitMQ
  * Added a new way to obtain data flow configuration from Apollo (hot update is not supported yet)

* 2024-08-30 v0.1.3
  * Added MySQL Sink mode
  * Adjust the loading mode of config

* 2024-08-29 v0.1.2
  * Added Oracle's Sink mode

* 2024-08-27 v0.1.1
  * Added RocketMQ, RabbitMQ Sink mode
  * Update go.mod and golang minimum compatible versions

* 2024-08-26 v0.1.0
  * Added RocketMQ source compatibility (unit test to be added)
  * Adjust the configuration level of Source Config
  * Update example configuration and default configuration

* 2024-08-23 v0.0.9
  * Update the archive location of default config.json
  * Update Makefile to be compatible with windows and linux/mac usage requirements

* 2024-08-22 v0.0.8
  * Added committer module for MQ data source submission tasks
    * No good solution has been thought of for the batch submission of the sink part yet, and the others have been implemented.
  * Add Makefile

* 2024-08-21 v0.0.7
  * Update go.mod
  * Complete the refactoring and testing of transform
  * Complete metadata penetration

* 2024-08-15 v0.0.6
  * Complete the extension of PostgresSQLSink
  * Complete the extension of LocalFileSink (except parquet)

* 2024-08-14 v0.0.5
  * Complete the extension of RedisSink
  * Adjust some codes (implemented by sinking to Base layer)
  * Set docs directory @liziwei
  * Set example directory @liziwei

* 2024-08-13 v0.0.5
  * Complete the extension of KafkaSink
  * Complete Prometheus monitoring and metrics slots
  * Complete pprof monitoring access

* 2024-08-12 v0.0.4
  * Complete the extension of ConsoleSink, HTTPSink
  * Completed optimization and adjustment of some structures
  * Complete basic unit testing
  * Complete the logic of automatic table creation for Clickhouse Sink

* 2024-08-09 v0.0.3
  * Complete basic end-to-end consumption testing

* 2024-08-06 v0.0.2
  * Completed Clickhouse Sink module
  * Complete the Stream and Transform structure modules

* 2024-08-05 v0.0.1
  * Complete the Kafka Source module
