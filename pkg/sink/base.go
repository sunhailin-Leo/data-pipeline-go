package sink

import (
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/commiter"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// Sink interface
type Sink interface {
	SinkName() string                  // sink name
	SetSinkAliasName(aliasName string) // set Sink alias name
	GetSinkAliasName() string          // get Sink alias name
	WriteData()                        // write data

	InitSink()  // init sink
	CloseSink() // close sink

	SetFromTransformChan()                              // set from transform module channel
	GetFromTransformChan() chan *models.TransformOutput // get from transform module channel
	SetMetricsHooks(metrics *middlewares.Metrics)       // set metrics hook
	SetStreamName(streamName string)                    // set stream name
	SetDebugMode(debug bool)                            // set debug mode
}

type BaseSink struct {
	DebugMode     bool
	ChanSize      int
	StreamName    string
	SinkAliasName string
	StreamConfig  *config.StreamConfig
	Metrics       *middlewares.Metrics

	fromTransformChan chan *models.TransformOutput
}

func (b *BaseSink) SinkName() string {
	panic("implement SinkName")
}

func (b *BaseSink) SetSinkAliasName(aliasName string) {
	b.SinkAliasName = aliasName
}

func (b *BaseSink) GetSinkAliasName() string {
	return b.SinkAliasName
}

func (b *BaseSink) WriteData() {
	panic("implement WriteData")
}

func (b *BaseSink) InitSink() {
	panic("implement InitSink")
}

func (b *BaseSink) CloseSink() {
	panic("implement CloseSink")
}

func (b *BaseSink) SetFromTransformChan() {
	b.fromTransformChan = make(chan *models.TransformOutput, b.ChanSize)
}

func (b *BaseSink) GetFromTransformChan() chan *models.TransformOutput {
	return b.fromTransformChan
}

func (b *BaseSink) SetMetricsHooks(metrics *middlewares.Metrics) {
	b.Metrics = metrics
}

func (b *BaseSink) SetStreamName(streamName string) {
	b.StreamName = streamName
}

func (b *BaseSink) SetDebugMode(debug bool) {
	b.DebugMode = debug
}

func (b *BaseSink) MessageCommit(client, message interface{}, configName string, params ...interface{}) {
	if b.StreamConfig.SourceAckMode == utils.AckModeInSink {
		commiter.MessageCommit(client, message, configName, params)
	}
}

func (b *BaseSink) Close() {
	if b.fromTransformChan != nil {
		close(b.fromTransformChan)
	}
}
