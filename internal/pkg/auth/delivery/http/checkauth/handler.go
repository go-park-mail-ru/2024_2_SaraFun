package checkauth

import (
	"context"
	"encoding/json"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=checkauth_mocks . SessionService

type SessionClient interface {
	CheckSession(ctx context.Context, in *generatedAuth.CheckSessionRequest) (*generatedAuth.CheckSessionResponse, error)
}

type Handler struct {
	sessionClient  generatedAuth.AuthClient
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(service generatedAuth.AuthClient,
	paymentsClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{sessionClient: service, paymentsClient: paymentsClient, logger: logger}
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
	_, err = h.sessionClient.CheckSession(ctx, &generatedAuth.CheckSessionRequest{SessionID: cookie.Value})
	if err != nil {
		h.logger.Error("checkauth sessionClient check session error", zap.Error(err))
		http.Error(w, "bad session", http.StatusUnauthorized)
	}

	cookie, err = r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("cookie error", zap.Error(err))
		http.Error(w, "Вы не авторизованы!", http.StatusUnauthorized)
		return
	}
	getIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getIDRequest)

	updateActivityReq := &generatedPayments.UpdateActivityRequest{UserID: userId.UserId}
	ans, err := h.paymentsClient.UpdateActivity(ctx, updateActivityReq)
	if err != nil {
		h.logger.Error("checkauth paymentsClient update activity error", zap.Error(err))
		http.Error(w, "что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(ans.Answer)
	if err != nil {
		h.logger.Error("json marshal error", zap.Error(err))
		http.Error(w, "что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write response error", zap.Error(err))
		http.Error(w, "что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	h.logger.Info("good update activity", zap.String("session", cookie.Value))
}
