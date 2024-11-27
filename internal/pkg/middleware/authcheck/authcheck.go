package authcheck

import (
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type sessionUsecase interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Middleware struct {
	usecase generatedAuth.AuthClient
	logger  *zap.Logger
}

func New(usecase generatedAuth.AuthClient, logger *zap.Logger) *Middleware {
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
		getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
		userID, err := m.usecase.GetUserIDBySessionID(r.Context(), getUserIdRequest)
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

		ctx := context.WithValue(r.Context(), "userID", userID.UserId)
		m.logger.Info("good authcheck")
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
