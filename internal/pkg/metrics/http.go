package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type HttpMetrics struct {
	HitsTotal *prometheus.CounterVec
	name      string
	Times     *prometheus.HistogramVec
	Errors    *prometheus.CounterVec
}

func NewHttpMetrics(name string) (*HttpMetrics, error) {
	var metr HttpMetrics
	metr.HitsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_hits_total",
		Help: "Number of total hits",
	},
		[]string{"path", "service", "status"},
	)
	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, fmt.Errorf("http metrics: unable to register hits total metric: %w", err)
	}

	metr.Errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "number of total errors",
	},
		[]string{"path", "service", "status"},
	)

	if err := prometheus.Register(metr.Errors); err != nil {
		return nil, fmt.Errorf("http metrics: unable to register errors metric: %w", err)
	}

	metr.name = name

	metr.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_total_times",
	},
		[]string{"path", "status"},
	)

	if err := prometheus.Register(metr.Times); err != nil {
		return nil, fmt.Errorf("http metrics: unable to register times metric: %w", err)
	}
	return &metr, nil

}
func (m *HttpMetrics) IncreaseHits(path string, status string) {
	m.HitsTotal.WithLabelValues(path, m.name, status).Inc()
}

func (m *HttpMetrics) IncreaseErrors(path string, status string) {
	m.Errors.WithLabelValues(path, m.name, status).Inc()
}

func (m *HttpMetrics) ObserveResponseTime(status int, path string, observeTime float64) {
	m.Times.WithLabelValues(path, strconv.Itoa(status)).Observe(observeTime)
}
