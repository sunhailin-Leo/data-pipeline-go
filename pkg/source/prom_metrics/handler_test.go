package source

import (
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewPromMetricSourceHandler(t *testing.T) {
	t.Helper()
	// Pre-Test
	initLogger()
	// Source - Metrics Puller
	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		SourceAliasName: "prom-metrics-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourcePromMetricsTagName,
			SourceName: "prom-metrics-1",
			PromMetrics: config.PromMetricsSourceConfig{
				Address:  "http://localhost:8080/metrics",
				Interval: 60,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}

	pm := NewPromMetricSourceHandler(baseSource)
	// pm.SetDebugMode(true)
	c := pm.GetToTransformChan()

	// Consumer - Fetch Data
	go pm.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from PromMetrics failed")
	}

	pm.CloseSource()
}
