package getmatches

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

type ReactionService interface {
	GetMatchList(ctx context.Context, userId int) ([]int, error)
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type ProfileService interface {
	GetProfile(ctx context.Context, id int64) (models.Profile, error)
}

type UserService interface {
	GetUsernameByUserId(ctx context.Context, userId int) (string, error)
}

//type Response struct {
//	Matches []models.MatchedUser `json:"matches"`
//}

type Handler struct {
	reactionService ReactionService
	sessionService  SessionService
	profileService  ProfileService
	userService     UserService
	logger          *zap.Logger
}

func NewHandler(reactionService ReactionService, sessionService SessionService, profileService ProfileService, userService UserService, logger *zap.Logger) *Handler {
	return &Handler{reactionService: reactionService, sessionService: sessionService, profileService: profileService, userService: userService, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	var authors []int
	authors, err = h.reactionService.GetMatchList(ctx, userId)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting authors ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	var matches []models.MatchedUser
	for _, author := range authors {
		var matchedUser models.MatchedUser
		matchedUser.Profile, err = h.profileService.GetProfile(ctx, int64(author))
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "session not found", http.StatusUnauthorized)
			return
		}
		matchedUser.UserId = author
		matchedUser.Username, err = h.userService.GetUsernameByUserId(ctx, author)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting username", zap.Error(err))
			http.Error(w, "session not found", http.StatusUnauthorized)
		}
		matches = append(matches, matchedUser)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(matches)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad marshalling json", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("GetMatches Handler: error writing response", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
	}
	h.logger.Info("GetMatches Handler: success")

}
