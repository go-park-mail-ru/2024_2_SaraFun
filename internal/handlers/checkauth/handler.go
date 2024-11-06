package checkauth

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=checkauth_mocks . SessionService
type SessionService interface {
	CheckSession(ctx context.Context, sessionID string) error
}

type Handler struct {
	sessionService SessionService
	logger         *zap.Logger
}

func NewHandler(service SessionService, logger *zap.Logger) *Handler {
	return &Handler{sessionService: service, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	if r.Method != http.MethodGet {
		h.logger.Error("unexpected method", zap.String("method", r.Method))
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad getting cookie", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	check := h.sessionService.CheckSession(ctx, cookie.Value)
	if check != nil {
		h.logger.Error("session check", zap.String("session", cookie.Value), zap.Error(check))
		http.Error(w, "session is not valid", http.StatusUnauthorized)
		return
	}
	h.logger.Info("good session check", zap.String("session", cookie.Value))
	fmt.Fprintf(w, "ok")
}
