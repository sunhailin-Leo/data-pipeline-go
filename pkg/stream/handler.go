package stream

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/transform"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

type Handler struct {
	BaseStream

	streamMap  map[string]transform.Transform
	sourcesMap map[string]source.Source
	sinksMap   map[string]sink.Sink
}

// GetSource init source
func (s *Handler) GetSource(streamConfig *config.StreamConfig) map[string]chan *models.SourceOutput {

	sourceInputChanMap := make(map[string]chan *models.SourceOutput)
	for _, sourceConfig := range streamConfig.Source {
		baseSource := source.BaseSource{
			ChanSize:        streamConfig.ChannelSize,
			AckMode:         streamConfig.SourceAckMode,
			StreamName:      streamConfig.Name,
			SourceAliasName: sourceConfig.SourceName,
			SourceConfig:    sourceConfig,
			Metrics:         s.metrics,
		}

		src := NewSource(sourceConfig.Type, baseSource)
		if src == nil {
			logger.Logger.Error(utils.LogServiceName + "[Source-Init]Source configuration error! Type: " + sourceConfig.Type)
			return sourceInputChanMap
		}
		if s.submitTask(src.FetchData, sourceConfig.Type+"-Source-"+sourceConfig.SourceName) {
			sourceInputChanMap[sourceConfig.SourceName] = src.GetToTransformChan()
			s.sourcesMap[sourceConfig.Type+":"+src.GetSourceAliasName()] = src
		}
	}

	return sourceInputChanMap
}

// GetTransform init transform
func (s *Handler) GetTransform(inputChan chan *models.SourceOutput, outputChanMap map[string]chan *models.TransformOutput, streamConfig *config.StreamConfig) transform.Transform {
	t := transform.NewTransformHandler(inputChan, outputChanMap)
	t.SetMetricsHooks(s.metrics).SetStreamConfig(streamConfig).InitTransform(streamConfig.Transform, streamConfig.ChannelSize)
	return t
}

// GetSink init sink
func (s *Handler) GetSink(streamConfig *config.StreamConfig) map[string]chan *models.TransformOutput {
	sinkOutputChanMap := make(map[string]chan *models.TransformOutput)

	for _, sinkConfig := range streamConfig.Sink {
		baseSink := sink.BaseSink{
			Metrics:       s.metrics,
			StreamName:    streamConfig.Name,
			SinkAliasName: sinkConfig.SinkName,
			ChanSize:      streamConfig.ChannelSize,
			StreamConfig:  streamConfig,
		}

		sinker := NewSink(sinkConfig.Type, baseSink, sinkConfig)
		if sinker == nil {
			logger.Logger.Error(utils.LogServiceName + "[Sink-Init]Sink configuration error! Type: " + sinkConfig.Type)
			return sinkOutputChanMap
		}
		if s.submitTask(sinker.WriteData, sinkConfig.Type+"-Sink-"+sinkConfig.SinkName) {
			sinkOutputChanMap[sinker.GetSinkAliasName()] = sinker.GetFromTransformChan()
			s.sinksMap[sinker.GetSinkAliasName()] = sinker
		}
	}

	return sinkOutputChanMap
}

// InitStream init whole stream
func (s *Handler) InitStream() {
	s.startMetrics()
	s.setupWorkPool()

	// generate stream
	for _, streamConfig := range s.streamsConfig {
		if streamConfig.Enable {
			sourceChanMap := s.GetSource(streamConfig)
			// TODO - 未来兼容多输入源的时候需要改造
			sourceChan := sourceChanMap[streamConfig.Source[0].SourceName]
			sinkChanMap := s.GetSink(streamConfig)
			transformObj := s.GetTransform(sourceChan, sinkChanMap, streamConfig)
			s.streamMap[streamConfig.Name] = transformObj
			logger.Logger.Info(utils.LogServiceName + "[Stream-Init]Stream: " + streamConfig.Name + " initialize successful!")
		}
	}

	logger.Logger.Info(utils.LogServiceName + "[Stream-Init]Stream initialize successful!")
}

// CloseStream close stream
func (s *Handler) CloseStream() {
	// close all source channel
	s.closeSource(s.sourcesMap)

	// close all transform channel
	for _, stream := range s.streamMap {
		stream.CloseTransform()
	}

	// close all sink channel
	s.closeSink(s.sinksMap)

	s.releaseWorkPool()
}

// Start stream handler start
func (s *Handler) Start() {
	for streamName, stream := range s.streamMap {
		s.submitTask(stream.From, "Stream-From-"+streamName)
		s.submitTask(stream.Convert, "Stream-Convert-"+streamName)
		s.submitTask(stream.To, "Stream-To-"+streamName)
	}

	logger.Logger.Info(utils.LogServiceName + "service start successful!")
	s.spin()
}

// start prometheus and pprof http service
func (s *Handler) startMetrics() {
	http.Handle(utils.PromHTTPRoute, s.metrics.Handler())
	go func() {
		if err := http.ListenAndServe(":"+utils.PromHTTPServerPort, nil); err != nil {
			logger.Logger.Fatal(utils.LogServiceName + "failed to start metrics service! Error reason: " + err.Error())
		}
	}()
}

// runs the service until catching os.Signal.
func (s *Handler) spin() {
	// until catching os.Signal
	quitChan := make(chan os.Signal, 2)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitChan
	logger.Logger.Info(utils.LogServiceName + "receive quit signal, shutdown service")
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.CloseStream()
	}()

	select {
	case <-quitChan:
		logger.Logger.Info(utils.LogServiceName + "receive quit signal twice, force shutdown")
	case <-done:
	}
}

// NewStreamHandler create stream handler
func NewStreamHandler() *Handler {
	// init logger
	logger.NewZapLogger()
	// init config loader
	config.NewTunnelConfigLoader()
	// init stream handler
	handler := &Handler{
		BaseStream: BaseStream{
			metrics:       middlewares.NewMetrics("data"),
			streamsConfig: config.TunnelCfg.Streams,
		},
		streamMap:  make(map[string]transform.Transform),
		sourcesMap: make(map[string]source.Source),
		sinksMap:   make(map[string]sink.Sink),
	}
	handler.InitStream()
	return handler
}
