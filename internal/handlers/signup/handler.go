package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"sparkit/internal/utils/hashing"
	"time"
)

type UserService interface {
	RegisterUser(ctx context.Context, user models.User) error
}

type SessionService interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
}

type Handler struct {
	userService    UserService
	sessionService SessionService
}

func NewHandler(userService UserService, sessionsService SessionService) *Handler {
	return &Handler{
		userService:    userService,
		sessionService: sessionsService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := models.User{}
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	hashedPass, err := hashing.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "bad password", http.StatusBadRequest)
		return
	}
	user.Password = hashedPass
	if err := h.userService.RegisterUser(ctx, user); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if session, err := h.sessionService.CreateSession(ctx, user); err != nil {
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
	fmt.Fprintf(w, "ok")
}
