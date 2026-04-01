package source

import (
	"context"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewPulsarSourceHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	pulsarAddr := testutil.GetEnvOrDefault(testutil.EnvPulsarAddr, "localhost:6650")
	testTopic := "integration-test-source-topic"

	// Producer
	client, clientErr := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://" + pulsarAddr})
	if clientErr != nil {
		t.Fatalf("Failed to create Pulsar client: %v", clientErr)
	}

	producer, producerErr := client.CreateProducer(pulsar.ProducerOptions{Topic: testTopic})
	if producerErr != nil {
		t.Fatalf("Failed to create Pulsar producer: %v", producerErr)
	}

	// Source - Consumer
	baseSource := source.BaseSource{
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "pulsar-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourcePulsarTagName,
			SourceName: "pulsar-1",
			Pulsar: config.PulsarSourceConfig{
				Address:          pulsarAddr,
				Topic:            testTopic,
				SubscriptionName: utils.ServiceName,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	p := NewPulsarSourceHandler(&baseSource)
	c := p.GetToTransformChan()

	// Producer - Send Data
	if _, sendErr := producer.Send(context.Background(), &pulsar.ProducerMessage{Payload: []byte(`{"name": "test"}`)}); sendErr != nil {
		t.Fatalf("record had a produce error while synchronously producing: %v\n", sendErr)
	}
	producer.Close()
	client.Close()

	// Consumer - Fetch Data
	go p.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from Pulsar failed")
	}

	fetchMessage := fetchData.SourceData.(pulsar.Message)
	t.Logf("Received message: %s", string(fetchMessage.Payload()))

	p.CloseSource()
}
