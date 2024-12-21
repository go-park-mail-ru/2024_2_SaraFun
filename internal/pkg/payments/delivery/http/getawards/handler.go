package getawards

import (
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

type Response struct {
	Responses []*models.Award `json:"responses"`
}

type Handler struct {
	authClient     generatedAuth.AuthClient
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient,
	paymentsClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:     authClient,
		paymentsClient: paymentsClient,
		logger:         logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "bad cookie", http.StatusUnauthorized)
		return
	}
	getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	_, err = h.authClient.GetUserIDBySessionID(ctx, getUserIDReq)
	if err != nil {
		h.logger.Error("get user id by session id", zap.Error(err))
		http.Error(w, "get user id by session id", http.StatusUnauthorized)
		return
	}

	getAwards := &generatedPayments.GetAwardsRequest{}

	grpcAwards, err := h.paymentsClient.GetAwards(ctx, getAwards)
	if err != nil {
		h.logger.Error("get awards", zap.Error(err))
		http.Error(w, "get awards", http.StatusInternalServerError)
		return
	}

	var awards []*models.Award
	for _, awd := range grpcAwards.Awards {
		award := &models.Award{
			DayNumber: int(awd.DayNumber),
			Type:      awd.Type,
			Count:     int(awd.Count),
		}
		awards = append(awards, award)
	}
	w.Header().Set("Content-Type", "application/json")
	response := Response{Responses: awards}

	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("bad marshal response", zap.Error(err))
		http.Error(w, "marshal response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write response", zap.Error(err))
		http.Error(w, "write response error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("get awards success")

}
