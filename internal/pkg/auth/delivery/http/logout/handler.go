package logout

import (
	"context"
	"fmt"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
type SessionService interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

type SessionClient interface {
	DeleteSession(ctx context.Context, in *generatedAuth.DeleteSessionRequest) (*generatedAuth.DeleteSessionResponse, error)
}

type Handler struct {
	client generatedAuth.AuthClient
	logger *zap.Logger
}

func NewHandler(client generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{client: client, logger: logger}
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
	if _, err := h.client.DeleteSession(ctx, &generatedAuth.DeleteSessionRequest{SessionID: cookie.Value}); err != nil {
		h.logger.Error("logout client delete session error", zap.Error(err))
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    consts.SessionCookie,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	h.logger.Info("session deleted success", zap.String("session", cookie.Value))
	fmt.Fprintf(w, "Вы успешно вышли из учетной записи")
}
