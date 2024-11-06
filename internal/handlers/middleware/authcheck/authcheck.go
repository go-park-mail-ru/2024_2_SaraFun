package authcheck

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/utils/consts"
	"time"
)

type sessionUsecase interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Middleware struct {
	usecase sessionUsecase
	logger  *zap.Logger
}

func New(usecase sessionUsecase, logger *zap.Logger) *Middleware {
	return &Middleware{usecase: usecase, logger: logger}
}

func (m *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(consts.SessionCookie)
		if err != nil {
			m.logger.Error("bad getting cookie from request", zap.Error(err))
			http.Error(w, "session not found", http.StatusUnauthorized)
			return
		}
		userID, err := m.usecase.GetUserIDBySessionID(r.Context(), cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    consts.SessionCookie,
				Value:   "",
				Expires: time.Now().AddDate(0, 0, -1),
			})
			m.logger.Error("getting user id from session error", zap.Error(err))
			http.Error(w, "session is not valid", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		m.logger.Info("good authcheck")
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
