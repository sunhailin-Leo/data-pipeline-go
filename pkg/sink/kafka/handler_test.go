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

func TestNewKafkaSinkHandler(t *testing.T) {
	t.Helper()
	// init logger
	initLogger()
	// Sink config
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
			Address:     "<test kafka hosts>",
			Topic:       "<test kafka topic>",
			MessageMode: "json",
		},
	}
	// init KafkaSinkHandler
	kafkaClient := NewKafkaSinkHandler(base, testSinkConfig.Kafka)
	// Sink Write
	go kafkaClient.WriteData()
	// Channel
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
	// for waiting data insert
	time.Sleep(10 * time.Second)

	kafkaClient.CloseSink()
}
