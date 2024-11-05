package middleware

//
//import (
//	"context"
//	"net/http"
//	"time"
//
//	"github.com/google/uuid"
//	"go.uber.org/zap"
//)
//
//type contextKey string
//
//const RequestIDKey = contextKey("request-id")
//
//type AccessLogMiddleware struct {
//	Logger *zap.Logger
//}
//
//func NewAccessLogMiddleware(logger *zap.Logger) *AccessLogMiddleware {
//	return &AccessLogMiddleware{Logger: logger}
//}
//
//func (alm *AccessLogMiddleware) Handler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		start := time.Now()
//
//		requestID := r.Header.Get("X-Request-ID")
//		if requestID == "" {
//			requestID = uuid.New().String()
//		}
//
//		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
//		r = r.WithContext(ctx)
//
//		w.Header().Set("X-Request-ID", requestID)
//
//		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
//
//		next.ServeHTTP(lrw, r)
//
//		duration := time.Since(start)
//
//		alm.Logger.Infow("HTTP Request",
//			"method", r.Method,
//			"url", r.URL.Path,
//			"remote_addr", r.RemoteAddr,
//			"status", lrw.statusCode,
//			"duration", duration,
//			"request_id", requestID,
//		)
//	})
//}
//
//type loggingResponseWriter struct {
//	http.ResponseWriter
//	statusCode int
//}
//
//func (lrw *loggingResponseWriter) WriteHeader(code int) {
//	lrw.statusCode = code
//	lrw.ResponseWriter.WriteHeader(code)
//}
