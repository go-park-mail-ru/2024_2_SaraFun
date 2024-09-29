package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/hashing"
)

type UserService interface {
	RegisterUser(ctx context.Context, user models.User) error
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{
		userService: userService,
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

	fmt.Fprintf(w, "ok")
}
