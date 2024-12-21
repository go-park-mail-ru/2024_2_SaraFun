package topUpBalance

import (
	"bytes"
	"encoding/json"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	authClient generatedAuth.AuthClient
	logger     *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient: authClient,
		logger:     logger,
	}
}

//go:generate easyjson -all handler.go

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnUrl string `json:"return_url"`
}

type Request struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

type Response struct {
	RedirectLink string `json:"redirect_link"`
}

type APIRequest struct {
	Amount       Amount       `json:"amount"`
	Capture      string       `json:"capture"`
	Confirmation Confirmation `json:"confirmation"`
	Description  string       `json:"description"`
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

	var requestData Request
	if err := easyjson.UnmarshalFromReader(r.Body, &requestData); err != nil {
		h.logger.Error("unmarshal request", zap.Error(err))
		http.Error(w, "unmarshal request", http.StatusBadRequest)
		return
	}

	url := "https://api.yookassa.ru/v3/payments"
	returnUrl := "https://spark-it.site/shop"

	apiRequest := &APIRequest{
		Amount: Amount{
			Value:    requestData.Price,
			Currency: "RUB",
		},
		Capture: "true",
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnUrl: returnUrl,
		},
		Description: strconv.Itoa(int(userID.UserId)),
	}
	apiData, err := easyjson.Marshal(apiRequest)
	if err != nil {
		h.logger.Error("marshal api request", zap.Error(err))
		http.Error(w, "marshal api request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(apiData))
	if err != nil {
		h.logger.Error("bad create api request", zap.Error(err))
		http.Error(w, "bad create api request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.New().String())

	shopID := os.Getenv("SHOP_ID")
	secretKey := os.Getenv("SECRET_SHOP_KEY")
	h.logger.Info("create api request", zap.String("shop_id", shopID))
	h.logger.Info("create api request", zap.String("secret_key", secretKey))
	req.SetBasicAuth(shopID, secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("bad api request", zap.Error(err))
		http.Error(w, "bad api request", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("bad api response reading body", zap.Error(err))
		http.Error(w, "bad api response reading body", http.StatusInternalServerError)
		return
	}
	var apiResponse map[string]interface{}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		h.logger.Error("unmarshal api response body", zap.Error(err))
		http.Error(w, "unmarshal api response body", http.StatusInternalServerError)
		return
	}

	h.logger.Info("api response", zap.Any("apiResponse", apiResponse))
	confirmation := apiResponse["confirmation"].(map[string]interface{})
	response := Response{RedirectLink: confirmation["confirmation_url"].(string)}
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
