package sink

import (
	"context"
	"fmt"
	"net/url"

	"github.com/wagslane/go-rabbitmq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type RabbitMQSinkHandler struct {
	sink.BaseSink

	sinkAddress      string
	sinkUsername     string
	sinkPassword     string
	sinkRoutingKey   string
	sinkExchangeName string

	rmqConn      *rabbitmq.Conn
	rmqPublisher *rabbitmq.Publisher
}

// SinkName returns the name of the RabbitMQ sink
func (r *RabbitMQSinkHandler) SinkName() string {
	return utils.SinkRabbitMQTagName
}

// WriteData writes data to RabbitMQ
func (r *RabbitMQSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]开始发送数据...")

	for {
		if r.rmqConn == nil || r.rmqPublisher == nil {
			logger.Logger.Error(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]RabbitMQ Client 已关闭或未配置!")
			return
		}

		data, ok := <-r.GetFromTransformChan()
		r.Metrics.OnSinkInput(r.StreamName, r.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]数据源已关闭!")
			return
		}
		r.Metrics.OnSinkInputSuccess(r.StreamName, r.SinkAliasName)
		r.Metrics.OnSinkOutput(r.StreamName, r.SinkAliasName)

		confirms, publishErr := r.rmqPublisher.PublishWithDeferredConfirmWithContext(
			context.Background(),
			data.Data[0].([]byte),
			[]string{r.sinkRoutingKey},
			rabbitmq.WithPublishOptionsContentType(utils.HTTPContentTypeJSON),
			rabbitmq.WithPublishOptionsExchange(r.sinkExchangeName),
			rabbitmq.WithPublishOptionsPersistentDelivery,
			rabbitmq.WithPublishOptionsMandatory)
		if publishErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]发送数据失败! 错误原因: " + publishErr.Error())
			continue
		}
		if len(confirms) == 0 || confirms[0] == nil {
			logger.Logger.Error(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]发送数据还未确认!")
			continue
		}

		isConfirm, confirmErr := confirms[0].WaitContext(context.Background())
		if confirmErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]发送数据确认失败! 错误原因: " + confirmErr.Error())
		}

		if isConfirm {
			r.Metrics.OnSinkOutputSuccess(r.StreamName, r.SinkAliasName)
			r.MessageCommit(data.SourceObj, data.SourceData, r.SinkAliasName)
		}
	}
}

// InitSink initializes the RabbitMQ sink
func (r *RabbitMQSinkHandler) InitSink() {
	rmqConn, rmqConnErr := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s/", url.QueryEscape(r.sinkUsername), url.QueryEscape(r.sinkPassword), r.sinkAddress),
		rabbitmq.WithConnectionOptionsLogging)
	if rmqConnErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]RabbitMQ Connection 创建失败! 错误原因: " + rmqConnErr.Error())
		return
	}

	publisher, publisherErr := rabbitmq.NewPublisher(
		rmqConn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(r.sinkExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsConfirm)
	if publisherErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RabbitMQ-Sink][Current config: " + r.SinkAliasName + "]RabbitMQ Publisher 创建失败! 错误原因: " + publisherErr.Error())
		return
	}
	r.rmqConn = rmqConn
	r.rmqPublisher = publisher
}

// CloseSink closes the RabbitMQ sink
func (r *RabbitMQSinkHandler) CloseSink() {
	if r.rmqPublisher != nil {
		r.rmqPublisher.Close()
	}
	if r.rmqConn != nil {
		_ = r.rmqConn.Close()
	}
	r.Close()
}

// NewRabbitMQSinkHandler creates a new RabbitMQ sink handler
func NewRabbitMQSinkHandler(baseSink sink.BaseSink, sinkRabbitMQCfg config.RabbitMQSinkConfig) *RabbitMQSinkHandler {
	handler := &RabbitMQSinkHandler{
		BaseSink:         baseSink,
		sinkAddress:      sinkRabbitMQCfg.Address,
		sinkUsername:     sinkRabbitMQCfg.Username,
		sinkPassword:     sinkRabbitMQCfg.Password,
		sinkRoutingKey:   sinkRabbitMQCfg.RoutingKey,
		sinkExchangeName: sinkRabbitMQCfg.Exchange,
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
