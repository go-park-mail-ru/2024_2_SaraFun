package sendmessage

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

//type MessageService interface {
//	AddMessage(ctx context.Context, message *models.Message) (int, error)
//}

//type SessionService interface {
//	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
//}

type WebSocketService interface {
	WriteMessage(ctx context.Context, authorID int, receiverID int, message string) error
}

type Handler struct {
	messageClient        generatedMessage.MessageClient
	sessionClient        generatedAuth.AuthClient
	communicationsClient generatedCommunications.CommunicationsClient
	ws                   WebSocketService
	logger               *zap.Logger
}

func NewHandler(messageClient generatedMessage.MessageClient, ws WebSocketService, sessionClient generatedAuth.AuthClient, communicationsClient generatedCommunications.CommunicationsClient, logger *zap.Logger) *Handler {
	return &Handler{
		messageClient:        messageClient,
		ws:                   ws,
		sessionClient:        sessionClient,
		communicationsClient: communicationsClient,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	msg := &models.Message{}
	msg.Sanitize()
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		h.logger.Info("Error decoding message", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Info("Error getting session cookie", zap.Error(err))
		http.Error(w, "bad request", http.StatusUnauthorized)
		return
	}

	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Info("Error getting user ID", zap.Error(err))
		http.Error(w, "bad request", http.StatusUnauthorized)
		return
	}
	msg.Author = int(userId.UserId)
	reqMessage := &generatedMessage.ChatMessage{
		ID:       int32(msg.ID),
		Author:   int32(msg.Author),
		Receiver: int32(msg.Receiver),
		Body:     msg.Body,
	}
	addMessageRequest := &generatedMessage.AddMessageRequest{Message: reqMessage}
	_, err = h.messageClient.AddMessage(ctx, addMessageRequest)
	if err != nil {
		h.logger.Info("Error adding message", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	err = h.ws.WriteMessage(ctx, msg.Author, msg.Receiver, msg.Body)
	if err != nil {
		h.logger.Info("Error writing message", zap.Error(err))
		//http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
