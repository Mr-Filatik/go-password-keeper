// Package metrics provides functionality for working with metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// BaseMetrics provides a base type for working with metrics.
//
// Used strictly through embedding (BaseMetrics) into other types.
type BaseMetrics struct {
	namespace   string
	constLabels prometheus.Labels
	reg         prometheus.Registerer
}

// NewBaseMetrics creates a new BaseMetrics instance.
//
// By default, it uses the prometheus.DefaultRegisterer registrar.
// You can set your own using the SetRegisterer(reg prometheus.Registerer) function.
//
// Parameters:
// - namespace string: namespace for metrics;
// - constLabels map[string]string: global labels for metrics.
func NewBaseMetrics(namespace string, constLabels map[string]string) *BaseMetrics {
	return &BaseMetrics{
		namespace:   namespace,
		constLabels: cloneLabels(constLabels),
		reg:         prometheus.DefaultRegisterer,
	}
}

// SetRegisterer allows you to install a custom registrar.
//
// Parameters:
// - reg prometheus.Registerer: the registrar for the metrics collection (default is prometheus.DefaultRegisterer).
func (b *BaseMetrics) SetRegisterer(reg prometheus.Registerer) *BaseMetrics {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	b.reg = reg

	return b
}

// CommonOpt describes the general parameters required to create a metric.
type CommonOpt struct {
	// Subsystem - the name of the subsystem.
	//
	// The string must contain underscores instead of spaces.
	Subsystem string

	// Name - metric name.
	//
	// The string must contain underscores instead of spaces.
	//
	// A special rule applies to counters:
	// Must be specified in the plural and with the suffix _total.
	// Examples: experiments_total, requests_total.
	Name string

	// Help - information about the metric.
	Help string

	// LabelNames - label names for the metric.
	LabelNames []string
}

// CounterOpt describes the parameters needed to create a counter.
type CounterOpt struct {
	CommonOpt
}

// CreateCounter creates a new metrics counter.
//
// Parameters:
// - opt CounterOpt: parameters for creating a counter.
func (b *BaseMetrics) CreateCounter(opt CounterOpt) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			ConstLabels: b.constLabels,
			Namespace:   b.namespace,
			Subsystem:   opt.Subsystem,
			Name:        opt.Name,
			Help:        opt.Help,
		},
		cloneLabelNames(opt.LabelNames),
	)

	//nolint:godox
	// TODO: MustRegister panics when re-registering a metric.
	b.reg.MustRegister(counter)

	return counter
}

// Gauge

// HistogramOpt describes the parameters needed to create a histogram.
type HistogramOpt struct {
	CommonOpt

	// Buckets - buckets for the histogram.
	//
	// If you don't specify a value, prometheus.DefBuckets will be selected by default.
	Buckets []float64
}

// CreateHistogram creates a new metrics histogram.
//
// Parameters:
// - opt HistogramOpt: parameters for creating a histogram.
func (b *BaseMetrics) CreateHistogram(opt HistogramOpt) *prometheus.HistogramVec {
	if opt.Buckets == nil {
		opt.Buckets = prometheus.DefBuckets
	}

	histogram := prometheus.NewHistogramVec(
		//nolint:exhaustruct // Native Histogram is not used
		prometheus.HistogramOpts{
			ConstLabels: b.constLabels,
			Namespace:   b.namespace,
			Subsystem:   opt.Subsystem,
			Name:        opt.Name,
			Help:        opt.Help,
			Buckets:     opt.Buckets,
		},
		cloneLabelNames(opt.LabelNames),
	)

	//nolint:godox
	// TODO: MustRegister panics when re-registering a metric.
	b.reg.MustRegister(histogram)

	return histogram
}

// Summary/Native Histogram

// cloneLabels creates a copy of the passed map[string]string (prometheus.Labels),
// to avoid side effects from modifying it after passing.
func cloneLabels(src map[string]string) prometheus.Labels {
	if src == nil {
		return prometheus.Labels{}
	}

	dst := make(prometheus.Labels, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

// cloneLabelNames creates a copy of the slice of label names.
// This prevents the original array from being accidentally modified from the outside.
func cloneLabelNames(src []string) []string {
	if src == nil {
		return nil
	}

	dst := make([]string, len(src))
	copy(dst, src)

	return dst
}
