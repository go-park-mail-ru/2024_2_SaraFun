package signin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"strconv"
	"time"
)

type UserService interface {
	CheckPassword(ctx context.Context, username string, password string) (models.User, error)
}
type SessionService interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
}

type Handler struct {
	userService    UserService
	sessionService SessionService
}

func NewHandler(userService UserService, sessionService SessionService) *Handler {
	return &Handler{
		userService:    userService,
		sessionService: sessionService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userData := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	if user, err := h.userService.CheckPassword(ctx, userData.Username, userData.Password); err != nil {
		if errors.Is(err, sparkiterrors.ErrWrongCredentials) {
			http.Error(w, "wrong credentials", http.StatusPreconditionFailed)
		}
	} else {
		if session, err := h.sessionService.CreateSession(ctx, user); err != nil {
			http.Error(w, "Не удалось создать сессию", http.StatusInternalServerError)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_id",
				Value:   strconv.Itoa(session.ID),
				Expires: time.Now().Add(time.Hour * 24),
			})
		}
	}

}
