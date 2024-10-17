package checkauth

import (
	"context"
	"fmt"
	"net/http"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=checkauth_mocks . SessionService
type SessionService interface {
	CheckSession(ctx context.Context, sessionID string) error
}

type Handler struct {
	sessionService SessionService
}

func NewHandler(service SessionService) *Handler {
	return &Handler{sessionService: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	check := h.sessionService.CheckSession(ctx, cookie.Value)
	if check != nil {
		http.Error(w, "session is not valid", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "ok")
}
