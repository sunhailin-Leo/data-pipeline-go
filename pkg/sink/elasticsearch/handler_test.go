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
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewElasticsearchSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	esAddr := testutil.GetEnvOrDefault(testutil.EnvESAddr, "http://localhost:9200")
	esUser := testutil.GetEnvOrDefault(testutil.EnvESUser, "elastic")
	esPass := testutil.GetEnvOrDefault(testutil.EnvESPass, "testpass")

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
			Address:   esAddr,
			Username:  esUser,
			Password:  esPass,
			IndexName: "integration-test-index",
			DocIdName: "id",
			Version:   "7.X",
		},
	}

	elasticsearchClient := NewElasticsearchSinkHandler(base, testSinkConfig.Elasticsearch)
	go elasticsearchClient.WriteData()

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

	time.Sleep(5 * time.Second)
	elasticsearchClient.CloseSink()
}
