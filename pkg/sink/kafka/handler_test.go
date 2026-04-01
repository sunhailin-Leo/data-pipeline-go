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

func TestNewKafkaSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	kafkaAddr := testutil.GetEnvOrDefault(testutil.EnvKafkaAddr, "localhost:9092")

	base := sink.BaseSink{
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamName:    "",
		SinkAliasName: "kafka-1",
		ChanSize:      100,
		StreamConfig:  &config.StreamConfig{},
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkKafkaTagName,
		SinkName: "kafka-1",
		Kafka: config.KafkaSinkConfig{
			Address:     kafkaAddr,
			Topic:       "integration-test-topic",
			MessageMode: "json",
		},
	}

	kafkaClient := NewKafkaSinkHandler(base, testSinkConfig.Kafka)
	go kafkaClient.WriteData()

	k := kafkaClient.GetFromTransformChan()
	for i := 1; i < 5; i++ {
		k <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "kafka-1",
		}
	}

	time.Sleep(5 * time.Second)
	kafkaClient.CloseSink()
}
