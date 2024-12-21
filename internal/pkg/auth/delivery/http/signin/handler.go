package signin

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=signin_mocks . UserService
//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=signin_mocks . SessionService

type UserClient interface {
	CheckPassword(ctx context.Context, in *generatedPersonalities.CheckPasswordRequest) (*generatedPersonalities.CheckPasswordResponse, error)
}

type SessionClient interface {
	CreateSession(ctx context.Context, in *generatedAuth.CreateSessionRequest) (*generatedAuth.CreateSessionResponse, error)
}

type Handler struct {
	userClient    generatedPersonalities.PersonalitiesClient
	sessionClient generatedAuth.AuthClient
	logger        *zap.Logger
}

func NewHandler(userClient generatedPersonalities.PersonalitiesClient, sessionClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{
		userClient:    userClient,
		sessionClient: sessionClient,
		logger:        logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	if r.Method != http.MethodPost {
		h.logger.Error("bad method", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userData := models.User{}

	//if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
	//	h.logger.Error("failed to decode body", zap.Error(err))
	//	http.Error(w, "Неверный формат данных", http.StatusBadRequest)
	//	return
	//}
	err := easyjson.UnmarshalFromReader(r.Body, &userData)
	if err != nil {
		h.logger.Error("failed to parse body", zap.Error(err))
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if user, err := h.userClient.CheckPassword(ctx,
		&generatedPersonalities.CheckPasswordRequest{Username: userData.Username, Password: userData.Password}); err != nil {
		h.logger.Error("invalid credentials", zap.Error(err))
		http.Error(w, "Неверные данные!", http.StatusPreconditionFailed)
		return

	} else {
		req_user := &generatedAuth.User{
			ID:       user.User.ID,
			Username: user.User.Username,
			Email:    user.User.Email,
			Password: user.User.Password,
			Profile:  user.User.Profile,
		}
		if session, err := h.sessionClient.CreateSession(ctx, &generatedAuth.CreateSessionRequest{User: req_user}); err != nil {
			h.logger.Error("failed to create session", zap.Error(err))
			http.Error(w, "Не удалось создать сессию", http.StatusInternalServerError)
			return
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     consts.SessionCookie,
				Value:    session.Session.SessionID,
				Expires:  time.Now().Add(time.Hour * 24),
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})
		}
	}
	h.logger.Info("user created success", zap.String("username", userData.Username))
	fmt.Fprintf(w, "Вы успешно вошли!")
}
