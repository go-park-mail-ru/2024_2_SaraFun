package changepassword

import (
	"encoding/json"
	"fmt"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

//go:generate easyjson -all handler.go

type Request struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ErrResponse struct {
	Message string `json:"message"`
}

//easyjson:skip
type Handler struct {
	authClient          generatedAuth.AuthClient
	personalitiesClient generatedPersonalities.PersonalitiesClient
	logger              *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, personalitiesClient generatedPersonalities.PersonalitiesClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:          authClient,
		personalitiesClient: personalitiesClient,
		logger:              logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req Request
	//err := json.NewDecoder(r.Body).Decode(&req)
	//if err != nil {
	//	h.logger.Error("json decoding error", zap.Error(err))
	//	http.Error(w, "Нам не нравится ваш запрос :(", http.StatusBadRequest)
	//	return
	//}
	err := easyjson.UnmarshalFromReader(r.Body, &req)
	if err != nil {
		h.logger.Error("json decoding error", zap.Error(err))
		http.Error(w, "Нам не нравится ваш запрос :(", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("cookie error", zap.Error(err))
		http.Error(w, "Вы не авторизованы!", http.StatusUnauthorized)
		return
	}
	getIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.authClient.GetUserIDBySessionID(ctx, getIDRequest)

	getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: userId.UserId}
	username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
	if err != nil {
		h.logger.Error("personalitiesClient.GetUsernameByUserID error", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	checkRequest := &generatedPersonalities.CheckPasswordRequest{
		Username: username.Username,
		Password: req.CurrentPassword,
	}

	_, err = h.personalitiesClient.CheckPassword(ctx, checkRequest)
	if err != nil {
		if grpcError, ok := status.FromError(err); ok {
			if grpcError.Code() == codes.InvalidArgument {
				h.logger.Error("bad password", zap.Error(err))
				errResponse := ErrResponse{
					Message: "Неправильный текущий пароль",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusPreconditionFailed)
				err := json.NewEncoder(w).Encode(errResponse)
				if err != nil {
					h.logger.Error("json encoding error", zap.Error(err))
					http.Error(w, "что-то пошло не так :(", http.StatusInternalServerError)
					return
				}
				return
			}
		}
		h.logger.Error("personalitiesClient.CheckPassword error", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	changePasswordRequest := &generatedPersonalities.ChangePasswordRequest{UserID: userId.UserId, Password: req.NewPassword}
	_, err = h.personalitiesClient.ChangePassword(ctx, changePasswordRequest)
	if err != nil {
		h.logger.Error("change password error", zap.Error(err))
		http.Error(w, "Не получилось поменять пароль :(", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Вы успешно поменяли пароль!")
}
