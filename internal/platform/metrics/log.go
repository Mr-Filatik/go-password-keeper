package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type LogMetrics struct {
	BaseMetrics

	logCounter *prometheus.CounterVec
}

func NewLogMetrics(base BaseMetrics) *LogMetrics {
	subsystemName := "log"

	logCounter := base.CreateCounter(CounterOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "logs_total",
			Help:       "Total number of logs.",
			LabelNames: []string{"level"},
		},
	})

	return &LogMetrics{
		BaseMetrics: base,
		logCounter:  logCounter,
	}
}

type LogLabel struct {
	Level string
}

func (p *LogMetrics) IncLogsCounter(labels LogLabel) {
	lbls := prometheus.Labels{
		"level": labels.Level,
	}

	p.logCounter.With(lbls).Inc()
}
