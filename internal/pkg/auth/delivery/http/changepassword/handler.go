package changepassword

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

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

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error("json decoding error", zap.Error(err))
		http.Error(w, "bad parse", http.StatusBadRequest)
	}
	user.Sanitize()
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("cookie error", zap.Error(err))
		http.Error(w, "bad cookie", http.StatusUnauthorized)
	}
	getIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.authClient.GetUserIDBySessionID(ctx, getIDRequest)

	changePasswordRequest := &generatedPersonalities.ChangePasswordRequest{UserID: userId.UserId, Password: user.Password}
	_, err = h.personalitiesClient.ChangePassword(ctx, changePasswordRequest)
	if err != nil {
		h.logger.Error("change password error", zap.Error(err))
		http.Error(w, "bad change password", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "ok")
}
