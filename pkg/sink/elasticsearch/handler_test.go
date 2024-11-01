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
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewElasticsearchSinkHandler(t *testing.T) {
	t.Helper()
	// init logger
	initLogger()
	// Sink config
	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "elasticsearch-1",
		ChanSize:      100,
		StreamConfig:  &config.StreamConfig{},
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkElasticsearchTagName,
		SinkName: "elasticsearch-1",
		Elasticsearch: config.ElasticsearchSinkConfig{
			Address:   "<test address>",
			Username:  "admin",
			Password:  "<test password>",
			IndexName: "<test document>",
			DocIdName: "id",
			Version:   "7.X",
		},
	}
	// init ElasticsearchSinkHandler
	elasticsearchClient := NewElasticsearchSinkHandler(base, testSinkConfig.Elasticsearch)
	// Sink Write
	go elasticsearchClient.WriteData()
	// Channel
	e := elasticsearchClient.GetFromTransformChan()
	for i := 1; i < 5; i++ {
		e <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "elasticsearch-1",
		}
	}
	// for waiting data insert
	time.Sleep(10 * time.Second)

	elasticsearchClient.CloseSink()
}
