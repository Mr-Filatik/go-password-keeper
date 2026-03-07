package app

import (
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// WriteStartMetric writes information about the component's start to metrics.
func WriteStartMetric(comp IComponent, status metrics.AppStartStatus, startTime time.Time, metr IAppMetrics) {
	if metr == nil {
		return
	}

	_, nok := comp.(*nopComponent)
	if nok {
		return
	}

	_, pok := comp.(*ParallelComponent)
	if pok {
		return
	}

	_, sok := comp.(*SequentialComponent)
	if sok {
		return
	}

	metr.SetStartDuration(metrics.AppStartLabel{
		Status:    status,
		Component: comp.GetName(),
	}, time.Since(startTime))
}

// WriteStopMetric writes information about the component's stop to metrics.
func WriteStopMetric(comp IComponent, status metrics.AppStopStatus, startTime time.Time, metr IAppMetrics) {
	if metr == nil {
		return
	}

	_, nok := comp.(*nopComponent)
	if nok {
		return
	}

	_, pok := comp.(*ParallelComponent)
	if pok {
		return
	}

	_, sok := comp.(*SequentialComponent)
	if sok {
		return
	}

	metr.SetStopDuration(metrics.AppStopLabel{
		Status:    status,
		Component: comp.GetName(),
	}, time.Since(startTime))
}
