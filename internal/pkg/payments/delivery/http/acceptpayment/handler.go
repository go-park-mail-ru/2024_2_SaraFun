package acceptpayment

import (
	"encoding/json"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	authClient     generatedAuth.AuthClient
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, paymentsClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:     authClient,
		paymentsClient: paymentsClient,
		logger:         logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var jsonData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		h.logger.Error("decode json", zap.Error(err))
		http.Error(w, "decode json error", http.StatusBadRequest)
		return
	}
	amount := jsonData["amount"].(map[string]interface{})

	price, err := strconv.Atoi(amount["value"].(string))
	if err != nil {
		h.logger.Error("parse json price", zap.Error(err))
		http.Error(w, "parse json error", http.StatusBadRequest)
		return
	}
	payerID, err := strconv.Atoi(jsonData["description"].(string))
	if err != nil {
		h.logger.Error("parse json payer id", zap.Error(err))
		http.Error(w, "parse json error", http.StatusBadRequest)
		return
	}

	changeBalanceReq := generatedPayments.ChangeBalanceRequest{
		UserID: int32(payerID),
		Amount: int32(price),
	}
	_, err = h.paymentsClient.ChangeBalance(ctx, &changeBalanceReq)
	if err != nil {
		h.logger.Error("change balance", zap.Error(err))
		http.Error(w, "change balance error", http.StatusUnauthorized)
		return
	}
}