package authcheck

import (
	"context"
	"net/http"
	"sparkit/internal/utils/consts"
	"time"
)

type sessionUsecase interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Middleware struct {
	usecase sessionUsecase
}

func New(usecase sessionUsecase) *Middleware {
	return &Middleware{usecase: usecase}
}

func (m *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(consts.SessionCookie)
		if err != nil {
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
			http.Error(w, "session is not valid", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
