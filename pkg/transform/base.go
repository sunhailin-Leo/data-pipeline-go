package transform

import (
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
)

type Transform interface {
	From()    // read data from source
	Convert() // convert/transform data
	To()      // write data to sink

	InitTransform(transformConfig config.TransformConfig, chanSize int) Transform                       // init transform module
	ConvertModeSelector(beforeConvertData *models.TransformBeforeConvert) *models.TransformAfterConvert // Convert
	CloseTransform()                                                                                    // close transform module

	SetStreamConfig(streamConfig *config.StreamConfig) Transform // set stream config
	SetMetricsHooks(metrics *middlewares.Metrics) Transform      // set metrics hooks
}

type BaseTransform struct {
	configs           config.TransformConfig
	streamConfig      *config.StreamConfig
	metrics           *middlewares.Metrics
	beforeConvertChan chan *models.TransformBeforeConvert
	afterConvertChan  chan *models.TransformAfterConvert
}

func (b *BaseTransform) ConvertModeSelector(_ *models.TransformBeforeConvert) *models.TransformAfterConvert {
	panic("implement me")
}

func (b *BaseTransform) SetStreamConfig(streamConfig *config.StreamConfig) Transform {
	b.streamConfig = streamConfig
	return b
}

func (b *BaseTransform) From() {
	panic("implement From")
}

func (b *BaseTransform) Convert() {
	panic("implement Convert")
}

func (b *BaseTransform) To() {
	panic("implement To")
}

func (b *BaseTransform) InitTransform(_ config.TransformConfig, _ int) Transform {
	panic("implement InitTransform")
}

func (b *BaseTransform) CloseTransform() {
	panic("implement CloseTransform")
}

func (b *BaseTransform) SetMetricsHooks(metrics *middlewares.Metrics) Transform {
	b.metrics = metrics
	return b
}

func (b *BaseTransform) initChannel(chanSize int) {
	b.beforeConvertChan = make(chan *models.TransformBeforeConvert, chanSize)
	b.afterConvertChan = make(chan *models.TransformAfterConvert, chanSize)
}

func (b *BaseTransform) close() {
	// 按顺序关闭 channel
	close(b.beforeConvertChan)
	close(b.afterConvertChan)
}
