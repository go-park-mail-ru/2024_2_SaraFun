package addProduct

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	authClient     generatedAuth.AuthClient
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, paymentsClient generatedPayments.PaymentClient,
	logger *zap.Logger) *Handler {
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
	var data models.Product
	err = easyjson.UnmarshalFromReader(r.Body, &data)
	if err != nil {
		h.logger.Error("unmarshal data", zap.Error(err))
		http.Error(w, "unmarshal data", http.StatusBadRequest)
		return
	}
	h.logger.Info("count", zap.Any("data", data))

	reqProduct := &generatedPayments.Product{
		Title:       data.Title,
		Description: data.Description,
		ImageLink:   data.ImageLink,
		Price:       int32(data.Price),
		Count:       int32(data.Count),
	}
	createProductReq := &generatedPayments.CreateProductRequest{Product: reqProduct}
	_, err = h.paymentsClient.CreateProduct(ctx, createProductReq)
	if err != nil {
		h.logger.Error("create product error", zap.Error(err))
		http.Error(w, "create product", http.StatusInternalServerError)
		return
	}
	h.logger.Info("add product success")
	fmt.Fprintf(w, "ok")
}
