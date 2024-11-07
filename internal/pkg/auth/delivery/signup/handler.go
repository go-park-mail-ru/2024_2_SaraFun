package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"sparkit/internal/utils/hashing"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService
type UserService interface {
	RegisterUser(ctx context.Context, user models.User) (int, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
type SessionService interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService
type ProfileService interface {
	CreateProfile(ctx context.Context, profile models.Profile) (int, error)
}

type Handler struct {
	userService    UserService
	sessionService SessionService
	profileService ProfileService
	logger         *zap.Logger
}

type Request struct {
	User    models.User
	Profile models.Profile
}

func NewHandler(userService UserService, sessionsService SessionService, profileService ProfileService, logger *zap.Logger) *Handler {
	return &Handler{
		userService:    userService,
		sessionService: sessionsService,
		profileService: profileService,
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
	request := Request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.User.Sanitize()
	request.Profile.Sanitize()
	exists, err := h.userService.CheckUsernameExists(ctx, request.User.Username)
	if err != nil {
		h.logger.Error("failed to check username exists", zap.Error(err))
		http.Error(w, "failed to check username exists", http.StatusInternalServerError)
		return
	}
	if exists {
		h.logger.Error("user already exists", zap.String("username", request.User.Username))
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	profileId, err := h.profileService.CreateProfile(ctx, request.Profile)
	if err != nil {
		h.logger.Error("failed to create Profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.User.Profile = profileId
	hashedPass, err := hashing.HashPassword(request.User.Password)
	if err != nil {
		h.logger.Error("failed to hash password", zap.Error(err))
		http.Error(w, "bad password", http.StatusBadRequest)
		return
	}
	request.User.Password = hashedPass
	id, err := h.userService.RegisterUser(ctx, request.User)
	if err != nil {
		h.logger.Error("failed to register User", zap.Error(err))
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	request.User.ID = id

	if session, err := h.sessionService.CreateSession(ctx, request.User); err != nil {
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
	h.logger.Info("good signup", zap.String("username", request.User.Username))
	fmt.Fprintf(w, "ok")
}
