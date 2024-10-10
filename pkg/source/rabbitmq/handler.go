package source

import (
	"fmt"
	"net/url"

	"github.com/wagslane/go-rabbitmq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type RabbitMQSourceHandler struct {
	source.BaseSource

	sourceAddress      string
	sourceUsername     string
	sourcePassword     string
	sourceQueueName    string
	sourceRoutingKey   string
	sourceExchangeName string

	rmqConn     *rabbitmq.Conn
	rmqConsumer *rabbitmq.Consumer
}

// SourceName returns the name of the RabbitMQ source
func (r *RabbitMQSourceHandler) SourceName() string {
	return utils.SourceRabbitMQTagName
}

// SourceTopic returns the topic of the RabbitMQ
func (r *RabbitMQSourceHandler) SourceTopic() string {
	return r.sourceQueueName
}

// FetchData fetches data from RabbitMQ
func (r *RabbitMQSourceHandler) FetchData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Rabbit-Source][Current config: " + r.SourceAliasName + "]Start waiting for data to be written...")

	runErr := r.rmqConsumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		r.Metrics.OnSourceInput(r.StreamName, r.SourceAliasName)
		r.Metrics.OnSourceInputSuccess(r.StreamName, r.SourceAliasName)
		if r.DebugMode || r.GetToTransformChan() == nil {
			logger.Logger.Info(utils.LogServiceName +
				"[Rabbit-Source][Current config: " + r.SourceAliasName + "]RocketMQ consume data: " + string(d.Body))
		} else {
			logger.Logger.Debug(utils.LogServiceName +
				"[Rabbit-Source][Current config: " + r.SourceAliasName + "]RocketMQ consume data: " + string(d.Body))
			r.GetToTransformChan() <- &models.SourceOutput{
				MetaData:   r.MetaData,
				SourceData: d,
			}
		}
		r.Metrics.OnSourceOutput(r.StreamName, r.SourceAliasName)
		r.Metrics.OnSourceOutputSuccess(r.StreamName, r.SourceAliasName)
		r.MessageCommit(d, nil, r.SourceAliasName)

		return rabbitmq.Manual
	})

	if runErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Rabbit-Source][Current config: " + r.SourceAliasName + "]Failed to consume, Reason for exception: " + runErr.Error())
	}
}

// InitSource initializes the RabbitMQ source
func (r *RabbitMQSourceHandler) InitSource() {
	rmqConn, rmqConnErr := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s/",
			url.QueryEscape(r.sourceUsername),
			url.QueryEscape(r.sourcePassword),
			r.sourceAddress),
		rabbitmq.WithConnectionOptionsLogging)
	if rmqConnErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RabbitMQ-Sink][Current config: " + r.SourceAliasName + "]Failed to create RabbitMQ connection! Reason for exception: " + rmqConnErr.Error())
		return
	}

	rmqConsumer, consumerErr := rabbitmq.NewConsumer(
		rmqConn,
		r.sourceQueueName,
		rabbitmq.WithConsumerOptionsRoutingKey(r.sourceRoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(r.sourceExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare)
	if consumerErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RabbitMQ-Sink][Current config: " + r.SourceAliasName + "]Failed to create RabbitMQ consumer! Reason for exception: " + consumerErr.Error())
		return
	}

	r.rmqConn = rmqConn
	r.rmqConsumer = rmqConsumer
}

func (r *RabbitMQSourceHandler) CloseSource() {
	if r.rmqConsumer != nil {
		r.rmqConsumer.Close()
	}
	if r.rmqConn != nil {
		_ = r.rmqConn.Close()
	}
	r.Close()
}

// NewRabbitMQSource initializes a new RabbitMQ source handler
func NewRabbitMQSource(baseSource source.BaseSource) *RabbitMQSourceHandler {
	handler := &RabbitMQSourceHandler{
		BaseSource:         baseSource,
		sourceAddress:      baseSource.SourceConfig.RabbitMQ.Address,
		sourceUsername:     baseSource.SourceConfig.RabbitMQ.Username,
		sourcePassword:     baseSource.SourceConfig.RabbitMQ.Password,
		sourceQueueName:    baseSource.SourceConfig.RabbitMQ.Queue,
		sourceRoutingKey:   baseSource.SourceConfig.RabbitMQ.RoutingKey,
		sourceExchangeName: baseSource.SourceConfig.RabbitMQ.Exchange,
	}
	handler.InitSource()
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[RabbitMQ-Source][Current config: " + baseSource.SourceConfig.SourceName + "]RabbitMQ init successful!")
	return handler
}
