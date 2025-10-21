// Package metrics provides functionality for working with metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ExperimentMetrics provides a type for working with AB experiment metrics.
type ExperimentMetrics struct {
	BaseMetrics

	distributionsCounter *prometheus.CounterVec
	executionsCounter    *prometheus.CounterVec
}

// NewExperimentMetrics creates a new ExperimentMetrics instance.
//
// Parameters:
//   - base BaseMetrics: a basic metric type that contains common data.
func NewExperimentMetrics(base BaseMetrics) *ExperimentMetrics {
	subsystemName := "experiment"

	distributionsCounter := base.CreateCounter(CounterOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "distributions_total",
			Help:       "Total number of experiment distributions",
			LabelNames: []string{"experiment_name", "branch_name", "distributor"},
		},
	})

	executionsCounter := base.CreateCounter(CounterOpt{
		CommonOpt: CommonOpt{
			Subsystem:  subsystemName,
			Name:       "executions_total",
			Help:       "Total number of experiment executions",
			LabelNames: []string{"experiment_name", "branch_name", "executor", "status"},
		},
	})

	return &ExperimentMetrics{
		BaseMetrics:          base,
		distributionsCounter: distributionsCounter,
		executionsCounter:    executionsCounter,
	}
}

// ExperimentDistributionLabel describes the data required to record the metric.
type ExperimentDistributionLabel struct {
	// ExperimentName - experiment name.
	ExperimentName string

	// BranchName - experiment branch name.
	BranchName string

	// Distributor - distributor.
	Distributor string
}

// IncDistributionsCounter increments the counter by one, specifying the labels.
//
// Parameters:
//   - labels ExperimentDistributionLabel: labels.
func (p *ExperimentMetrics) IncDistributionsCounter(labels ExperimentDistributionLabel) {
	lbls := prometheus.Labels{
		"experiment_name": labels.ExperimentName,
		"branch_name":     labels.BranchName,
		"distributor":     labels.Distributor,
	}

	p.distributionsCounter.With(lbls).Inc()
}

// ExperimentExecutionLabel describes the data required to record the metric.
type ExperimentExecutionLabel struct {
	// ExperimentName - experiment name.
	ExperimentName string

	// BranchName - experiment branch name.
	BranchName string

	// Executor - executor.
	Executor string

	// Status - status (success or failed).
	Status string
}

// IncExecutionsCounter increments the counter by one, specifying the labels.
//
// Parameters:
//   - labels ExperimentExecutionLabel: labels.
func (p *ExperimentMetrics) IncExecutionsCounter(labels ExperimentExecutionLabel) {
	lbls := prometheus.Labels{
		"experiment_name": labels.ExperimentName,
		"branch_name":     labels.BranchName,
		"executor":        labels.Executor,
		"status":          labels.Status,
	}

	p.executionsCounter.With(lbls).Inc()
}
