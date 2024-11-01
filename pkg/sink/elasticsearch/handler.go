package sink

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cenkalti/backoff/v4"
	elasticsearchV7 "github.com/elastic/go-elasticsearch/v7"
	elasticsearchV8 "github.com/elastic/go-elasticsearch/v8"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/sink"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

const (
	SupportES7VersionTagName string = "7.X"
	SupportES8VersionTagName string = "8.X"
)

type ElasticsearchSinkHandler struct {
	sink.BaseSink

	sinkElasticsearchCfg config.ElasticsearchSinkConfig
	sinkES7Client        *elasticsearchV7.Client
	sinkES8Client        *elasticsearchV8.Client
}

func (e *ElasticsearchSinkHandler) SinkName() string {
	return utils.SinkElasticsearchTagName
}

// write data to elasticsearch
func (e *ElasticsearchSinkHandler) write(data *models.TransformOutput) error {
	docBytes := data.Data[0].([]byte)

	var documentId string
	// Getting the id field from nested json data is not supported yet.
	idNode, getIdErr := sonic.Get(docBytes, e.sinkElasticsearchCfg.DocIdName)
	if getIdErr != nil {
		logger.Logger.Warn(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]" +
			"Unable to retrieve the id field from the transformed data, which may cause the document write to fail...")
	} else {
		documentId, _ = idNode.String()
	}

	switch e.sinkElasticsearchCfg.Version {
	case SupportES7VersionTagName:
		res, indexErr := e.sinkES7Client.Index(
			e.sinkElasticsearchCfg.IndexName,
			bytes.NewReader(docBytes),
			e.sinkES7Client.Index.WithDocumentID(documentId),
			e.sinkES7Client.Index.WithContext(context.Background()),
			e.sinkES7Client.Index.WithRefresh("true"))
		defer func() {
			if res != nil {
				_ = res.Body.Close()
			}
		}()
		if indexErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]write data error! Reason: " + indexErr.Error())
			return indexErr
		}
		return nil
	case SupportES8VersionTagName:
		res, indexErr := e.sinkES8Client.Index(
			e.sinkElasticsearchCfg.IndexName,
			bytes.NewReader(docBytes),
			e.sinkES8Client.Index.WithDocumentID(documentId),
			e.sinkES8Client.Index.WithContext(context.Background()),
			e.sinkES8Client.Index.WithRefresh("true"))
		defer func() {
			if res != nil {
				_ = res.Body.Close()
			}
		}()
		if indexErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]write data error! Reason: " + indexErr.Error())
			return indexErr
		}
		return nil
	default:
		return errors.New("elasticsearch version is empty or cannot support")
	}
}

// WriteData writes data to Elasticsearch
func (e *ElasticsearchSinkHandler) WriteData() {
	logger.Logger.Info(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]start waiting data writing...")
	for {
		if e.sinkES7Client == nil && e.sinkES8Client == nil {
			logger.Logger.Error(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]Elasticsearch client already closed or not initialize!")
			return
		}
		data, ok := <-e.GetFromTransformChan()
		e.Metrics.OnSinkInput(e.StreamName, e.SinkAliasName)
		if !ok {
			logger.Logger.Error(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]data source is already closed!")
			return
		}
		e.Metrics.OnSinkInputSuccess(e.StreamName, e.SinkAliasName)
		e.Metrics.OnSinkOutput(e.StreamName, e.SinkAliasName)
		if writeErr := e.write(data); writeErr != nil {
			logger.Logger.Error(utils.LogServiceName + "[Kafka-Sink][Current config: " + e.SinkAliasName + "]send data error! Reason: " + writeErr.Error())
		} else {
			e.Metrics.OnSinkOutputSuccess(e.StreamName, e.SinkAliasName)
			e.MessageCommit(data.SourceObj, data.SourceData, e.SinkAliasName)
		}
	}
}

// InitSink initializes the Elasticsearch sink
func (e *ElasticsearchSinkHandler) InitSink() {
	retryBackoff := backoff.NewExponentialBackOff()
	elasticsearchTransportDialer := &net.Dialer{
		Timeout:   time.Second,
		KeepAlive: 30 * time.Second,
	}
	elasticsearchRetryStatus := []int{
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusTooManyRequests,
	}
	elasticsearchRetryFunc := func(i int) time.Duration {
		if i == 1 {
			retryBackoff.Reset()
		}
		return retryBackoff.NextBackOff()
	}
	// TODO update to fasthttp
	elasticsearchTransport := &http.Transport{
		MaxIdleConnsPerHost:   100,
		ResponseHeaderTimeout: time.Second,
		DialContext:           elasticsearchTransportDialer.DialContext,
	}

	// Version
	switch e.sinkElasticsearchCfg.Version {
	case SupportES7VersionTagName:
		es7Config := elasticsearchV7.Config{
			Addresses: strings.Split(e.sinkElasticsearchCfg.Address, ","),
			// Retry on 429, 502, 503, 504
			RetryOnStatus: elasticsearchRetryStatus,
			// Configure the backoff function
			RetryBackoff: elasticsearchRetryFunc,
			Username:     e.sinkElasticsearchCfg.Username,
			Password:     e.sinkElasticsearchCfg.Password,
			Transport:    elasticsearchTransport,
		}
		es7Client, clientErr := elasticsearchV7.NewClient(es7Config)
		if clientErr != nil {
			logger.Logger.Fatal(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]Failed to connect Elasticsearch! Reason for exception: " + clientErr.Error())
			return
		}
		e.sinkES7Client = es7Client
	case SupportES8VersionTagName:
		es8Config := elasticsearchV8.Config{
			Addresses: strings.Split(e.sinkElasticsearchCfg.Address, ","),
			// Retry on 429, 502, 503, 504
			RetryOnStatus: elasticsearchRetryStatus,
			// Configure the backoff function
			RetryBackoff: elasticsearchRetryFunc,
			Username:     e.sinkElasticsearchCfg.Username,
			Password:     e.sinkElasticsearchCfg.Password,
			Transport:    elasticsearchTransport,
		}
		es8Client, clientErr := elasticsearchV8.NewClient(es8Config)
		if clientErr != nil {
			logger.Logger.Fatal(utils.LogServiceName + "[Elasticsearch-Sink][Current config: " + e.SinkAliasName + "]Failed to connect Elasticsearch! Reason for exception: " + clientErr.Error())
			return
		}
		e.sinkES8Client = es8Client
	}
}

// CloseSink closes the ElasticsearchElasticsearch sink
func (e *ElasticsearchSinkHandler) CloseSink() {
	e.Close()
}

// NewElasticsearchSinkHandler creates a new Elasticsearch sink handler
func NewElasticsearchSinkHandler(baseSink sink.BaseSink, sinkElasticsearchCfg config.ElasticsearchSinkConfig) *ElasticsearchSinkHandler {
	handler := &ElasticsearchSinkHandler{
		BaseSink:             baseSink,
		sinkElasticsearchCfg: sinkElasticsearchCfg,
	}
	handler.InitSink()
	handler.SetFromTransformChan()
	return handler
}
