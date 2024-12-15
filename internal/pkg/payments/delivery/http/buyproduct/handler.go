package buyproduct

import (
	"fmt"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"net/http"
)

//go:generate easyjson -all handler.go

type Request struct {
	Title string `json:"title"`
	Price int    `json:"price"`
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
	var data Request
	err = easyjson.UnmarshalFromReader(r.Body, &data)
	if err != nil {
		h.logger.Error("json unmarshal", zap.Error(err))
		http.Error(w, "json unmarshal error", http.StatusBadRequest)
		return
	}

	buyLikesReq := &generatedPayments.BuyLikesRequest{
		Title:  data.Title,
		Amount: int32(data.Price),
		UserID: userID.UserId,
	}
	_, err = h.paymentsClient.BuyLikes(ctx, buyLikesReq)
	if err != nil {
		st, ok := status.FromError(err)
		h.logger.Info("status code", zap.String("code", st.String()))
		if ok && st.Message() == "Недостаточно средств" {
			h.logger.Error("buy likes failed", zap.Error(err))
			http.Error(w, "У вас недостаточно средств. Срочно пополните его!", http.StatusBadRequest)
			return
		} else if ok && st.Message() == "Суммы не хватает даже на один лайк" {
			h.logger.Error("buy likes failed", zap.Error(err))
			http.Error(w, "Суммы не хватает даже на один лайк! Потратьте больше денег!", http.StatusBadRequest)
			return
		}
		h.logger.Error("buy likes", zap.Error(err))
		http.Error(w, "buy likes", http.StatusInternalServerError)
		return
	}
	h.logger.Info("buy product success")
	fmt.Fprintf(w, "ok")
}
