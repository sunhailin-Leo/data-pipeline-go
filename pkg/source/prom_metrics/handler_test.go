package source

import (
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/middlewares"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/testutil"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNewPromMetricSourceHandler(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	initLogger()

	baseSource := source.BaseSource{
		DebugMode:       false,
		ChanSize:        100,
		SourceAliasName: "prom-metrics-1",
		SourceConfig: &config.SourceConfig{
			Type:       utils.SourcePromMetricsTagName,
			SourceName: "prom-metrics-1",
			PromMetrics: config.PromMetricsSourceConfig{
				Address:  "http://localhost:8080/metrics",
				Interval: 5,
			},
		},
		Metrics: middlewares.NewMetrics("data_tunnel"),
	}

	pm := NewPromMetricSourceHandler(baseSource)
	c := pm.GetToTransformChan()

	go pm.FetchData()
	fetchData, ok := <-c
	if !ok || fetchData == nil {
		t.Fatalf("Fetch data from PromMetrics failed")
	}

	pm.CloseSource()
}
