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

func TestNewOracleSinkHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink Oracle Test
	base := sink.BaseSink{
		ChanSize:      100,
		StreamName:    "",
		SinkAliasName: "Oracle-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
	}
	testOracleConfig := config.OracleSinkConfig{
		Address:   "<test address>",
		Username:  "<test username>",
		Password:  "<test password>",
		Database:  "orcl",
		TableName: "<test table>",
		BulkSize:  5,
	}
	oraClient := NewOracleSinkHandler(base, testOracleConfig)
	// Sink Write
	go oraClient.WriteData()

	// channel
	c := oraClient.GetFromTransformChan()
	for i := 1; i < 20; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				i,
				"col" + cast.ToString(i),
				1.1 + cast.ToFloat64(i),
				time.Now(),
			},
			SinkName: "oracle-1",
		}
	}

	// for waiting data insert
	time.Sleep(20 * time.Second)
	oraClient.CloseSink()
}
