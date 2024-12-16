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
	h.logger.Info("handle request", zap.Any("jsonData", jsonData))
	object := jsonData["object"].(map[string]interface{})
	amount := object["amount"].(map[string]interface{})
	h.logger.Info("amount", zap.Any("amount", amount))
	price, err := strconv.ParseFloat(amount["value"].(string), 32)
	if err != nil {
		h.logger.Error("parse json price", zap.Error(err))
		http.Error(w, "parse json error", http.StatusBadRequest)
		return
	}
	h.logger.Info("price", zap.Any("price", price))
	payerID, err := strconv.Atoi(object["description"].(string))
	if err != nil {
		h.logger.Error("parse json payer id", zap.Error(err))
		http.Error(w, "parse json error", http.StatusBadRequest)
		return
	}
	h.logger.Info("payer", zap.Any("payerID", payerID))
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
