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

func TestNewPulsarSinkHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	pulsarAddr := testutil.GetEnvOrDefault(testutil.EnvPulsarAddr, "localhost:6650")

	base := sink.BaseSink{
		DebugMode:     false,
		StreamName:    "",
		SinkAliasName: "pulsar-1",
		Metrics:       middlewares.NewMetrics("data_tunnel"),
		StreamConfig:  &config.StreamConfig{},
	}
	testSinkConfig := &config.SinkConfig{
		Type:     utils.SinkPulsarTagName,
		SinkName: "pulsar-1",
		Pulsar: config.PulsarSinkConfig{
			Address: pulsarAddr,
			Topic:   "integration-test-topic",
		},
	}

	pulsarClient := NewPulsarSinkHandler(base, testSinkConfig.Pulsar)
	go pulsarClient.WriteData()

	p := pulsarClient.GetFromTransformChan()
	for i := 0; i < 5; i++ {
		p <- &models.TransformOutput{
			SourceOutput: &models.SourceOutput{},
			Data: []any{
				[]byte(`{"id":` + cast.ToString(i) + `,"name":"test"}`),
			},
			SinkName: "pulsar-1",
		}
	}

	time.Sleep(5 * time.Second)
	pulsarClient.CloseSink()
}
