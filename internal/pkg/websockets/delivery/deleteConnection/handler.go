package deleteConnection

import (
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_usecase.go -package=mocks  UseCase
//go:generate mockgen -destination=./mocks/mock_authClient.go -package=mocks  AuthClient

type UseCase interface {
	DeleteConnection(ctx context.Context, userId int) error
}

type Handler struct {
	useCase    UseCase
	authClient generatedAuth.AuthClient
	logger     *zap.Logger
}

func NewHandler(uc UseCase, authClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{
		useCase:    uc,
		authClient: authClient,
		logger:     logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("setConnection request id", zap.String("request_id", req_id))

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Info("setConnection get cookie error", zap.Error(err))
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return
	}
	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.authClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Info("setConnection get user id error", zap.Error(err))
		http.Error(w, "get user id error", http.StatusUnauthorized)
	}
	err = h.useCase.DeleteConnection(ctx, int(userId.UserId))
	if err != nil {
		h.logger.Info("setConnection usecase add connection error", zap.Error(err))
		http.Error(w, "add connection error", http.StatusInternalServerError)
		return
	}
}
