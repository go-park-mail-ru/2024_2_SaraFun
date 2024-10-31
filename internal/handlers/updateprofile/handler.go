package updateprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, id int64, profile models.Profile) error
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type UserService interface {
	GetProfileIdByUserId(ctx context.Context, userId int) (int64, error)
}

type Handler struct {
	profileService ProfileService
	sessionService SessionService
	userService    UserService
}

func NewHandler(profileService ProfileService, sessionService SessionService, userService UserService) *Handler {
	return &Handler{profileService: profileService, sessionService: sessionService, userService: userService}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	profileId, err := h.userService.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		http.Error(w, "profile not found", http.StatusUnauthorized)
		return
	}
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.profileService.UpdateProfile(ctx, profileId, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok")
}
