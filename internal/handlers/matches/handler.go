package matchhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_MatchService.go -package=matchhandler_mocks . MatchService
//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=matchhandler_mocks . SessionService

type MatchService interface {
	GetUserMatches(ctx context.Context, userID int) ([]models.Match, error)
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Handler struct {
	matchService   MatchService
	sessionService SessionService
}

func NewHandler(matchService MatchService, sessionService SessionService) *Handler {
	return &Handler{
		matchService:   matchService,
		sessionService: sessionService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	matches, err := h.matchService.GetUserMatches(ctx, userID)
	if err != nil {
		if errors.Is(err, sparkiterrors.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(matches); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		return 0, err
	}
	sessionID := cookie.Value
	userID, err := h.sessionService.GetUserIDBySessionID(r.Context(), sessionID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
