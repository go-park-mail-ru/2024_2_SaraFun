package like

import (
	"context"
	"encoding/json"
	"net/http"
)

type LikeService interface {
	AddLike(ctx context.Context, userID int, itemID string) error
}

type Handler struct {
	likeService LikeService
}

func NewHandler(likeService LikeService) *Handler {
	return &Handler{likeService: likeService}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := ctx.Value("userID").(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	itemID := r.URL.Query().Get("itemID")
	if itemID == "" {
		http.Error(w, "itemID is required", http.StatusBadRequest)
		return
	}

	if err := h.likeService.AddLike(ctx, userID, itemID); err != nil {
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "like added"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
