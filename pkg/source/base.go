package source

import (
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/commiter"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// Source interface
type Source interface {
	SourceName() string                  // source name
	SourceTopic() string                 // source topic (for MQ)
	SetSourceAliasName(aliasName string) // set source alias name
	GetSourceAliasName() string          // get source alias name
	FetchData()                          // fetch data

	InitSource()  // init source
	CloseSource() // close source

	SetToTransformChan()                           // set to transform module channel
	GetToTransformChan() chan *models.SourceOutput // get transform module channel
	SetMetricsHooks(metrics *middlewares.Metrics)  // set metrics hook
	SetStreamName(streamName string)               // set stream name
	SetDebugMode(debug bool)                       // set debug mode
}

type BaseSource struct {
	DebugMode       bool
	ChanSize        int
	AckMode         int
	StreamName      string
	SourceAliasName string
	SourceConfig    *config.SourceConfig
	Metrics         *middlewares.Metrics
	MetaData        *models.MetaData

	toTransformChan chan *models.SourceOutput
}

func (b *BaseSource) SourceName() string {
	panic("implement SourceName")
}

func (b *BaseSource) SourceTopic() string {
	panic("implement SourceTopic")
}

func (b *BaseSource) SetSourceAliasName(aliasName string) {
	b.SourceAliasName = aliasName
}

func (b *BaseSource) GetSourceAliasName() string {
	return b.SourceAliasName
}

func (b *BaseSource) FetchData() {
	panic("implement FetchData")
}

func (b *BaseSource) InitSource() {
	panic("implement InitSource")
}

func (b *BaseSource) CloseSource() {
	panic("implement CloseSource")
}

func (b *BaseSource) SetToTransformChan() {
	b.toTransformChan = make(chan *models.SourceOutput, b.ChanSize)
}

func (b *BaseSource) GetToTransformChan() chan *models.SourceOutput {
	return b.toTransformChan
}

func (b *BaseSource) SetMetricsHooks(metrics *middlewares.Metrics) {
	b.Metrics = metrics
}

func (b *BaseSource) SetStreamName(streamName string) {
	b.StreamName = streamName
}

func (b *BaseSource) SetDebugMode(debug bool) {
	b.DebugMode = debug
}

func (b *BaseSource) MessageCommit(client, message interface{}, configName string, params ...interface{}) {
	if b.AckMode == utils.AckModeInSource {
		commiter.MessageCommit(client, message, configName, params)
	}
}

func (b *BaseSource) Close() {
	if b.toTransformChan != nil {
		close(b.toTransformChan)
	}
}
