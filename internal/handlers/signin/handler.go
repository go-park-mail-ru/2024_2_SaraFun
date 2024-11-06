package signin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=signin_mocks . UserService
//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=signin_mocks . SessionService

type UserService interface {
	CheckPassword(ctx context.Context, username string, password string) (models.User, error)
}
type SessionService interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
}

type Handler struct {
	userService    UserService
	sessionService SessionService
	logger         *zap.Logger
}

func NewHandler(userService UserService, sessionService SessionService, logger *zap.Logger) *Handler {
	return &Handler{
		userService:    userService,
		sessionService: sessionService,
		logger:         logger,
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

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		h.logger.Error("failed to decode body", zap.Error(err))
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	if user, err := h.userService.CheckPassword(ctx, userData.Username, userData.Password); err != nil {
		if errors.Is(err, sparkiterrors.ErrWrongCredentials) {
			h.logger.Error("invalid credentials", zap.Error(err))
			http.Error(w, "wrong credentials", http.StatusPreconditionFailed)
			return
		}
	} else {
		if session, err := h.sessionService.CreateSession(ctx, user); err != nil {
			h.logger.Error("failed to create session", zap.Error(err))
			http.Error(w, "Не удалось создать сессию", http.StatusInternalServerError)
			return
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     consts.SessionCookie,
				Value:    session.SessionID,
				Expires:  time.Now().Add(time.Hour * 24),
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Path:     "/",
			})
		}
	}
	h.logger.Info("user created success", zap.String("username", userData.Username))
	fmt.Fprintf(w, "ok")
}
