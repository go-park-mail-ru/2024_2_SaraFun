package likehandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/utils/consts"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_LikeService.go -package=likehandler_mocks . LikeService
//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=likehandler_mocks . SessionService

type LikeService interface {
	ProcessLike(ctx context.Context, fromUserID int, toUserID int, isLike bool) error
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Handler struct {
	likeService    LikeService
	sessionService SessionService
}

func NewHandler(likeService LikeService, sessionService SessionService) *Handler {
	return &Handler{
		likeService:    likeService,
		sessionService: sessionService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var likeReq LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&likeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.likeService.ProcessLike(ctx, userID, likeReq.TargetUserID, likeReq.IsLike)
	if err != nil {
		if errors.Is(err, sparkiterrors.ErrUserNotFound) {
			http.Error(w, "Target user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, sparkiterrors.ErrCannotLikeSelf) {
			http.Error(w, "Cannot like/dislike self", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to process like/dislike", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

type LikeRequest struct {
	TargetUserID int  `json:"target_user_id" validate:"required"`
	IsLike       bool `json:"is_like"`
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
