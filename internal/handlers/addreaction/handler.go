package addreaction

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService
type ReactionService interface {
	AddReaction(ctx context.Context, reaction models.Reaction) error
}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

//type ProfileService interface {
//	GetProfile(ctx context.Context, id int) (models.Profile, error)
//}

//type UserService interface {
//
//}

type Handler struct {
	reactionService ReactionService
	sessionService  SessionService
	logger          *zap.Logger
}

func NewHandler(reactionService ReactionService, sessionService SessionService, logger *zap.Logger) *Handler {
	return &Handler{reactionService: reactionService, sessionService: sessionService, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	var reaction models.Reaction
	err = json.NewDecoder(r.Body).Decode(&reaction)
	reaction.Author = userId
	if err != nil {
		h.logger.Error("AddReaction Handler: bad decoding ", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = h.reactionService.AddReaction(ctx, reaction)
	if err != nil {
		h.logger.Error("AddReaction Handler: error adding reaction", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("AddReaction Handler: added reaction", zap.Any("reaction", reaction))
	fmt.Fprintf(w, "ok")
}
