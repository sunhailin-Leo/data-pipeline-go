package sink

import (
	"testing"
	"time"

	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewMySQLSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	mysqlAddr := testutil.GetEnvOrDefault(testutil.EnvMySQLAddr, "localhost:3306")
	mysqlUser := testutil.GetEnvOrDefault(testutil.EnvMySQLUser, "root")
	mysqlPass := testutil.GetEnvOrDefault(testutil.EnvMySQLPass, "testpass")
	mysqlDB := testutil.GetEnvOrDefault(testutil.EnvMySQLDB, "integration_test")

	base := sink.BaseSink{
		ChanSize:      100,
		StreamName:    "",
		SinkAliasName: "MySQL-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
	}
	testMySQLConfig := config.MySQLSinkConfig{
		Address:   mysqlAddr,
		Username:  mysqlUser,
		Password:  mysqlPass,
		Database:  mysqlDB,
		TableName: "integration_test_table",
		BulkSize:  5,
	}

	mysqlClient := NewMySQLSinkHandler(base, testMySQLConfig)
	go mysqlClient.WriteData()

	c := mysqlClient.GetFromTransformChan()
	for i := 1; i < 20; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				i,
				"col" + cast.ToString(i),
				1.1 + cast.ToFloat64(i),
				time.Now(),
			},
			SinkName: "mysql-1",
		}
	}

	time.Sleep(5 * time.Second)
	mysqlClient.CloseSink()
}
