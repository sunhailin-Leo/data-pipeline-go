package sink

import (
	"context"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type RocketMQSinkHandler struct {
	sink.BaseSink

	sinkAddress string
	sinkTopic   string

	rmqClient rocketmq.Producer
}

// SinkName returns the name of the RocketMQ sink
func (r *RocketMQSinkHandler) SinkName() string {
	return utils.SinkRocketMQTagName
}

// WriteData writes data to RocketMQ
func (r *RocketMQSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Start waiting for data to be written...")
	if r.rmqClient == nil {
		logger.Logger.Error(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]RocketMQ client connection closed or not connected!!")
		return
	}
	startErr := r.rmqClient.Start()
	if startErr != nil {
		if startErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Failed to start RocketMQ consumer, Reason for exception: " + startErr.Error())
			return
		}
	}

	for {
		if r.rmqClient == nil {
			logger.Logger.Error(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]RocketMQ Client is already closed or not connected!!")
			return
		}

		data, ok := <-r.GetFromTransformChan()
		r.Metrics.OnSinkInput(r.StreamName, r.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Sink is already closed!")
			return
		}
		r.Metrics.OnSinkInputSuccess(r.StreamName, r.SinkAliasName)
		r.Metrics.OnSinkOutput(r.StreamName, r.SinkAliasName)
		msg := &primitive.Message{Topic: r.sinkTopic, Body: data.Data[0].([]byte)}
		if sendRes, sendErr := r.rmqClient.SendSync(context.Background(), msg); sendErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Failed to send data! Reason for exception: " + sendErr.Error())
		} else {
			logger.Logger.Debug(utils.LogServiceName + "[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Send data successful! Message content: " + sendRes.String())
			r.Metrics.OnSinkOutputSuccess(r.StreamName, r.SinkAliasName)
			r.MessageCommit(data.SourceObj, data.SourceData, r.SinkAliasName)
		}
	}
}

// InitSink initializes the RocketMQ sink
func (r *RocketMQSinkHandler) InitSink() {
	rmqClient, rmqClientErr := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(strings.Split(r.sinkAddress, ","))),
		producer.WithRetry(2))
	if rmqClientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RocketMQ-Sink][Current config: " + r.SinkAliasName + "]Failed to create RocketMQ Client! Reason for exception: " + rmqClientErr.Error())
		return
	}
	r.rmqClient = rmqClient
}

// CloseSink closes the RocketMQ sink
func (r *RocketMQSinkHandler) CloseSink() {
	if r.rmqClient != nil {
		_ = r.rmqClient.Shutdown()
	}
	r.Close()
}

// NewRocketMQSinkHandler creates a new RocketMQ sink handler
func NewRocketMQSinkHandler(baseSink sink.BaseSink, sinkRocketMQCfg config.RocketMQSinkConfig) *RocketMQSinkHandler {
	handler := &RocketMQSinkHandler{
		BaseSink:    baseSink,
		sinkAddress: sinkRocketMQCfg.Address,
		sinkTopic:   sinkRocketMQCfg.Topic,
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
