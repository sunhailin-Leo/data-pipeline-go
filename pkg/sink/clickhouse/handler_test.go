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

func TestNewClickhouseSink(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Sink Clickhouse Test
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "Clickhouse-1",
		ChanSize:      100,
	}
	testClickhouseConfig := config.ClickhouseSinkConfig{
		Address:   "<test address>",
		Username:  "<test username>",
		Password:  "<test password>",
		Database:  "ai_group_test",
		TableName: "test",
		Columns: []config.ClickhouseTableColumn{
			{Name: "col1", Type: "Int32"},
			{Name: "col2", Type: "Float32"},
			{Name: "col3", Type: "String"},
		},
		BulkSize: 5,
	}
	ckClient := NewClickhouseSink(base, testClickhouseConfig)
	// Sink Write
	go ckClient.WriteData()

	// channel
	c := ckClient.GetFromTransformChan()
	for i := 1; i < 100; i++ {
		c <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				i,
				cast.ToFloat32(i),
				cast.ToString(i),
			},
			SinkName: "clickhouse-1",
		}
	}

	// for waiting data insert
	time.Sleep(20 * time.Second)

	ckClient.CloseSink()
}
