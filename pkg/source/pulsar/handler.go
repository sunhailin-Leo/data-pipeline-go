package source

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type PulsarSourceHandler struct {
	source.BaseSource

	sourceAddress          string
	sourceTopic            string
	sourceSubscriptionName string

	client   pulsar.Client
	consumer pulsar.Consumer
}

// SourceName returns the name of the Pulsar source
func (p *PulsarSourceHandler) SourceName() string {
	return utils.SourcePulsarTagName
}

// SourceTopic returns the topic of the Pulsar source
func (p *PulsarSourceHandler) SourceTopic() string {
	return p.sourceTopic
}

// FetchData fetches data from Pulsar
func (p *PulsarSourceHandler) FetchData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Start waiting for data to be written...")
	for {
		if p.consumer == nil || p.client == nil {
			logger.Logger.Fatal(utils.LogServiceName +
				"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Pulsar consumer/client is already closed or not connected!")
			return
		}

		receiveMsg, receiveMsgErr := p.consumer.Receive(context.Background())
		if receiveMsgErr != nil {
			logger.Logger.Error(utils.LogServiceName +
				"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Failed to read message in Pulsar, Reason for exception: " + receiveMsgErr.Error())
			return
		}

		if p.DebugMode || p.GetToTransformChan() == nil {
			logger.Logger.Info(utils.LogServiceName +
				"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Pulsar consume data: " + string(receiveMsg.Payload()))
		} else {
			logger.Logger.Debug(utils.LogServiceName +
				"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Pulsar consume data: " + string(receiveMsg.Payload()))
			// 往 Transform 管道写数据
			p.GetToTransformChan() <- &models.SourceOutput{
				MetaData:   p.MetaData,
				SourceData: receiveMsg,
			}
		}

		p.MessageCommit(p.consumer, receiveMsg, p.SourceAliasName)
		p.Metrics.OnSourceOutput(p.StreamName, p.SourceAliasName)
		p.Metrics.OnSourceOutputSuccess(p.StreamName, p.SourceAliasName)
	}
}

// InitSource initializes the Pulsar source
func (p *PulsarSourceHandler) InitSource() {
	client, clientErr := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://" + p.sourceAddress})
	if clientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Failed to create Pulsar client! Reason for exception: " + clientErr.Error())
		return
	}

	consumer, createConsumerErr := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            p.sourceTopic,
		SubscriptionName: p.sourceSubscriptionName,
		Type:             pulsar.Shared,
	})
	if createConsumerErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Pulsar-Source][Current config: " + p.SourceAliasName + "]Failed to create Pulsar consumer! Reason for exception: " + createConsumerErr.Error())
		return
	}

	p.client = client
	p.consumer = consumer
	p.MetaData = &models.MetaData{
		StreamName:    p.StreamName,
		SourceTagName: p.SourceName(),
		AliasName:     p.SourceAliasName,
		SourceObj:     p.consumer,
	}
}

// CloseSource closes the Pulsar source
func (p *PulsarSourceHandler) CloseSource() {
	if p.consumer != nil {
		p.consumer.Close()
	}
	if p.client != nil {
		p.client.Close()
	}
	p.Close()
}

// NewPulsarSourceHandler initializes a new Pulsar source handler
func NewPulsarSourceHandler(baseSource source.BaseSource) *PulsarSourceHandler {
	handler := &PulsarSourceHandler{
		BaseSource:             baseSource,
		sourceAddress:          baseSource.SourceConfig.Pulsar.Address,
		sourceTopic:            baseSource.SourceConfig.Pulsar.Topic,
		sourceSubscriptionName: baseSource.SourceConfig.Pulsar.SubscriptionName,
	}
	handler.InitSource()
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[Pulsar-Source][Current config: " + baseSource.SourceConfig.SourceName + "]Pulsar init successful!")
	return handler
}
