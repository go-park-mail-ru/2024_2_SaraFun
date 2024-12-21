package setConnection

import (
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	websocket "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_usecase.go -package=mocks UseCase
//go:generate mockgen -destination=./mocks/mock_authClient.go -package=mocks AuthClient

type UseCase interface {
	AddConnection(ctx context.Context, conn *websocket.Conn, userId int) error
	DeleteConnection(ctx context.Context, userId int) error
}

type Handler struct {
	useCase    UseCase
	authClient generatedAuth.AuthClient
	logger     *zap.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
		return
	}
	h.logger.Info("before upgrade")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Info("setConnection websocket upgrade error", zap.Error(err))
		http.Error(w, "upgrade websocket error", http.StatusInternalServerError)
		return
	}
	err = h.useCase.AddConnection(ctx, ws, int(userId.UserId))
	if err != nil {
		h.logger.Info("setConnection usecase add connection error", zap.Error(err))
		http.Error(w, "add connection error", http.StatusInternalServerError)
		return
	}

	go func() {
		defer func() {
			h.logger.Info("deleteConnection defer func")
			err := h.useCase.DeleteConnection(ctx, int(userId.UserId))
			if err != nil {
				h.logger.Error("setConnection delete connection error", zap.Error(err))
			}
		}()

		for {
			messageType, msg, err := ws.ReadMessage()
			if err != nil {
				h.logger.Info("setConnection websocket read error", zap.Error(err))
				break
			}
			if messageType == websocket.CloseMessage {
				h.logger.Info("setConnection websocket close", zap.String("message_type", string(msg)))
				break
			}
		}
	}()
}
