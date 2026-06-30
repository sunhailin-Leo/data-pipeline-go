package source

import (
	"strings"

	"github.com/nsqio/go-nsq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type NSQSourceHandler struct {
	source.BaseSource

	sourceAddress        string
	sourceLookupdAddress string
	sourceTopic          string
	sourceChannel        string
	sourceMaxInFlight    int

	nsqConsumer *nsq.Consumer
}

// SourceName returns the name of the NSQ source.
func (n *NSQSourceHandler) SourceName() string {
	return utils.SourceNSQTagName
}

// SourceTopic returns the topic of the NSQ source.
func (n *NSQSourceHandler) SourceTopic() string {
	return n.sourceTopic
}

func (n *NSQSourceHandler) HandleMessage(message *nsq.Message) error {
	// 关闭 go-nsq 自动确认，交给 ack_mode 控制确认时机。
	message.DisableAutoResponse()

	if n.Metrics != nil {
		n.Metrics.OnSourceInput(n.StreamName, n.SourceAliasName)
		n.Metrics.OnSourceInputSuccess(n.StreamName, n.SourceAliasName)
	}

	if n.DebugMode || n.GetToTransformChan() == nil {
		logger.Logger.Info(utils.LogServiceName +
			"[NSQ-Source][Current config: " + n.SourceAliasName + "]NSQ consume data: " + string(message.Body))
		// Debug 模式没有下游消费，直接确认避免消息堆积。
		message.Finish()
		return nil
	}

	logger.Logger.Debug(utils.LogServiceName +
		"[NSQ-Source][Current config: " + n.SourceAliasName + "]NSQ consume data: " + string(message.Body))
	n.GetToTransformChan() <- &models.SourceOutput{
		MetaData: &models.MetaData{
			StreamName:    n.StreamName,
			SourceTagName: n.SourceName(),
			AliasName:     n.SourceAliasName,
			// NSQ 每条消息自己 Finish，后续 ACK 阶段需要拿到当前 message。
			SourceObj: message,
		},
		SourceData: message,
	}

	n.MessageCommit(message, message, n.SourceAliasName)

	if n.Metrics != nil {
		n.Metrics.OnSourceOutput(n.StreamName, n.SourceAliasName)
		n.Metrics.OnSourceOutputSuccess(n.StreamName, n.SourceAliasName)
	}
	return nil
}

// FetchData keeps the source task alive while the NSQ consumer is running.
func (n *NSQSourceHandler) FetchData() {
	if n.nsqConsumer == nil {
		logger.Logger.Error(utils.LogServiceName +
			"[NSQ-Source][Current config: " + n.SourceAliasName + "]NSQ consumer is closed or not configured!")
		return
	}

	logger.Logger.Info(utils.LogServiceName +
		"[NSQ-Source][Current config: " + n.SourceAliasName + "]Start waiting for data to be written...")
	// NSQ 的消息消费由 go-nsq 回调 HandleMessage 完成；这里阻塞等待，保持 Source 任务存活。
	<-n.nsqConsumer.StopChan
}

// InitSource initializes the NSQ source.
func (n *NSQSourceHandler) InitSource() {
	nsqConfig := nsq.NewConfig()
	if n.sourceMaxInFlight > 0 {
		nsqConfig.MaxInFlight = n.sourceMaxInFlight
	}

	consumer, createErr := nsq.NewConsumer(n.sourceTopic, n.sourceChannel, nsqConfig)
	if createErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[NSQ-Source][Current config: " + n.SourceAliasName + "]Failed to create NSQ consumer! Reason for exception: " + createErr.Error())
		return
	}
	consumer.AddHandler(n)

	var connectErr error
	// 优先使用 nsqlookupd；未配置时直连 nsqd。
	if n.sourceLookupdAddress != "" {
		connectErr = consumer.ConnectToNSQLookupds(strings.Split(n.sourceLookupdAddress, ","))
	} else {
		connectErr = consumer.ConnectToNSQDs(strings.Split(n.sourceAddress, ","))
	}
	if connectErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[NSQ-Source][Current config: " + n.SourceAliasName + "]Failed to connect NSQ! Reason for exception: " + connectErr.Error())
		return
	}

	n.nsqConsumer = consumer
	n.MetaData = &models.MetaData{
		StreamName:    n.StreamName,
		SourceTagName: n.SourceName(),
		AliasName:     n.SourceAliasName,
		SourceObj:     n.nsqConsumer,
	}
}

// CloseSource closes the NSQ source.
func (n *NSQSourceHandler) CloseSource() {
	if n.nsqConsumer != nil {
		n.nsqConsumer.Stop()
	}
	n.Close()
}

// NewNSQSource initializes a new NSQ source handler.
func NewNSQSource(baseSource *source.BaseSource) *NSQSourceHandler {
	handler := &NSQSourceHandler{
		BaseSource:           *baseSource,
		sourceAddress:        baseSource.SourceConfig.NSQ.Address,
		sourceLookupdAddress: baseSource.SourceConfig.NSQ.LookupdAddress,
		sourceTopic:          baseSource.SourceConfig.NSQ.Topic,
		sourceChannel:        baseSource.SourceConfig.NSQ.Channel,
		sourceMaxInFlight:    baseSource.SourceConfig.NSQ.MaxInFlight,
	}
	handler.InitSource()
	if handler.nsqConsumer == nil {
		return nil
	}
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[NSQ-Source][Current config: " + baseSource.SourceConfig.SourceName + "]NSQ init successful!")
	return handler
}
