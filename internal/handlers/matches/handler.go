package match

import (
	"context"
	"encoding/json"
	"net/http"
	"sparkit/internal/models"
)

type MatchService interface {
	GetMatches(ctx context.Context, userID int) ([]models.Match, error)
}

type Handler struct {
	matchService MatchService
}

func NewHandler(matchService MatchService) *Handler {
	return &Handler{matchService: matchService}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := ctx.Value("userID").(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	matches, err := h.matchService.GetMatches(ctx, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}
