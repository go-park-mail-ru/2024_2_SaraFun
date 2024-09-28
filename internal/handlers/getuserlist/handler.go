package getuserlist

import (
	"context"
	"encoding/json"
	"net/http"
	"sparkit/internal/models"
)

type UserUsecase interface {
	GetUserList(ctx context.Context) ([]models.User, error)
}

type Handler struct {
	usecase UserUsecase
}

func NewHandler(usecase UserUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//получить список пользователей
	var users []models.User
	users, err := h.usecase.GetUserList(ctx)
	if err != nil {
		http.Error(w, "ошибка в получении списка пользователей", http.StatusInternalServerError)
	}
	//перевести в формат json
	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "ошибка в сериализации в json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, "не получилось записать json", http.StatusInternalServerError)
	}
}
