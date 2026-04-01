package testutil

import (
	"os"
	"testing"
)

const (
	EnvIntegrationTest = "INTEGRATION_TEST"

	EnvRedisAddr      = "INTEGRATION_REDIS_ADDR"
	EnvKafkaAddr      = "INTEGRATION_KAFKA_ADDR"
	EnvMySQLAddr      = "INTEGRATION_MYSQL_ADDR"
	EnvMySQLUser      = "INTEGRATION_MYSQL_USER"
	EnvMySQLPass      = "INTEGRATION_MYSQL_PASS"
	EnvMySQLDB        = "INTEGRATION_MYSQL_DB"
	EnvPostgresAddr   = "INTEGRATION_POSTGRES_ADDR"
	EnvPostgresUser   = "INTEGRATION_POSTGRES_USER"
	EnvPostgresPass   = "INTEGRATION_POSTGRES_PASS"
	EnvPostgresDB     = "INTEGRATION_POSTGRES_DB"
	EnvClickhouseAddr = "INTEGRATION_CLICKHOUSE_ADDR"
	EnvClickhouseUser = "INTEGRATION_CLICKHOUSE_USER"
	EnvClickhousePass = "INTEGRATION_CLICKHOUSE_PASS"
	EnvClickhouseDB   = "INTEGRATION_CLICKHOUSE_DB"
	EnvRabbitMQAddr   = "INTEGRATION_RABBITMQ_ADDR"
	EnvRabbitMQUser   = "INTEGRATION_RABBITMQ_USER"
	EnvRabbitMQPass   = "INTEGRATION_RABBITMQ_PASS"
	EnvRocketMQAddr   = "INTEGRATION_ROCKETMQ_ADDR"
	EnvPulsarAddr     = "INTEGRATION_PULSAR_ADDR"
	EnvESAddr         = "INTEGRATION_ES_ADDR"
	EnvESUser         = "INTEGRATION_ES_USER"
	EnvESPass         = "INTEGRATION_ES_PASS"
)

// SkipIfNotIntegration skips the test if INTEGRATION_TEST env var is not set.
func SkipIfNotIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv(EnvIntegrationTest) == "" {
		t.Skip("Skipping integration test: set INTEGRATION_TEST=true to run")
	}
}

// GetEnvOrDefault returns the value of the environment variable or a default value.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
