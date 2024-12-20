package sendmessage

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

//go:generate mockgen -destination=./mocks/mock_WebSocketService.go -package=sign_up_mocks . WebSocketService
type WebSocketService interface {
	WriteMessage(ctx context.Context, authorID int, receiverID int, message string, username string) error
}

type ErrResponse struct {
	Info string `json:"info"`
}

//easyjson:skip
type Handler struct {
	messageClient        generatedMessage.MessageClient
	sessionClient        generatedAuth.AuthClient
	communicationsClient generatedCommunications.CommunicationsClient
	personalitiesClient  generatedPersonalities.PersonalitiesClient
	ws                   WebSocketService
	logger               *zap.Logger
}

func NewHandler(messageClient generatedMessage.MessageClient,
	ws WebSocketService, sessionClient generatedAuth.AuthClient,
	communicationsClient generatedCommunications.CommunicationsClient,
	personalitiesClient generatedPersonalities.PersonalitiesClient, logger *zap.Logger) *Handler {
	return &Handler{
		messageClient:        messageClient,
		ws:                   ws,
		sessionClient:        sessionClient,
		communicationsClient: communicationsClient,
		personalitiesClient:  personalitiesClient,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	msg := &models.Message{}
	msg.Sanitize()
	if err := easyjson.UnmarshalFromReader(r.Body, msg); err != nil {
		h.logger.Info("Error decoding message", zap.Error(err))
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Info("Error getting session cookie", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}

	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Info("Error getting user ID", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}

	getReportRequest := &generatedMessage.CheckUsersBlockNotExistsRequest{
		FirstUserID:  userId.UserId,
		SecondUserID: int32(msg.Receiver),
	}
	status, err := h.messageClient.CheckUsersBlockNotExists(ctx, getReportRequest)
	h.logger.Info("status", zap.Any("status", status))
	if err != nil {
		h.logger.Info("Error checking users block exists", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	if status.Status != "" {
		errResponse := &ErrResponse{
			Info: status.Status,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _, err := easyjson.MarshalToHTTPResponseWriter(errResponse, w)
		if err != nil {
			h.logger.Info("Error encoding message", zap.Error(err))
			http.Error(w, "что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
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
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: userId.UserId}
	username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
	if err != nil {
		h.logger.Info("Error getting username by userID", zap.Error(err))
		return
	}
	err = h.ws.WriteMessage(ctx, msg.Author, msg.Receiver, msg.Body, username.Username)
	if err != nil {
		h.logger.Info("Error writing message", zap.Error(err))
		return
	}
}
