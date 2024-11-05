package getmatches

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService
type ReactionService interface {
	GetMatchList(ctx context.Context, userId int) ([]int, error)
}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService
type ProfileService interface {
	GetProfile(ctx context.Context, id int) (models.Profile, error)
}

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService
type UserService interface {
	GetUsernameByUserId(ctx context.Context, userId int) (string, error)
}

type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//type Response struct {
//	Matches []models.PersonCard `json:"matches"`
//}

type Handler struct {
	reactionService ReactionService
	sessionService  SessionService
	profileService  ProfileService
	userService     UserService
	imageService    ImageService
	logger          *zap.Logger
}

func NewHandler(reactionService ReactionService, sessionService SessionService, profileService ProfileService, userService UserService, imageService ImageService, logger *zap.Logger) *Handler {
	return &Handler{reactionService: reactionService, sessionService: sessionService, profileService: profileService, userService: userService, imageService: imageService, logger: logger}
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

	var matches []models.PersonCard
	for _, author := range authors {
		var matchedUser models.PersonCard
		matchedUser.Profile, err = h.profileService.GetProfile(ctx, author)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "bad get profile", http.StatusInternalServerError)
			return
		}
		var links []models.Image
		links, err = h.imageService.GetImageLinksByUserId(ctx, userId)
		if err != nil {
			h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchedUser.Images = links
		matchedUser.UserId = author
		matchedUser.Username, err = h.userService.GetUsernameByUserId(ctx, author)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting username", zap.Error(err))
			http.Error(w, "bad get username", http.StatusInternalServerError)
			return
		}
		matches = append(matches, matchedUser)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(matches)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad marshalling json", zap.Error(err))
		http.Error(w, "bad marshalling json", http.StatusInternalServerError)
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("GetMatches Handler: error writing response", zap.Error(err))
		http.Error(w, "error writing json response", http.StatusUnauthorized)
	}
	h.logger.Info("GetMatches Handler: success")

}
