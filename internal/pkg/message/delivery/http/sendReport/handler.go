package sendReport

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	sessionClient generatedAuth.AuthClient
	messageClient generatedMessage.MessageClient
	logger        *zap.Logger
}

func NewHandler(sessionClient generatedAuth.AuthClient, messageClient generatedMessage.MessageClient, logger *zap.Logger) *Handler {
	return &Handler{
		sessionClient: sessionClient,
		messageClient: messageClient,
		logger:        logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report := models.Report{}
	report.Sanitize()
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		h.logger.Error("bad json decode", zap.Error(err))
		http.Error(w, "bad json decode", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "bad cookie", http.StatusBadRequest)
		return
	}
	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Error("bad get user id", zap.Error(err))
		http.Error(w, "bad get user id", http.StatusBadRequest)
	}
	report.Author = int(userId.UserId)

	requestReport := &generatedMessage.Report{
		ID:       int32(report.ID),
		Author:   int32(report.Author),
		Receiver: int32(report.Receiver),
		Body:     report.Body,
	}
	addReportRequest := &generatedMessage.AddReportRequest{Report: requestReport}
	_, err = h.messageClient.AddReport(ctx, addReportRequest)
	if err != nil {
		h.logger.Error("bad add report", zap.Error(err))
		http.Error(w, "bad add report", http.StatusInternalServerError)
		return
	}
	h.logger.Info("sendReport successfully")
}
