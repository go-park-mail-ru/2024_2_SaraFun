package CSPMiddleware

import (
	"go.uber.org/zap"
	"net/http"
)

type Middleware struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Middleware {
	return &Middleware{logger: logger}
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';base-uri 'self';form-action 'self'")
		m.logger.Info("Handling request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
		next.ServeHTTP(w, r)
	})
}
