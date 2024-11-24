package grpcMetricsMiddleware

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type Middleware struct {
	metrics *metrics.GrpcMetrics
	logger  *zap.Logger
}

func NewMiddleware(metrics *metrics.GrpcMetrics, logger *zap.Logger) *Middleware {
	return &Middleware{
		metrics: metrics,
		logger:  logger,
	}
}

func (m *Middleware) ServerMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	tm := time.Since(start).Seconds()
	if err != nil {
		m.metrics.IncreaseErrors(info.FullMethod)
	}
	m.metrics.IncreaseHits(info.FullMethod)
	m.metrics.ObserveResponseTime(info.FullMethod, tm)
	return resp, err
}
