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

func TestNewClickhouseSink(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	clickhouseAddr := testutil.GetEnvOrDefault(testutil.EnvClickhouseAddr, "localhost:9000")
	clickhouseUser := testutil.GetEnvOrDefault(testutil.EnvClickhouseUser, "default")
	clickhousePass := testutil.GetEnvOrDefault(testutil.EnvClickhousePass, "testpass")
	clickhouseDB := testutil.GetEnvOrDefault(testutil.EnvClickhouseDB, "integration_test")

	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "Clickhouse-1",
		ChanSize:      100,
	}
	testClickhouseConfig := config.ClickhouseSinkConfig{
		Address:           clickhouseAddr,
		Username:          clickhouseUser,
		Password:          clickhousePass,
		Database:          clickhouseDB,
		TableName:         "integration_test_table",
		IsAutoCreateTable: true,
		Columns: []config.ClickhouseTableColumn{
			{Name: "col1", Type: "Int32"},
			{Name: "col2", Type: "Float32"},
			{Name: "col3", Type: "String"},
		},
		Engine:   "MergeTree",
		OrderBy:  []string{"col1"},
		BulkSize: 5,
	}

	ckClient := NewClickhouseSink(base, &testClickhouseConfig)
	go ckClient.WriteData()

	c := ckClient.GetFromTransformChan()
	for i := 1; i < 20; i++ {
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

	time.Sleep(5 * time.Second)
	ckClient.CloseSink()
}
