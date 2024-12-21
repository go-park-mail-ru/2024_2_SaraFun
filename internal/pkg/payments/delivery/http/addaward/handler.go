package addaward

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(paymentsClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{
		paymentsClient: paymentsClient,
		logger:         logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var award models.Award
	err := easyjson.UnmarshalFromReader(r.Body, &award)
	if err != nil {
		h.logger.Error("easyjson bad parsing from body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	reqAward := &generatedPayments.Award{
		DayNumber: int32(award.DayNumber),
		Type:      award.Type,
		Count:     int32(award.Count),
	}
	addAwardReq := &generatedPayments.AddAwardRequest{
		Award: reqAward,
	}

	_, err = h.paymentsClient.AddAward(ctx, addAwardReq)
	if err != nil {
		h.logger.Error("add award failed", zap.Error(err))
		http.Error(w, "Add award failed", http.StatusInternalServerError)
		return
	}

	h.logger.Info("add award success")
	fmt.Fprintf(w, "ok")
}
