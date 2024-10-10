package source

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/cast"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type RocketMQSourceHandler struct {
	source.BaseSource

	sourceAddress string
	sourceTopic   string
	sourceGroup   string

	rmqClient rocketmq.PullConsumer
}

// SourceName returns the name of the RocketMQ source
func (r *RocketMQSourceHandler) SourceName() string {
	return utils.SourceRocketMQTagName
}

// SourceTopic returns the topic of the RocketMQ
func (r *RocketMQSourceHandler) SourceTopic() string {
	return r.sourceTopic
}

// FetchData fetches data from RocketMQ
func (r *RocketMQSourceHandler) FetchData() {
	startErr := r.rmqClient.Start()
	if startErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Failed to create RocketMQ consume, Reason for exception: " + startErr.Error())
		return
	}
	logger.Logger.Info(utils.LogServiceName +
		"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Start waiting for data to be written...")
	for {
		resp, respErr := r.rmqClient.Pull(context.Background(), 1)
		if respErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]pull error: " + respErr.Error())
			return
		}

		switch resp.Status {
		case primitive.PullFound:
			logger.Logger.Info(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]RocketMQ : " +
				fmt.Sprintf("pull message successfully MinOffset:%d, MaxOffset:%d, nextOffset: %d, len:%d",
					resp.MinOffset, resp.MaxOffset, resp.NextBeginOffset, len(resp.GetMessages())))

			if len(resp.GetMessages()) == 0 {
				return
			}

			var queue *primitive.MessageQueue
			for _, msg := range resp.GetMessageExts() {
				// parse messageQueue
				queue = msg.Queue
				r.Metrics.OnSourceInput(r.StreamName, r.SourceAliasName)
				r.Metrics.OnSourceInputSuccess(r.StreamName, r.SourceAliasName)
				if r.DebugMode || r.GetToTransformChan() == nil {
					logger.Logger.Info(utils.LogServiceName +
						"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]RocketMQ consume data: " + string(msg.Body))
				} else {
					logger.Logger.Debug(utils.LogServiceName +
						"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]RocketMQ consume data: " + string(msg.Body))
					r.GetToTransformChan() <- &models.SourceOutput{
						MetaData:   r.MetaData,
						SourceData: msg,
					}
				}
				r.Metrics.OnSourceOutput(r.StreamName, r.SourceAliasName)
				r.Metrics.OnSourceOutputSuccess(r.StreamName, r.SourceAliasName)
			}

			r.MessageCommit(r.rmqClient, queue, r.SourceAliasName, resp.NextBeginOffset)

		case primitive.PullNoNewMsg, primitive.PullNoMsgMatched:
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]no pull message, next=" + cast.ToString(resp.NextBeginOffset))
			return
		case primitive.PullBrokerTimeout:
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]pull broker timeout, next=" + cast.ToString(resp.NextBeginOffset))
			return
		case primitive.PullOffsetIllegal:
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]pull offset illegal, next=" + cast.ToString(resp.NextBeginOffset))
			return
		default:
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]pull error, next=" + cast.ToString(resp.NextBeginOffset))
		}
	}
}

// InitSource initializes the RocketMQ source
func (r *RocketMQSourceHandler) InitSource() {
	rmqClient, rmqClientErr := rocketmq.NewPullConsumer(
		consumer.WithAutoCommit(false),
		consumer.WithGroupName(r.SourceConfig.RocketMQ.Group),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(strings.Split(r.sourceAddress, ","))),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: r.SourceConfig.RocketMQ.AccessKey,
			SecretKey: r.SourceConfig.RocketMQ.SecretKey,
		}),
		consumer.WithMaxReconsumeTimes(2))
	if rmqClientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Failed to create RocketMQ client! Reason for exception: " + rmqClientErr.Error())
		return
	}
	subscribeErr := rmqClient.Subscribe(r.sourceTopic, consumer.MessageSelector{})
	if subscribeErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Failed to subscribe RocketMQ topic, Reason for exception: " + subscribeErr.Error())
		return
	}
	r.rmqClient = rmqClient
	r.MetaData = &models.MetaData{
		StreamName:    r.StreamName,
		SourceTagName: r.SourceName(),
		AliasName:     r.SourceAliasName,
		SourceObj:     rmqClient,
	}
}

// CloseSource closes the RocketMQ source
func (r *RocketMQSourceHandler) CloseSource() {
	if r.rmqClient != nil {
		unsubscribeErr := r.rmqClient.Unsubscribe(r.sourceTopic)
		if unsubscribeErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Failed unsubscribe RocketMQ! Reason for exception: " + unsubscribeErr.Error())
		}
		shutdownErr := r.rmqClient.Shutdown()
		if shutdownErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[RocketMQ-Source][Current config: " + r.SourceAliasName + "]Failed shutdown RocketMQ client! Reason for exception: " + shutdownErr.Error())
		}
	}
	r.Close()
}

// NewRocketMQSource initializes a new RocketMQ source handler
func NewRocketMQSource(baseSource source.BaseSource) *RocketMQSourceHandler {
	sourceGroup := baseSource.SourceConfig.RocketMQ.Group
	if sourceGroup == "" {
		sourceGroup = utils.ServiceName
	}

	handler := &RocketMQSourceHandler{
		BaseSource:    baseSource,
		sourceAddress: baseSource.SourceConfig.RocketMQ.Address,
		sourceTopic:   baseSource.SourceConfig.RocketMQ.Topic,
		sourceGroup:   sourceGroup,
	}
	handler.InitSource()
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[RocketMQ-Source][Current config: " + baseSource.SourceConfig.SourceName + "]RocketMQ init successful!")
	return handler
}
