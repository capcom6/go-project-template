package agent

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricsNamespace = "template"
	metricsSubsystem = "agent"
)

type Metrics struct {
	tasksProcessed *prometheus.CounterVec
	llmCallsTotal  prometheus.Counter
	llmDuration    prometheus.Histogram
	loopIterations prometheus.Counter
	queueDepth     prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		tasksProcessed: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "tasks_processed_total",
			Help:      "Total number of tasks processed by status",
		}, []string{"status"}),
		llmCallsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "llm_calls_total",
			Help:      "Total number of LLM calls made",
		}),
		llmDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "llm_duration_seconds",
			Help:      "Duration of LLM calls in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		loopIterations: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "loop_iterations_total",
			Help:      "Total number of agent loop iterations",
		}),
		queueDepth: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "queue_depth",
			Help:      "Current number of pending tasks in the queue",
		}),
	}
}

func (m *Metrics) IncTasksProcessed(status string) {
	m.tasksProcessed.WithLabelValues(status).Inc()
}

func (m *Metrics) IncLLMCalls() {
	m.llmCallsTotal.Inc()
}

func (m *Metrics) ObserveLLMDuration(duration float64) {
	m.llmDuration.Observe(duration)
}

func (m *Metrics) IncLoopIterations() {
	m.loopIterations.Inc()
}

func (m *Metrics) SetQueueDepth(depth int) {
	m.queueDepth.Set(float64(depth))
}
