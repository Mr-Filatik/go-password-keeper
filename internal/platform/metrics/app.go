// Package metrics provides functionality for working with metrics.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// AppMetrics provides a type for working with core application metrics.
type AppMetrics struct {
	BaseMetrics

	buildInfoGauge  *prometheus.GaugeVec // information about the current build of the application
	deployInfoGauge *prometheus.GaugeVec // information about the current deploy of the application

	startInfoGauge *prometheus.GaugeVec // information about launching application components
	stopInfoGauge  *prometheus.GaugeVec // information about stopping application components
}

// NewAppMetrics creates a new *AppMetrics instance and initializes the base metrics.
//
// Parameters:
//   - base BaseMetrics: a basic metric type that contains common data.
//
// List of basic app metrics:
//   - build_info: Information about the application build;
//   - deploy_info: Information about application deploy;
//   - start_duration_seconds: Duration of all services startup in seconds;
//   - stop_duration_seconds: Duration of stopping all services in seconds.
func NewAppMetrics(base BaseMetrics) *AppMetrics {
	subsystemName := "app"

	buildInfoGauge := base.CreateGauge(GaugeOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "build_info",
			Help:       "Information about the application build.",
			LabelNames: []string{"version", "date", "commit"},
		},
	})

	deployInfoGauge := base.CreateGauge(GaugeOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "deploy_info",
			Help:       "Information about application deploy.",
			LabelNames: []string{"number"},
		},
	})

	startInfoGauge := base.CreateGauge(GaugeOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "start_duration_seconds",
			Help:       "Duration of all services startup in seconds.",
			LabelNames: []string{"status"},
		},
	})

	stopInfoGauge := base.CreateGauge(GaugeOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "stop_duration_seconds",
			Help:       "Duration of stopping all services in seconds.",
			LabelNames: []string{"status"},
		},
	})

	return &AppMetrics{
		BaseMetrics:     base,
		buildInfoGauge:  buildInfoGauge,
		deployInfoGauge: deployInfoGauge,
		startInfoGauge:  startInfoGauge,
		stopInfoGauge:   stopInfoGauge,
	}
}

// AppBuildLabel labels for describing information about the application build.
type AppBuildLabel struct {
	Version string
	Date    string
	Commit  string
}

// SetBuildInfo sets information about a build by specifying labels.
//
// Parameters:
//   - labels AppBuildLabel: labels.
func (p *AppMetrics) SetBuildInfo(labels AppBuildLabel) {
	lbls := prometheus.Labels{
		"version": labels.Version,
		"date":    labels.Date,
		"commit":  labels.Commit,
	}

	p.buildInfoGauge.With(lbls).Add(1)
}

// AppDeployLabel labels for describing information about the application deploy.
type AppDeployLabel struct {
	Number string
}

// SetDeployInfo sets information about a deploy by specifying labels.
//
// Parameters:
//   - labels AppDeployLabel: labels.
func (p *AppMetrics) SetDeployInfo(labels AppDeployLabel) {
	lbls := prometheus.Labels{
		"number": labels.Number,
	}

	p.deployInfoGauge.With(lbls).Add(1)
}

// AppStartStatus describes the statuses when starting application components.
type AppStartStatus string

const (
	// StartStatusSuccess - components launched without errors.
	StartStatusSuccess AppStartStatus = "success"

	// StartStatusFailed - components started with errors.
	StartStatusFailed AppStartStatus = "failed"
)

// AppStartLabel labels for describing information about the startup of application components.
type AppStartLabel struct {
	Status AppStartStatus // component startup status
}

// SetStartDuration sets the startup time of components by specifying labels.
//
// Parameters:
//   - labels AppStartLabel: labels.
//   - duration time.Duration: labels.
func (p *AppMetrics) SetStartDuration(labels AppStartLabel, duration time.Duration) {
	lbls := prometheus.Labels{
		"status": string(labels.Status),
	}

	p.startInfoGauge.With(lbls).Add(duration.Seconds())
}

// AppStopStatus describes the statuses when application components are stopped.
type AppStopStatus string

const (
	// StopStatusSuccess - components stopped without errors.
	StopStatusSuccess AppStopStatus = "success"

	// StopStatusFailed - components stopped with errors.
	StopStatusFailed AppStopStatus = "failed"
)

// AppStopLabel labels for describing information about stopping application components.
type AppStopLabel struct {
	Status AppStopStatus // component stop status
}

// SetStopDuration sets the stop time of components by specifying labels.
//
// Parameters:
//   - labels AppStartLabel: labels.
//   - duration time.Duration: labels.
func (p *AppMetrics) SetStopDuration(labels AppStopLabel, duration time.Duration) {
	lbls := prometheus.Labels{
		"status": string(labels.Status),
	}

	p.stopInfoGauge.With(lbls).Add(duration.Seconds())
}
