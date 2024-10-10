package source

import (
	"bytes"
	"crypto/tls"
	"time"

	"github.com/bytedance/sonic"
	"github.com/prometheus/common/expfmt"
	"github.com/valyala/fasthttp"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	defaultReadTimeout int64 = 5
)

type PromMetricSourceHandler struct {
	source.BaseSource

	sourceAddress string

	httpClient    *fasthttp.Client
	ticker        *time.Ticker
	metricsParser expfmt.TextParser
}

// SourceName returns the name of the PromMetrics source
func (p *PromMetricSourceHandler) SourceName() string {
	return utils.SourcePromMetricsTagName
}

// SourceTopic returns the topic of the PromMetrics
func (p *PromMetricSourceHandler) SourceTopic() string {
	return "prom-metrics"
}

func (p *PromMetricSourceHandler) getMetrics() ([]byte, error) {
	// Request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	// set URL
	req.SetRequestURI(p.sourceAddress)
	// Response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	// Do
	reqErr := fasthttp.DoTimeout(req, resp, time.Duration(defaultReadTimeout)*time.Second)
	if reqErr != nil {
		return nil, reqErr
	}
	// Get Body
	metricsBody := resp.Body()
	return metricsBody, nil
}

func (p *PromMetricSourceHandler) parseMetrics(data []byte) ([]byte, error) {
	parsedMetrics, err := p.metricsParser.TextToMetricFamilies(bytes.NewReader(data))
	if err != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]Failed to parse metrics data, Reason for exception: " + err.Error())
		return nil, err
	}

	// TODO temporary serialize to json
	jsonMetrics, marshalErr := sonic.MarshalIndent(parsedMetrics, "", " ")
	if marshalErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]Analysis of indicator data anomalies, Reason for exception: " + marshalErr.Error())
		return nil, err
	}

	return jsonMetrics, nil
}

// FetchData fetches data from PromMetrics
func (p *PromMetricSourceHandler) FetchData() {
	logger.Logger.Info(utils.LogServiceName +
		"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]Start waiting for data to be written...")

	for {
		select {
		case <-p.ticker.C:
			p.Metrics.OnSourceInput(p.StreamName, p.SourceAliasName)
			p.Metrics.OnSourceInputSuccess(p.StreamName, p.SourceAliasName)

			metricsBody, getMetricsErr := p.getMetrics()
			if getMetricsErr != nil {
				logger.Logger.Error(utils.LogServiceName +
					"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]Get data error, Reason for exception: " + getMetricsErr.Error())
				return
			}
			parseResult, parseErr := p.parseMetrics(metricsBody)
			if parseErr != nil {
				logger.Logger.Error(utils.LogServiceName +
					"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]Parse data error, Reason for exception: " + parseErr.Error())
				return
			}
			if p.DebugMode || p.GetToTransformChan() == nil {
				logger.Logger.Info(utils.LogServiceName +
					"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]PromMetrics consume data: " + string(parseResult))
			} else {
				logger.Logger.Debug(utils.LogServiceName +
					"[PromMetrics-Source][Current config: " + p.SourceAliasName + "]PromMetrics consume data: " + string(parseResult))
				p.GetToTransformChan() <- &models.SourceOutput{
					MetaData:   p.MetaData,
					SourceData: parseResult,
				}
				p.Metrics.OnSourceOutput(p.StreamName, p.SourceAliasName)
				p.Metrics.OnSourceOutputSuccess(p.StreamName, p.SourceAliasName)
			}
		}
	}
}

// InitSource initializes the PromMetrics source
func (p *PromMetricSourceHandler) InitSource() {
	httpClient := &fasthttp.Client{ReadTimeout: time.Duration(defaultReadTimeout) * time.Second}
	httpClient.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	p.httpClient = httpClient
	p.metricsParser = expfmt.TextParser{}
}

// CloseSource closes the PromMetrics source
func (p *PromMetricSourceHandler) CloseSource() {
	if p.httpClient != nil {
		p.httpClient.CloseIdleConnections()
		p.httpClient = nil
	}
	if p.ticker != nil {
		p.ticker.Stop()
	}
	p.Close()
}

// NewPromMetricSourceHandler initializes a new PromMetrics source handler
func NewPromMetricSourceHandler(baseSource source.BaseSource) *PromMetricSourceHandler {
	handler := &PromMetricSourceHandler{
		BaseSource:    baseSource,
		sourceAddress: baseSource.SourceConfig.PromMetrics.Address,
		ticker:        time.NewTicker(time.Duration(baseSource.SourceConfig.PromMetrics.Interval) * time.Second),
	}
	handler.InitSource()
	handler.SetToTransformChan()
	logger.Logger.Info(utils.LogServiceName +
		"[PromMetrics-Source][Current config: " + baseSource.SourceConfig.SourceName + "]PromMetrics init successful!")
	return handler
}
