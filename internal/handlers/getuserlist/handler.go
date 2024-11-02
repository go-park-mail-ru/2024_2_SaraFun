package getuserlist

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
)

type UserUsecase interface {
	GetUserList(ctx context.Context) ([]models.User, error)
}

type Handler struct {
	usecase UserUsecase
	logger  *zap.Logger
}

func NewHandler(usecase UserUsecase, logger *zap.Logger) *Handler {
	return &Handler{usecase: usecase, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Error("bad method", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//получить список пользователей
	var users []models.User
	users, err := h.usecase.GetUserList(ctx)
	if err != nil {
		h.logger.Error("failed to get user list", zap.Error(err))
		http.Error(w, "ошибка в получении списка пользователей", http.StatusInternalServerError)
		return
	}
	//перевести в формат json
	jsonData, err := json.Marshal(users)
	if err != nil {
		h.logger.Error("failed to marshal user list", zap.Error(err))
		http.Error(w, "ошибка в сериализации в json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		h.logger.Error("failed to write jsonData", zap.Error(err))
		http.Error(w, "не получилось записать json", http.StatusInternalServerError)
		return
	}
}
