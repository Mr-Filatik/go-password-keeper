// Package metrics provides functionality for working with metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProviderExperiment struct {
	distributionsTotal prometheus.CounterVec
	executionTotal     prometheus.CounterVec
}

const subsystemExperiment = "experiment"

func createProviderExperiment(namespace string, constLabels prometheus.Labels) *ProviderExperiment {
	provider := &ProviderExperiment{
		distributionsTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				ConstLabels: constLabels,
				Namespace:   namespace,
				Subsystem:   subsystemExperiment,
				Name:        "distributions_total",
				Help:        "Total number of distributions.",
			},
			ExperimentDistributionLabelNames(),
		),
		executionTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				ConstLabels: constLabels,
				Namespace:   namespace,
				Subsystem:   subsystemExperiment,
				Name:        "executions_total",
				Help:        "Total number of executions.",
			},
			ExperimentExecutionLabelNames(),
		),
	}

	prometheus.MustRegister(
		provider.distributionsTotal,
		provider.executionTotal,
	)

	return provider
}

type ExperimentDistributionLabel struct {
	ExperimentName string

	BranchName string

	Distributor string
}

func ExperimentDistributionLabelNames() []string {
	return []string{"experiment_name", "branch_name", "distributor"}
}

type ExperimentExecutionLabel struct {
	ExperimentName string

	BranchName string
}

func ExperimentExecutionLabelNames() []string {
	return []string{"experiment_name", "branch_name"}
}

func (p *ProviderExperiment) IncDistributionTotal(labels ExperimentDistributionLabel) {
	lbls := prometheus.Labels{
		"experiment_name": labels.ExperimentName,
		"branch_name":     labels.BranchName,
		"distributor":     labels.Distributor,
	}

	p.distributionsTotal.With(lbls).Inc()
}

func (p *ProviderExperiment) IncExecutionTotal(labels ExperimentExecutionLabel) {
	lbls := prometheus.Labels{
		"experiment_name": labels.ExperimentName,
		"branch_name":     labels.BranchName,
	}

	p.executionTotal.With(lbls).Inc()
}
