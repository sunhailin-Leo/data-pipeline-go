package sink

import (
	"fmt"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type ConsoleSinkHandler struct {
	sink.BaseSink
}

// SinkName return Console name
func (c *ConsoleSinkHandler) SinkName() string {
	return utils.SinkConsoleTagName
}

// WriteData write Console data
func (c *ConsoleSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[Console-Sink][Current config: " + c.SinkAliasName + "]start waiting data writing...")
	for {
		data, ok := <-c.GetFromTransformChan()
		if !ok {
			logger.Logger.Error(utils.LogServiceName +
				"[Console-Sink][Current config: " + c.SinkAliasName + "]data source is close already!")
			return
		}

		c.Metrics.OnSinkInput(c.StreamName, c.SinkAliasName)
		c.Metrics.OnSinkInputSuccess(c.StreamName, c.SinkAliasName)

		logger.Logger.Info(utils.LogServiceName +
			"[Console-Sink][Current config: " + c.SinkAliasName + "]data: " + fmt.Sprintf("%v", data))

		c.Metrics.OnSinkOutput(c.StreamName, c.SinkAliasName)
		c.Metrics.OnSinkOutputSuccess(c.StreamName, c.SinkAliasName)
		c.MessageCommit(data.SourceObj, data.SourceData, c.SinkAliasName)
	}
}

// InitSink initialize Console Sink
func (c *ConsoleSinkHandler) InitSink() {
	logger.Logger.Info(utils.LogServiceName + "[Console-Sink][Current config: " + c.SinkAliasName + "]initialize successful!")
}

// CloseSink close Console Sink
func (c *ConsoleSinkHandler) CloseSink() {
	c.Close()
}

// NewConsoleSinkHandler initialize Console Sink
func NewConsoleSinkHandler(baseSink sink.BaseSink) *ConsoleSinkHandler {
	handler := &ConsoleSinkHandler{BaseSink: baseSink}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
