package stream

import (
	"strings"

	"github.com/panjf2000/ants/v2"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	clickhouse "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/clickhouse"
	console "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/console"
	http "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/http"
	kafkaSink "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/kafka"
	mysql "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/mysql"
	oracle "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/oracle"
	postgressql "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/postgressql"
	rabbitmqSink "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/rabbitmq"
	redis "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/redis"
	rocketmqSink "github.com/sunhailin-Leo/data-pipeline-go/pkg/sink/rocketmq"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	kafkaSource "github.com/sunhailin-Leo/data-pipeline-go/pkg/source/kafka"
	promMetrics "github.com/sunhailin-Leo/data-pipeline-go/pkg/source/prom_metrics"
	rabbitmqSource "github.com/sunhailin-Leo/data-pipeline-go/pkg/source/rabbitmq"
	rocketmqSource "github.com/sunhailin-Leo/data-pipeline-go/pkg/source/rocketmq"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/transform"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type Stream interface {
	GetSource(streamConfig *config.StreamConfig) map[string]chan *models.SourceOutput

	GetTransform(
		inputChan chan *models.SourceOutput,
		outputChanMap map[string]chan *models.TransformOutput,
		streamConfig *config.StreamConfig) transform.Transform

	GetSink(streamConfig *config.StreamConfig) map[string]chan *models.TransformOutput

	InitStream()
	CloseStream()
	Start()
}

type BaseStream struct {
	metrics        *middlewares.Metrics
	streamWorkPool *ants.Pool
	streamsConfig  []*config.StreamConfig
}

// SubmitTask submit stream sub task
func (b *BaseStream) submitTask(task func(), callerName string) bool {
	if err := b.streamWorkPool.Submit(task); err != nil {
		logger.Logger.Error(utils.LogServiceName + "[Stream-Submit]" + callerName + "任务提交失败! 失败原因: " + err.Error())
		return false
	}
	return true
}

// setupWorkPool init goroutine worker pool
func (b *BaseStream) setupWorkPool() {
	// count init goroutine then initialize it.
	poolCnt := 1
	for _, streamConfig := range b.streamsConfig {
		if streamConfig.Enable {
			poolCnt += len(streamConfig.Sink)
			poolCnt += len(streamConfig.Source)
			poolCnt += 3 // transform[From -> Convert -> To]
		}
	}
	streamWorkPool, createWorkPoolErr := ants.NewPool(poolCnt, ants.WithPreAlloc(true))
	if createWorkPoolErr != nil {
		panic(createWorkPoolErr)
	}
	b.streamWorkPool = streamWorkPool
}

// releaseWorkPool release goroutine worker pool
func (b *BaseStream) releaseWorkPool() {
	b.streamWorkPool.Release()
}

// closeSource close all source channel
func (b *BaseStream) closeSource(sourcesMap map[string]source.Source) {
	for srcName, src := range sourcesMap {
		sourceTagArr := strings.Split(srcName, ":")
		if len(sourceTagArr) == 0 {
			logger.Logger.Error(utils.LogServiceName + "[Stream-Close]Stream close error, sourceName: " + srcName)
			continue
		}
		src.CloseSource()
		logger.Logger.Info(utils.LogServiceName + "[Stream-Close]Source-" + sourceTagArr[1] + "already closed!")
	}
}

// closeSink close all sink channel
func (b *BaseStream) closeSink(sinksMap map[string]sink.Sink) {
	for sinkName, sinkObj := range sinksMap {
		sinkTagArr := strings.Split(sinkName, ":")
		if len(sinkTagArr) == 0 {
			logger.Logger.Error(utils.LogServiceName + "[Stream-Close]Sink close error, sinkName: " + sinkName)
			continue
		}
		sinkObj.CloseSink()
		logger.Logger.Info(utils.LogServiceName + "[Stream-Close]Sink-" + sinkTagArr[1] + "already closed!")
	}
}

func (b *BaseStream) InitStream() {
	panic("implement InitStream")
}

func (b *BaseStream) CloseStream() {
	panic("implement CloseStream")
}

func (b *BaseStream) Start() {
	panic("implement Start")
}

// NewSink init sink
func NewSink(sinkType string, baseSink sink.BaseSink, sinkConfig *config.SinkConfig) sink.Sink {
	switch sinkType {
	case utils.SinkClickhouseTagName:
		return clickhouse.NewClickhouseSink(baseSink, sinkConfig.Clickhouse)
	case utils.SinkConsoleTagName:
		return console.NewConsoleSinkHandler(baseSink)
	case utils.SinkHTTPTagName:
		return http.NewHTTPSinkHandler(baseSink, sinkConfig.HTTP)
	case utils.SinkKafkaTagName:
		return kafkaSink.NewKafkaSinkHandler(baseSink, sinkConfig.Kafka)
	case utils.SinkRedisTagName:
		return redis.NewRedisSinkHandler(baseSink, sinkConfig.Redis)
	case utils.SinkPostgresSQLTagName:
		return postgressql.NewPostgresSQLHandler(baseSink, sinkConfig.PostgresSQL)
	case utils.SinkRocketMQTagName:
		return rocketmqSink.NewRocketMQSinkHandler(baseSink, sinkConfig.RocketMQ)
	case utils.SinkRabbitMQTagName:
		return rabbitmqSink.NewRabbitMQSinkHandler(baseSink, sinkConfig.RabbitMQ)
	case utils.SinkOracleTagName:
		return oracle.NewOracleSinkHandler(baseSink, sinkConfig.Oracle)
	case utils.SinkMySQLTagName:
		return mysql.NewMySQLSinkHandler(baseSink, sinkConfig.MySQL)
	}
	return nil
}

// NewSource create source
func NewSource(sourceType string, baseSource source.BaseSource) source.Source {
	switch sourceType {
	case utils.SourceKafkaTagName:
		return kafkaSource.NewKafkaSource(baseSource)
	case utils.SourceRocketMQTagName:
		return rocketmqSource.NewRocketMQSource(baseSource)
	case utils.SourceRabbitMQTagName:
		return rabbitmqSource.NewRabbitMQSource(baseSource)
	case utils.SourcePromMetricsTagName:
		return promMetrics.NewPromMetricSourceHandler(baseSource)
	}

	return nil
}
