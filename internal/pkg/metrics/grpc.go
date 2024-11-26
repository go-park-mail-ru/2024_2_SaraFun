package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type GrpcMetrics struct {
	HitsTotal *prometheus.CounterVec
	Errors    *prometheus.CounterVec
	Times     *prometheus.HistogramVec
	name      string
}

func NewGrpcMetrics(name string) (*GrpcMetrics, error) {
	var metr GrpcMetrics
	metr.HitsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_hits_total",
		Help: "Number of total hits",
	},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, fmt.Errorf("grpc metrics: unable to register hits total metric: %w", err)
	}

	metr.Errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_errors_total",
		Help: "number of total errors",
	},
		[]string{"path", "service"},
	)

	if err := prometheus.Register(metr.Errors); err != nil {
		return nil, fmt.Errorf("grpc metrics: unable to register errors metric: %w", err)
	}

	metr.name = name

	metr.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "grpc_total_times",
	},
		[]string{"path"},
	)

	if err := prometheus.Register(metr.Times); err != nil {
		return nil, fmt.Errorf("grpc metrics: unable to register times metric: %w", err)
	}
	return &metr, nil
}

func (m *GrpcMetrics) IncreaseHits(path string) {
	m.HitsTotal.WithLabelValues(path, m.name).Inc()
}

func (m *GrpcMetrics) IncreaseErrors(path string) {
	m.Errors.WithLabelValues(path, m.name).Inc()
}

func (m *GrpcMetrics) ObserveResponseTime(path string, observeTime float64) {
	m.Times.WithLabelValues(path).Observe(observeTime)
}
