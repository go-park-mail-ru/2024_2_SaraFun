package logout

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/utils/consts"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
type SessionService interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

type Handler struct {
	service SessionService
	logger  *zap.Logger
}

func NewHandler(service SessionService, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	if r.Method != http.MethodGet {
		h.logger.Error("bad method", zap.String("method", r.Method))
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad getting cookie from request", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	if err := h.service.DeleteSession(ctx, cookie.Value); err != nil {
		h.logger.Error("deleting session", zap.Error(err))
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    consts.SessionCookie,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	h.logger.Info("session deleted success", zap.String("session", cookie.Value))
	fmt.Fprintf(w, "log out is complete")
}
