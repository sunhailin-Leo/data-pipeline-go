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

func TestNewPostgresSQLHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink PostgresSQL Test
	base := sink.BaseSink{
		ChanSize:      100,
		StreamName:    "",
		SinkAliasName: "PostgresSQL-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
	}
	testPostgresSQLConfig := config.PostgresSQLSinkConfig{
		Address:   "<test address>",
		Username:  "<test username>",
		Password:  "<test password>",
		Database:  "postgres",
		TableName: "<test table>",
		BulkSize:  5,
	}

	pgClient := NewPostgresSQLHandler(base, testPostgresSQLConfig)
	// Sink Write
	go pgClient.WriteData()

	// channel
	c := pgClient.GetFromTransformChan()
	for i := 1; i < 100; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				cast.ToString(i),
				i,
			},
			SinkName: "postgressql-1",
		}
	}

	// for waiting data insert
	time.Sleep(20 * time.Second)

	pgClient.CloseSink()
}
