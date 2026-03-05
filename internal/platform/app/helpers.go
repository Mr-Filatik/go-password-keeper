package app

import (
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/metrics"
)

// IAppMetrics описывает интерфейс для метрик о старте и остановке компонентов.
// Может быть перенести, но вроде бы это по гошному.
type IAppMetrics interface {
	SetStartDuration(labels metrics.AppStartLabel, duration time.Duration)
	SetStopDuration(labels metrics.AppStopLabel, duration time.Duration)
}

// WriteStartMetric writes information about the component's start to metrics.
func WriteStartMetric(comp IComponent, status metrics.AppStartStatus, startTime time.Time, metr IAppMetrics) {
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
