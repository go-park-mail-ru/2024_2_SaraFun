package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type DatabaseMetrics struct {
	db      string
	service string
	Errors  *prometheus.CounterVec
	Times   *prometheus.HistogramVec
}

func NewDatabaseMetrics(name string, service string) (*DatabaseMetrics, error) {
	m := &DatabaseMetrics{db: name, service: service}
	m.Errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "db_errors_total",
		Help: "number of total errors."},
		[]string{"query", "db", "service"})
	if err := prometheus.Register(m.Errors); err != nil {
		return nil, fmt.Errorf("db metrics: error registering errors collector: %w", err)
	}
	m.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: name + "db_total_times"},
		[]string{"query", "db", "service"},
	)
	if err := prometheus.Register(m.Times); err != nil {
		return nil, fmt.Errorf("db metrics: error registering times collector: %w", err)
	}
	return m, nil
}

func (m *DatabaseMetrics) IncreaseErrors(queryName string) {
	m.Errors.WithLabelValues(queryName, m.db, m.service).Inc()
}

func (m *DatabaseMetrics) ObserveResponseTime(queryName string, observeTime float64) {
	m.Times.WithLabelValues(queryName, m.db, m.service).Observe(observeTime)
}
