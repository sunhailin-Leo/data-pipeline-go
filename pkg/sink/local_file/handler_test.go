package sink

import (
	"testing"
	"time"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewLocalFileHandler(t *testing.T) {
	t.Helper()
	// 初始化日志
	initLogger()
	// Sink 配置
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "local-file-1",
		ChanSize:      100,
	}
	testSinkConfig := &config.SinkConfig{
		Type:     "LocalFile",
		SinkName: "local-file-1",
		LocalFile: config.LocalFileSinkConfig{
			FileName:       "test",
			FileFormatType: utils.LocalFileTextFormatType,
			RowDelimiter:   "\n",
		},
	}
	// 初始化 LocalFileSinkHandler
	localFileClient := NewLocalFileHandler(base, testSinkConfig.LocalFile)
	// Sink Write
	go localFileClient.WriteData()
	// Channel
	k := localFileClient.GetFromTransformChan()
	for i := 1; i < 5; i++ {
		k <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				"col1",
				2,
				3.1,
			},
			SinkName: "local-file-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	localFileClient.CloseSink()
}
