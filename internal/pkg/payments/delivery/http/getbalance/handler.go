package getbalance

import (
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go
type Response struct {
	DailyLikeBalance     int `json:"daily_like_balance"`
	PurchasedLikeBalance int `json:"purchased_like_balance"`
	MoneyBalance         int `json:"money_balance"`
}

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
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "bad cookie", http.StatusUnauthorized)
		return
	}
	getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userID, err := h.authClient.GetUserIDBySessionID(ctx, getUserIDReq)
	if err != nil {
		h.logger.Error("get user id by session id", zap.Error(err))
		http.Error(w, "get user id by session id", http.StatusUnauthorized)
		return
	}
	getBalanceReq := &generatedPayments.GetAllBalanceRequest{UserID: userID.UserId}
	balances, err := h.paymentsClient.GetAllBalance(ctx, getBalanceReq)
	if err != nil {
		h.logger.Error("get all balance", zap.Error(err))
		http.Error(w, "get all balance", http.StatusInternalServerError)
		return
	}
	response := Response{
		DailyLikeBalance:     int(balances.DailyLikeBalance),
		PurchasedLikeBalance: int(balances.PurchasedLikeBalance),
		MoneyBalance:         int(balances.MoneyBalance),
	}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("marshal json", zap.Error(err))
		http.Error(w, "marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write response", zap.Error(err))
		http.Error(w, "write response", http.StatusInternalServerError)
		return
	}
}
