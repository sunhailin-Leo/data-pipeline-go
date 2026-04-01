package source

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

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

func TestNewKafkaSource(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	kafkaAddr := testutil.GetEnvOrDefault(testutil.EnvKafkaAddr, "localhost:9092")
	testTopic := "integration-test-source-topic-direct"

	// Producer - ensure topic exists
	seeds := strings.Split(kafkaAddr, ",")
	testProducer, err := kgo.NewClient(kgo.SeedBrokers(seeds...), kgo.AllowAutoTopicCreation())
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer testProducer.Close()

	// Produce a message to ensure topic exists
	initRecord := &kgo.Record{Topic: testTopic, Value: []byte(`{"name": "init"}`)}
	if produceErr := testProducer.ProduceSync(context.Background(), initRecord).FirstErr(); produceErr != nil {
		t.Fatalf("record had a produce error while synchronously producing: %v\n", produceErr)
	}

	// Source - Consumer without consumer group (direct partition consumption)
	baseSource := source.BaseSource{
		Metrics:         middlewares.NewMetrics("data_tunnel"),
		SourceAliasName: "kafka-1",
		StreamName:      "",
		ChanSize:        100,
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourceKafkaTagName,
			SourceName: "Kafka-1",
			Kafka: config.KafkaSourceConfig{
				Address:  kafkaAddr,
				Group:    "",
				Topic:    testTopic,
				User:     "",
				Password: "",
			},
		},
	}
	kafka := NewKafkaSource(&baseSource)
	c := kafka.GetToTransformChan()

	// Start consumer
	go kafka.FetchData()

	// Wait for consumer to be ready
	time.Sleep(2 * time.Second)

	// Send test message
	record := &kgo.Record{Topic: testTopic, Value: []byte(`{"name": "test"}`)}
	if produceErr := testProducer.ProduceSync(context.Background(), record).FirstErr(); produceErr != nil {
		t.Fatalf("record had a produce error while synchronously producing: %v\n", produceErr)
	}

	// Consumer - Fetch Data with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case fetchData, ok := <-c:
		if !ok || fetchData == nil {
			t.Fatalf("Fetch data from Kafka failed")
		}
		fetchRecord := fetchData.SourceData.(*kgo.Record)
		t.Logf("Received message: %s", string(fetchRecord.Value))
	case <-ctx.Done():
		t.Fatalf("Test timed out waiting for Kafka message")
	}

	kafka.CloseSource()
}
