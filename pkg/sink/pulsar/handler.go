package sink

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type PulsarSinkHandler struct {
	sink.BaseSink

	sinkAddress string
	sinkTopic   string

	client   pulsar.Client
	producer pulsar.Producer
}

// SinkName returns the name of the Pulsar sink
func (p *PulsarSinkHandler) SinkName() string {
	return utils.SinkPulsarTagName
}

// WriteData writes data to Pulsar
func (p *PulsarSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName + "[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Start waiting for data to be written...")
	for {
		if p.producer == nil {
			logger.Logger.Error(utils.LogServiceName + "[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Pulsar producer is already closed or not connected!")
			return
		}
		data, ok := <-p.GetFromTransformChan()
		p.Metrics.OnSinkInput(p.StreamName, p.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName + "[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Sink is already closed!")
			return
		}
		p.Metrics.OnSinkInputSuccess(p.StreamName, p.SinkAliasName)
		p.Metrics.OnSinkOutput(p.StreamName, p.SinkAliasName)
		if msgId, sendErr := p.producer.Send(context.Background(), &pulsar.ProducerMessage{Payload: data.Data[0].([]byte)}); sendErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Send message error! Reason for exception: " + sendErr.Error())
		} else {
			logger.Logger.Debug(utils.LogServiceName + "[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Message ID: " + msgId.String())
			p.Metrics.OnSinkOutputSuccess(p.StreamName, p.SinkAliasName)
			p.MessageCommit(data.SourceObj, data.SourceData, p.SinkAliasName)
		}
	}
}

// InitSink initializes the Pulsar sink
func (p *PulsarSinkHandler) InitSink() {
	client, clientErr := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://" + p.sinkAddress,
	})
	if clientErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Failed to create Pulsar client! Reason for exception: " + clientErr.Error())
		return
	}

	producer, createProducerErr := client.CreateProducer(pulsar.ProducerOptions{Topic: p.sinkTopic})
	if createProducerErr != nil {
		logger.Logger.Fatal(utils.LogServiceName +
			"[Pulsar-Sink][Current config: " + p.SinkAliasName + "]Failed to create Pulsar producer! Reason for exception: " + createProducerErr.Error())
		return
	}

	p.client = client
	p.producer = producer
}

// CloseSink closes the Pulsar sink
func (p *PulsarSinkHandler) CloseSink() {
	if p.producer != nil {
		p.producer.Close()
	}
	if p.client != nil {
		p.client.Close()
	}
	p.Close()
}

// NewPulsarSinkHandler creates a new Pulsar sink handler
func NewPulsarSinkHandler(baseSink sink.BaseSink, sinkPulsarCfg config.PulsarSinkConfig) *PulsarSinkHandler {
	handler := &PulsarSinkHandler{
		BaseSink:    baseSink,
		sinkAddress: sinkPulsarCfg.Address,
		sinkTopic:   sinkPulsarCfg.Topic,
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
