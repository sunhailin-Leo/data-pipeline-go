package source

import (
	"context"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/plugin/kprom"
	"github.com/twmb/franz-go/plugin/kzap"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type KafkaSourceHandler struct {
	source.BaseSource

	sourceAddress string
	sourceTopic   string
	sourceGroup   string

	kafkaClient *kgo.Client
}

// SourceName returns the name of the Kafka source
func (k *KafkaSourceHandler) SourceName() string { return utils.SourceKafkaTagName }

// SourceTopic returns the topic of the Kafka source
func (k *KafkaSourceHandler) SourceTopic() string { return k.sourceTopic }

// FetchData fetches data from Kafka
func (k *KafkaSourceHandler) FetchData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Kafka-Source][Current config: " + k.SourceAliasName + "]开始消费数据...")
	for {
		if k.kafkaClient == nil {
			logger.Logger.Fatal(utils.LogServiceName +
				"[Kafka-Source][Current config: " + k.SourceAliasName + "]Kafka client 已关闭或未配置!")
			return
		}
		// PollFetches 会阻塞直到有数据到达
		fetches := k.kafkaClient.PollFetches(context.Background())
		if fetches == nil || fetches.IsClientClosed() {
			return
		}
		fetches.EachError(func(s string, i int32, err error) {
			logger.Logger.Fatal(utils.LogServiceName +
				"[Kafka-Source][Current config: " + k.SourceAliasName + "]Kafka 消费异常! 错误原因: " +
				"fetch err topic " + s + " partition " + string(i) + " err " + err.Error())
		})
		fetches.EachRecord(func(record *kgo.Record) {
			k.Metrics.OnSourceInput(k.StreamName, k.SourceAliasName)
			k.Metrics.OnSourceInputSuccess(k.StreamName, k.SourceAliasName)
			if k.DebugMode || k.GetToTransformChan() == nil {
				logger.Logger.Info(utils.LogServiceName +
					"[Kafka-Source][Current config: " + k.SourceAliasName + "]Kafka 消费数据: " + string(record.Value))
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[Kafka-Source][Current config: " + k.SourceAliasName + "]Kafka 消费数据: " + string(record.Value))
				// 往 Transform 管道写数据
				k.GetToTransformChan() <- &models.SourceOutput{
					MetaData:   k.MetaData,
					SourceData: record,
				}
			}

			k.MessageCommit(k.kafkaClient, record, k.SourceAliasName)
			k.Metrics.OnSourceOutput(k.StreamName, k.SourceAliasName)
			k.Metrics.OnSourceOutputSuccess(k.StreamName, k.SourceAliasName)
		})
	}
}

// InitSource initializes the Kafka source
func (k *KafkaSourceHandler) InitSource() {
	// kgo zap 日志和 prom 配置
	l := kzap.New(logger.Logger, kzap.Level(kgo.LogLevelInfo))
	// kprom
	metrics := kprom.NewMetrics("kafka_source", kprom.GoCollectors())

	// kgo 配置
	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(k.sourceAddress, ",")...),
		kgo.WithHooks(metrics),
		kgo.ConsumeTopics(k.sourceTopic),
		kgo.ConsumerGroup(k.sourceGroup),
		kgo.DisableAutoCommit(),
		kgo.WithLogger(l),
	}
	// kgo sasl
	if k.SourceConfig.Kafka.User != "" && k.SourceConfig.Kafka.Password != "" {
		opts = append(opts, kgo.SASL(plain.Auth{User: k.SourceConfig.Kafka.User, Pass: k.SourceConfig.Kafka.Password}.AsMechanism()))
	}

	kgoClient, clientErr := kgo.NewClient(opts...)
	if clientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Kafka-Source][Current config: " + k.SourceAliasName + "]Failed to create Kafka client! Reason for exception: " + clientErr.Error())
		return
	}
	k.kafkaClient = kgoClient
	k.MetaData = &models.MetaData{
		StreamName:    k.StreamName,
		SourceTagName: k.SourceName(),
		AliasName:     k.SourceAliasName,
		SourceObj:     k.kafkaClient,
	}
}

// CloseSource closes the Kafka source
func (k *KafkaSourceHandler) CloseSource() {
	if k.kafkaClient != nil {
		k.kafkaClient.Close()
	}
	k.Close()
}

// NewKafkaSource initializes a new Kafka source handler
func NewKafkaSource(baseSource source.BaseSource) *KafkaSourceHandler {
	sourceGroup := baseSource.SourceConfig.Kafka.Group
	if sourceGroup == "" {
		sourceGroup = utils.ServiceName
	}

	handler := &KafkaSourceHandler{
		BaseSource:    baseSource,
		sourceAddress: baseSource.SourceConfig.Kafka.Address,
		sourceTopic:   baseSource.SourceConfig.Kafka.Topic,
		sourceGroup:   sourceGroup,
	}
	handler.InitSource()
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[Kafka-Source][Current config: " + baseSource.SourceConfig.SourceName + "]Kafka 初始化成功!")
	return handler
}
