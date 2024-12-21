package getuserlist

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=getuserlist_mocks . SessionService

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=getuserlist_mocks . ProfileService

type PersonalitiesClient interface {
	GetProfile(ctx context.Context,
		in *generatedPersonalities.GetProfileRequest) (*generatedPersonalities.GetProfileResponse, error)
	GetUsernameByUserId(ctx context.Context,
		in *generatedPersonalities.GetUsernameByUserIDRequest) (*generatedPersonalities.GetUsernameByUserIDResponse, error)
	GetFeedList(ctx context.Context,
		in *generatedPersonalities.GetFeedListRequest) (*generatedPersonalities.GetFeedListResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=getuserlist_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=getuserlist_mocks . ReactionService

type CommunicationsClient interface {
	GetReactionList(ctx context.Context,
		in *generatedCommunications.GetReactionListRequest) (*generatedCommunications.GetReactionListResponse, error)
}

type Response struct {
	Responses []models.PersonCard
}

type Handler struct {
	sessionClient        generatedAuth.AuthClient
	personalitiesClient  generatedPersonalities.PersonalitiesClient
	imageService         ImageService
	communicationsClient generatedCommunications.CommunicationsClient
	logger               *zap.Logger
}

func NewHandler(sessionClient generatedAuth.AuthClient, personalitiesClient generatedPersonalities.PersonalitiesClient, imageService ImageService, communicationsClient generatedCommunications.CommunicationsClient, logger *zap.Logger) *Handler {
	return &Handler{sessionClient: sessionClient, personalitiesClient: personalitiesClient, imageService: imageService, communicationsClient: communicationsClient, logger: logger}
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
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("GetUser Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	getReactionListRequest := &generatedCommunications.GetReactionListRequest{UserId: userId.UserId}
	receivers, err := h.communicationsClient.GetReactionList(ctx, getReactionListRequest)
	if err != nil {
		http.Error(w, "reaction list failed", http.StatusUnauthorized)
		return
	}
	recs := make([]int32, len(receivers.Receivers))
	for i, rc := range receivers.Receivers {
		recs[i] = int32(rc)
	}
	//получить список пользователей
	getFeedRequest := &generatedPersonalities.GetFeedListRequest{
		UserID:    userId.UserId,
		Receivers: recs,
	}
	users, err := h.personalitiesClient.GetFeedList(ctx, getFeedRequest)
	if err != nil {
		h.logger.Error("failed to get feed list", zap.Error(err))
		http.Error(w, "ошибка в получении списка пользователей", http.StatusInternalServerError)
		return
	}

	var cards []models.PersonCard
	for _, user := range users.Users {
		var card models.PersonCard
		getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: user.ID}
		profile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "bad get profile", http.StatusInternalServerError)
			return
		}
		card.Profile = models.Profile{
			ID:        int(profile.Profile.ID),
			FirstName: profile.Profile.FirstName,
			LastName:  profile.Profile.LastName,
			Age:       int(profile.Profile.Age),
			Gender:    profile.Profile.Gender,
			Target:    profile.Profile.Target,
			About:     profile.Profile.About,
		}
		var links []models.Image
		links, err = h.imageService.GetImageLinksByUserId(ctx, int(user.ID))
		if err != nil {
			h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		card.Images = links
		card.UserId = int(user.ID)
		getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: user.ID}
		username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting username", zap.Error(err))
			http.Error(w, "bad get username", http.StatusInternalServerError)
			return
		}
		card.Username = username.Username
		cards = append(cards, card)
	}
	w.Header().Set("Content-Type", "application/json")
	response := Response{Responses: cards}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad marshalling json", zap.Error(err))
		http.Error(w, "bad marshalling json", http.StatusInternalServerError)
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("GetMatches Handler: error writing response", zap.Error(err))
		http.Error(w, "error writing json response", http.StatusInternalServerError)
	}
	h.logger.Info("GetMatches Handler: success")

}
