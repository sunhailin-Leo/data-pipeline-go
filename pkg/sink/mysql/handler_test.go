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
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewMySQLSinkHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink MySQL Test
	base := sink.BaseSink{
		ChanSize:      100,
		StreamName:    "",
		SinkAliasName: "MySQL-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
	}
	testMySQLConfig := config.MySQLSinkConfig{
		Address:   "<test address>",
		Username:  "<test username>",
		Password:  "<test password>",
		Database:  "ai_group",
		TableName: "<test table>",
		BulkSize:  5,
	}
	mysqlClient := NewMySQLSinkHandler(base, testMySQLConfig)
	// Sink Write
	go mysqlClient.WriteData()

	// channel
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

	// for waiting data insert
	time.Sleep(20 * time.Second)
	mysqlClient.CloseSink()
}
