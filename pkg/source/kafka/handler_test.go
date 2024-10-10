package source

import (
	"context"
	"strings"
	"testing"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	devKafkaHosts string = "<test kafka hosts>"
)

func initLogger() {
	logger.NewZapLogger()
}

func testKafkaSourceMock() *kgo.Client {
	seeds := strings.Split(devKafkaHosts, ",")

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation())
	if err != nil {
		panic(err)
	}
	return cl
}

func TestNewKafkaSource(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Producer
	testProducer := testKafkaSourceMock()
	defer testProducer.Close()
	// Source - Consumer
	testTopic := "alns-script"
	baseSource := source.BaseSource{
		Metrics:         middlewares.NewMetrics("data_tunnel"),
		SourceAliasName: "kafka-1",
		StreamName:      "",
		ChanSize:        100,
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceKafkaTagName,
			SourceName: "Kafka-1",
			Kafka: config.KafkaSourceConfig{
				Address:  devKafkaHosts,
				Group:    utils.ServiceName,
				Topic:    testTopic,
				User:     "",
				Password: "",
			},
		},
	}
	kafka := NewKafkaSource(baseSource)
	// kafka.SetDebugMode(true)
	c := kafka.GetToTransformChan()

	// Producer - Send Data
	record := &kgo.Record{Topic: testTopic, Value: []byte(`{"name": "test"}`)}
	if err := testProducer.ProduceSync(context.Background(), record).FirstErr(); err != nil {
		t.Fatalf("record had a produce error while synchronously producing: %v\n", err)
	}

	// Consumer - Fetch Data
	go kafka.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from Kafka failed")
	}

	// Parse Data
	fetchRecord := fetchData.SourceData.(*kgo.Record)
	println(string(fetchRecord.Value))

	kafka.CloseSource()
}
