package middleware

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

//type contextKey string
//
//const RequestIDKey = contextKey("request-id")

type AccessLogMiddleware struct {
	Logger *zap.SugaredLogger
}

func NewAccessLogMiddleware(logger *zap.SugaredLogger) *AccessLogMiddleware {
	return &AccessLogMiddleware{Logger: logger}
}

func (alm *AccessLogMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), consts.RequestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		alm.Logger.Infow("HTTP Request",
			"method", r.Method,
			"url", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"status", lrw.statusCode,
			"duration", duration,
			"request_id", requestID,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not support Hijacker")
	}
	return hijacker.Hijack()
}
