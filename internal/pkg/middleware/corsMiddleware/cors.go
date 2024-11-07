package corsMiddleware

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
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == http.MethodOptions {
			m.logger.Info("Handling request preflight", zap.String("method", r.Method), zap.String("url", r.URL.String()))
			return
		}
		m.logger.Info("Handling request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
		next.ServeHTTP(w, r)
	})
}
