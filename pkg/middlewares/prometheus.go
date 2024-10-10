package middlewares

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	reg *prometheus.Registry

	sourceInputCounter         *prometheus.CounterVec
	sourceInputSuccessCounter  *prometheus.CounterVec
	sourceOutputCounter        *prometheus.CounterVec
	sourceOutputSuccessCounter *prometheus.CounterVec

	transformInputCounter          *prometheus.CounterVec
	transformInputSuccessCounter   *prometheus.CounterVec
	transformConvertCounter        *prometheus.CounterVec
	transformConvertSuccessCounter *prometheus.CounterVec
	transformOutputCounter         *prometheus.CounterVec
	transformOutputSuccessCounter  *prometheus.CounterVec

	sinkInputCounter         *prometheus.CounterVec
	sinkInputSuccessCounter  *prometheus.CounterVec
	sinkOutputCounter        *prometheus.CounterVec
	sinkOutputSuccessCounter *prometheus.CounterVec
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})
}

func (m *Metrics) OnSourceInput(streamID, sourceID string) {
	m.sourceInputCounter.WithLabelValues(streamID, sourceID).Inc()
}

func (m *Metrics) OnSourceInputSuccess(streamID, sourceID string) {
	m.sourceInputSuccessCounter.WithLabelValues(streamID, sourceID).Inc()
}

func (m *Metrics) OnSourceOutput(streamID, sourceID string) {
	m.sourceOutputCounter.WithLabelValues(streamID, sourceID).Inc()
}

func (m *Metrics) OnSourceOutputSuccess(streamID, sourceID string) {
	m.sourceOutputSuccessCounter.WithLabelValues(streamID, sourceID).Inc()
}

func (m *Metrics) OnTransformInput(streamID string) {
	m.transformInputCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnTransformInputSuccess(streamID string) {
	m.transformInputSuccessCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnTransformConvert(streamID string) {
	m.transformConvertCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnTransformConvertSuccess(streamID string) {
	m.transformConvertSuccessCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnTransformOutput(streamID string) {
	m.transformOutputCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnTransformOutputSuccess(streamID string) {
	m.transformOutputSuccessCounter.WithLabelValues(streamID).Inc()
}

func (m *Metrics) OnSinkInput(streamID, sinkID string) {
	m.sinkInputCounter.WithLabelValues(streamID, sinkID).Inc()
}

func (m *Metrics) OnSinkInputSuccess(streamID, sinkID string) {
	m.sinkInputSuccessCounter.WithLabelValues(streamID, sinkID).Inc()
}

func (m *Metrics) OnSinkOutput(streamID, sinkID string) {
	m.sinkOutputCounter.WithLabelValues(streamID, sinkID).Inc()
}

func (m *Metrics) OnSinkOutputSuccess(streamID, sinkID string) {
	m.sinkOutputSuccessCounter.WithLabelValues(streamID, sinkID).Inc()
}

// NewMetrics creates a new Metrics instance
func NewMetrics(namespace string) (m *Metrics) {
	reg := prometheus.NewRegistry()
	factory := promauto.With(reg)

	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector(collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll)))

	return &Metrics{
		reg: reg,

		sourceInputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "source_input_total",
			Help:      "Total number of input records from source",
		}, []string{"stream_id", "source_id"}),

		sourceInputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "source_input_success_total",
			Help:      "Total number of successful input records from source",
		}, []string{"stream_id", "source_id"}),

		sourceOutputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "source_output_total",
			Help:      "Total number of output records from source",
		}, []string{"stream_id", "source_id"}),

		sourceOutputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "source_output_success_total",
			Help:      "Total number of successful output records from source",
		}, []string{"stream_id", "source_id"}),

		transformInputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_input_total",
			Help:      "Total number of input records from transform",
		}, []string{"stream_id"}),

		transformInputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_input_success_total",
			Help:      "Total number of successful input records from transform",
		}, []string{"stream_id"}),

		transformConvertCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_convert_total",
			Help:      "Total number of converted records from transform",
		}, []string{"stream_id"}),

		transformConvertSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_convert_success_total",
			Help:      "Total number of successful converted records from transform",
		}, []string{"stream_id"}),

		transformOutputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_output_total",
			Help:      "Total number of output records from transform",
		}, []string{"stream_id"}),

		transformOutputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "transform_output_success_total",
			Help:      "Total number of successful output records from transform",
		}, []string{"stream_id"}),

		sinkInputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "sink_input_total",
			Help:      "Total number of input records from sink",
		}, []string{"stream_id", "sink_id"}),

		sinkInputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "sink_input_success_total",
			Help:      "Total number of successful input records from sink",
		}, []string{"stream_id", "sink_id"}),

		sinkOutputCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "sink_output_total",
			Help:      "Total number of output records from sink",
		}, []string{"stream_id", "sink_id"}),

		sinkOutputSuccessCounter: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "sink_output_success_total",
			Help:      "Total number of successful output records from sink",
		}, []string{"stream_id", "sink_id"}),
	}
}
