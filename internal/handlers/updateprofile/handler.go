package updateprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, id int, profile models.Profile) error
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type UserService interface {
	GetProfileIdByUserId(ctx context.Context, userId int) (int, error)
}

type Handler struct {
	profileService ProfileService
	sessionService SessionService
	userService    UserService
	logger         *zap.Logger
}

func NewHandler(profileService ProfileService, sessionService SessionService, userService UserService, logger *zap.Logger) *Handler {
	return &Handler{profileService: profileService, sessionService: sessionService, userService: userService, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("error getting session cookie", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		h.logger.Error("error getting user id", zap.Error(err))
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	profileId, err := h.userService.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		h.logger.Error("error getting profile id", zap.Error(err))
		http.Error(w, "profile not found", http.StatusUnauthorized)
		return
	}
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		h.logger.Error("error decoding profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.profileService.UpdateProfile(ctx, profileId, profile); err != nil {
		h.logger.Error("error updating profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("profile updated sucessfully")
	fmt.Fprintf(w, "ok")
}
