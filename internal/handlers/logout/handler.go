package logout

import (
	"context"
	"fmt"
	"net/http"
	"sparkit/internal/utils/consts"
	"time"
)

type SessionService interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

type Handler struct {
	service SessionService
}

func NewHandler(service SessionService) *Handler {
	return &Handler{service: service}
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
	if err := h.service.DeleteSession(ctx, cookie.Value); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    consts.SessionCookie,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	fmt.Fprintf(w, "log out is complete")
}
