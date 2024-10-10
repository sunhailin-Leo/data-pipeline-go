package sink

import (
	"context"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kprom"
	"github.com/twmb/franz-go/plugin/kzap"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type KafkaSinkHandler struct {
	sink.BaseSink

	sinkAddress string
	sinkTopic   string

	kafkaClient *kgo.Client
}

// SinkName returns the name of the Kafka sink
func (k *KafkaSinkHandler) SinkName() string {
	return utils.SinkKafkaTagName
}

// WriteData writes data to Kafka
func (k *KafkaSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName + "[Kafka-Sink][Current config: " + k.SinkAliasName + "]start waiting data writing...")
	for {
		if k.kafkaClient == nil {
			logger.Logger.Error(utils.LogServiceName + "[Kafka-Sink][Current config: " + k.SinkAliasName + "]Kafka client already closed or not initialize!")
			return
		}
		data, ok := <-k.GetFromTransformChan()
		k.Metrics.OnSinkInput(k.StreamName, k.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName + "[Kafka-Sink][Current config: " + k.SinkAliasName + "]data source is already closed!")
			return
		}
		k.Metrics.OnSinkInputSuccess(k.StreamName, k.SinkAliasName)
		k.Metrics.OnSinkOutput(k.StreamName, k.SinkAliasName)
		if produceErr := k.kafkaClient.ProduceSync(context.Background(), kgo.SliceRecord(data.Data[0].([]byte))).FirstErr(); produceErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[Kafka-Sink][Current config: " + k.SinkAliasName + "]send data error! Reason: " + produceErr.Error())
		} else {
			k.Metrics.OnSinkOutputSuccess(k.StreamName, k.SinkAliasName)
			k.MessageCommit(data.SourceObj, data.SourceData, k.SinkAliasName)
		}
	}
}

// InitSink initializes the Kafka sink
func (k *KafkaSinkHandler) InitSink() {
	// kgo zap logger and prom config
	l := kzap.New(logger.Logger, kzap.Level(kgo.LogLevelInfo))

	// kprom metrics
	metrics := kprom.NewMetrics("kafka_sink", kprom.GoCollectors())

	// kgo options
	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(k.sinkAddress, ",")...),
		kgo.WithHooks(metrics),
		kgo.DefaultProduceTopic(k.sinkTopic),
		kgo.WithLogger(l),
	}

	kgoClient, clientErr := kgo.NewClient(opts...)
	if clientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Kafka-Sink][Current config: " + k.SinkAliasName + "]Kafka client failed to create! Reason: " + clientErr.Error())
		return
	}
	k.kafkaClient = kgoClient
	logger.Logger.Info(utils.LogServiceName + "[Kafka-Sink][Current config: " + k.SinkAliasName + "]initialize successful!")
}

// CloseSink closes the Kafka sink
func (k *KafkaSinkHandler) CloseSink() {
	if k.kafkaClient != nil {
		k.kafkaClient.Close()
	}
	k.Close()
}

// NewKafkaSinkHandler creates a new Kafka sink handler
func NewKafkaSinkHandler(baseSink sink.BaseSink, sinkKafkaCfg config.KafkaSinkConfig) *KafkaSinkHandler {
	handler := &KafkaSinkHandler{
		BaseSink:    baseSink,
		sinkAddress: sinkKafkaCfg.Address,
		sinkTopic:   sinkKafkaCfg.Topic,
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
