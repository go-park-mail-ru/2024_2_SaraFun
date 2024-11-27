package httpMetricsMiddleware

import (
	"bufio"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/metrics"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

type Middleware struct {
	metrics *metrics.HttpMetrics
	logger  *zap.Logger
}

func NewMiddleware(metrics *metrics.HttpMetrics, logger *zap.Logger) *Middleware {
	return &Middleware{
		metrics: metrics,
		logger:  logger,
	}
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		route := r.URL.Path
		statusCode := rw.statusCode
		if statusCode != http.StatusOK {
			m.metrics.IncreaseErrors(route, strconv.Itoa(statusCode))
		}
		m.metrics.IncreaseHits(route, strconv.Itoa(statusCode))
		m.metrics.ObserveResponseTime(statusCode, route, time.Since(start).Seconds())
	})
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not support Hijacker")
	}
	return hijacker.Hijack()
}
