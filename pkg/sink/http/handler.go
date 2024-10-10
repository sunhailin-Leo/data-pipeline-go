package sink

import (
	"time"

	"github.com/valyala/fasthttp"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultTimeout = 3 * time.Second
)

type HTTPSinkHandler struct {
	sink.BaseSink

	httpClient  *fasthttp.Client
	sinkHTTPCfg config.HTTPSinkConfig
}

// SinkName return HTTP sink name
func (h *HTTPSinkHandler) SinkName() string {
	return utils.SinkHTTPTagName
}

// sendData send http data
func (h *HTTPSinkHandler) sendData(data []any) {
	// request
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(h.sinkHTTPCfg.URL)
	req.Header.SetContentType(h.sinkHTTPCfg.ContentType)

	// temporary only support POST method
	req.Header.SetMethod(fasthttp.MethodPost)
	// For the time being, only fixed-parameter Headers are supported.
	if h.sinkHTTPCfg.Headers != nil {
		for hKey, hValue := range h.sinkHTTPCfg.Headers {
			req.Header.Set(hKey, hValue)
		}
	}
	req.SetBody(data[0].([]byte))

	h.Metrics.OnSinkOutput(h.StreamName, h.SinkAliasName)

	// response
	if err := h.httpClient.DoTimeout(req, nil, defaultTimeout); err != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[HTTP-Sink][Current config: " + h.SinkAliasName + "]send data error: " + err.Error())
		return
	}
	logger.Logger.Info(utils.LogServiceName + "[HTTP-Sink][Current config: " + h.SinkAliasName + "]send data successful!")
	h.Metrics.OnSinkOutputSuccess(h.StreamName, h.SinkAliasName)
}

func (h *HTTPSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName +
		"[HTTP-Sink][Current config: " + h.SinkAliasName + "]start waiting data writing...")

	for {
		if h.httpClient == nil {
			logger.Logger.Error(utils.LogServiceName +
				"[HTTP-Sink][Current config: " + h.SinkAliasName + "]HTTP client not initialised or closed!")
			return
		}

		data, ok := <-h.GetFromTransformChan()
		if !ok {
			logger.Logger.Info(utils.LogServiceName +
				"[HTTP-Sink][Current config: " + h.SinkAliasName + "]data write channel closed!")
			return
		}
		logger.Logger.Debug(utils.LogServiceName +
			"[HTTP-Sink][Current config: " + h.SinkAliasName + "]Receive data: " + string(data.Data[0].([]byte)))

		h.Metrics.OnSinkInput(h.StreamName, h.SinkAliasName)
		h.Metrics.OnSinkInputSuccess(h.StreamName, h.SinkAliasName)

		// send data
		h.sendData(data.Data)
		h.MessageCommit(data.SourceObj, data.SourceData, h.SinkAliasName)
	}
}

// InitSink initialize Sink
func (h *HTTPSinkHandler) InitSink() {
	h.httpClient = &fasthttp.Client{}
	if h.sinkHTTPCfg.ReadTimeoutSecs > 0 {
		h.httpClient.ReadTimeout = time.Duration(h.sinkHTTPCfg.ReadTimeoutSecs) * time.Second
	}
	if h.sinkHTTPCfg.WriteTimeoutSecs > 0 {
		h.httpClient.WriteTimeout = time.Duration(h.sinkHTTPCfg.WriteTimeoutSecs) * time.Second
	}
	if h.sinkHTTPCfg.MaxIdleConnDurationSecs > 0 {
		h.httpClient.MaxIdleConnDuration = time.Duration(h.sinkHTTPCfg.MaxIdleConnDurationSecs) * time.Second
	}
	if h.sinkHTTPCfg.MaxConnWaitTimeoutSecs > 0 {
		h.httpClient.MaxConnWaitTimeout = time.Duration(h.sinkHTTPCfg.MaxConnWaitTimeoutSecs) * time.Second
	}
	logger.Logger.Info(utils.LogServiceName + "[HTTP-Sink][Current config: " + h.SinkAliasName + "]初始化成功!")
}

// CloseSink close Sink
func (h *HTTPSinkHandler) CloseSink() {
	if h.httpClient != nil {
		h.httpClient.CloseIdleConnections()
		h.httpClient = nil
	}
	h.Close()
}

// NewHTTPSinkHandler initialize HTTP Sink
func NewHTTPSinkHandler(baseSink sink.BaseSink, sinkHTTPCfg config.HTTPSinkConfig) *HTTPSinkHandler {
	handler := &HTTPSinkHandler{BaseSink: baseSink, sinkHTTPCfg: sinkHTTPCfg}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
