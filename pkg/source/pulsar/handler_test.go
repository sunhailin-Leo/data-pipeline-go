package source

import (
	"context"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func testPulsarSourceMock(topic string) (pulsar.Client, pulsar.Producer) {
	client, clientErr := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://172.20.49.19:16650"})
	if clientErr != nil {
		panic(clientErr)
	}

	producer, producerErr := client.CreateProducer(pulsar.ProducerOptions{Topic: topic})
	if producerErr != nil {
		panic(producerErr)
	}

	return client, producer
}

func TestNewPulsarSourceHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Client, Producer
	client, producer := testPulsarSourceMock("alg_test")
	// Source - Consumer
	baseSource := source.BaseSource{
		ChanSize:        100,
		StreamName:      "",
		SourceAliasName: "pulsar-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourcePulsarTagName,
			SourceName: "pulsar-1",
			Pulsar: config.PulsarSourceConfig{
				Address:          "172.20.49.19:16650",
				Topic:            "alg_test",
				SubscriptionName: utils.ServiceName,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}
	p := NewPulsarSourceHandler(baseSource)
	// p.SetDebugMode(true)
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
	println(string(fetchMessage.Payload()))

	p.CloseSource()
}
