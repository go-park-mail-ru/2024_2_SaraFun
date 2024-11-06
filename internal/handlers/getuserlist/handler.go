package getuserlist

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=getuserlist_mocks . SessionService
type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=getuserlist_mocks . ProfileService
type ProfileService interface {
	GetProfile(ctx context.Context, id int) (models.Profile, error)
}

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=getuserlist_mocks . UserService
type UserService interface {
	GetUsernameByUserId(ctx context.Context, userId int) (string, error)
	GetFeedList(ctx context.Context, userId int, receivers []int) ([]models.User, error)
}

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=getuserlist_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=getuserlist_mocks . ReactionService
type ReactionService interface {
	GetReactionList(ctx context.Context, userId int) ([]int, error)
}

type Handler struct {
	sessionService  SessionService
	profileService  ProfileService
	userService     UserService
	imageService    ImageService
	reactionService ReactionService
	logger          *zap.Logger
}

func NewHandler(sessionService SessionService, profileService ProfileService, userService UserService, imageService ImageService, reactionService ReactionService, logger *zap.Logger) *Handler {
	return &Handler{sessionService: sessionService, profileService: profileService, userService: userService, imageService: imageService, reactionService: reactionService, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	if r.Method != http.MethodGet {
		h.logger.Error("bad method", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("GetUser Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		h.logger.Error("GetUser Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	receivers, err := h.reactionService.GetReactionList(ctx, userId)
	if err != nil {
		http.Error(w, "reaction list failed", http.StatusUnauthorized)
		return
	}
	//получить список пользователей
	var users []models.User
	users, err = h.userService.GetFeedList(ctx, userId, receivers)
	if err != nil {
		h.logger.Error("failed to get feed list", zap.Error(err))
		http.Error(w, "ошибка в получении списка пользователей", http.StatusInternalServerError)
		return
	}

	var cards []models.PersonCard
	for _, user := range users {
		var card models.PersonCard
		card.Profile, err = h.profileService.GetProfile(ctx, user.ID)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "bad get profile", http.StatusInternalServerError)
			return
		}
		var links []models.Image
		links, err = h.imageService.GetImageLinksByUserId(ctx, user.ID)
		if err != nil {
			h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		card.Images = links
		card.UserId = user.ID
		card.Username, err = h.userService.GetUsernameByUserId(ctx, user.ID)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting username", zap.Error(err))
			http.Error(w, "bad get username", http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(cards)
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

	//перевести в формат json
	//jsonData, err := json.Marshal(users)
	//if err != nil {
	//	h.logger.Error("failed to marshal user list", zap.Error(err))
	//	http.Error(w, "ошибка в сериализации в json", http.StatusInternalServerError)
	//	return
	//}
	//w.Header().Set("Content-Type", "application/json")
	//if _, err := w.Write(jsonData); err != nil {
	//	h.logger.Error("failed to write jsonData", zap.Error(err))
	//	http.Error(w, "не получилось записать json", http.StatusInternalServerError)
	//	return
	//}
}
