package sendReport

import (
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	sessionClient        generatedAuth.AuthClient
	messageClient        generatedMessage.MessageClient
	communicationsClient generatedCommunications.CommunicationsClient
	logger               *zap.Logger
}

func NewHandler(sessionClient generatedAuth.AuthClient, messageClient generatedMessage.MessageClient,
	communicationsClient generatedCommunications.CommunicationsClient,
	logger *zap.Logger) *Handler {
	return &Handler{
		sessionClient:        sessionClient,
		messageClient:        messageClient,
		communicationsClient: communicationsClient,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report := models.Report{}
	report.Sanitize()
	err := easyjson.UnmarshalFromReader(r.Body, &report)
	if err != nil {
		h.logger.Error("bad json decode", zap.Error(err))
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusBadRequest)
		return
	}
	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Error("bad get user id", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusBadRequest)
		return
	}
	report.Author = int(userId.UserId)
	h.logger.Info("report", zap.Any("report", report))
	requestReport := &generatedMessage.Report{
		ID:       int32(report.ID),
		Author:   int32(report.Author),
		Receiver: int32(report.Receiver),
		Reason:   report.Reason,
		Body:     report.Body,
	}
	addReportRequest := &generatedMessage.AddReportRequest{Report: requestReport}
	_, err = h.messageClient.AddReport(ctx, addReportRequest)
	if err != nil {
		h.logger.Error("bad add report", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	reqReaction := &generatedCommunications.Reaction{
		Author:   userId.UserId,
		Receiver: int32(report.Receiver),
		Type:     false,
	}
	updCreateRequest := &generatedCommunications.UpdateOrCreateReactionRequest{Reaction: reqReaction}
	_, err = h.communicationsClient.UpdateOrCreateReaction(ctx, updCreateRequest)
	if err != nil {
		h.logger.Error("bad update or create reaction", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	h.logger.Info("sendReport successfully")
}
