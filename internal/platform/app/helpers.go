package app

import (
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// WriteStartMetric writes information about the component's start to metrics.
func WriteStartMetric(
	component IComponent,
	status metrics.AppStartStatus,
	startTime time.Time,
	metricsProvider IAppMetrics,
) {
	if metricsProvider == nil {
		return
	}

	_, nok := component.(*nopComponent)
	if nok {
		return
	}

	_, pok := component.(*ParallelComponent)
	if pok {
		return
	}

	_, sok := component.(*SequentialComponent)
	if sok {
		return
	}

	metricsProvider.SetStartDuration(metrics.AppStartLabel{
		Status:    status,
		Component: component.GetName(),
	}, time.Since(startTime))
}

// WriteStopMetric writes information about the component's stop to metrics.
func WriteStopMetric(
	component IComponent,
	status metrics.AppStopStatus,
	startTime time.Time,
	metricsProvider IAppMetrics,
) {
	if metricsProvider == nil {
		return
	}

	_, nok := component.(*nopComponent)
	if nok {
		return
	}

	_, pok := component.(*ParallelComponent)
	if pok {
		return
	}

	_, sok := component.(*SequentialComponent)
	if sok {
		return
	}

	metricsProvider.SetStopDuration(metrics.AppStopLabel{
		Status:    status,
		Component: component.GetName(),
	}, time.Since(startTime))
}
