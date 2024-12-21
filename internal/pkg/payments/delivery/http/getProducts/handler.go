package getProducts

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
	Responses []models.Product `json:"responses"`
}

type Handler struct {
	authClient     generatedAuth.AuthClient
	paymentsClient generatedPayments.PaymentClient
	logger         *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, paymentClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:     authClient,
		paymentsClient: paymentClient,
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
	getProductsReq := &generatedPayments.GetProductsRequest{}
	products, err := h.paymentsClient.GetProducts(ctx, getProductsReq)
	if err != nil {
		h.logger.Error("get products", zap.Error(err))
		http.Error(w, "get products error", http.StatusInternalServerError)
		return
	}

	var prs []models.Product
	for _, product := range products.Products {
		pr := models.Product{
			Title:       product.Title,
			Description: product.Description,
			ImageLink:   product.ImageLink,
			Price:       int(product.Price),
			Count:       int(product.Count),
		}
		prs = append(prs, pr)
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{Responses: prs}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("GetProducts Handler: bad marshalling json", zap.Error(err))
		http.Error(w, "bad marshalling json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("GetProducts Handler: error writing response", zap.Error(err))
		http.Error(w, "error writing json response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("GetProducts Handler: success")
}
