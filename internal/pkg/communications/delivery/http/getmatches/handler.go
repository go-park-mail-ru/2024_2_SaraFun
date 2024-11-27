package getmatches

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

////go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService
//type ReactionService interface {
//	GetMatchList(ctx context.Context, userId int) ([]int, error)
//}

type CommunicationsClient interface {
	GetMatchList(ctx context.Context,
		in *generatedCommunications.GetMatchListRequest) (*generatedCommunications.GetMatchListResponse, error)
}

////go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
//type SessionService interface {
//	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
//}

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

////go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService
//type ProfileService interface {
//	GetProfile(ctx context.Context, id int) (models.Profile, error)
//}
//
////go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService
//type UserService interface {
//	GetUsernameByUserId(ctx context.Context, userId int) (string, error)
//}

type PersonalitiesClient interface {
	GetUsernameByUserId(ctx context.Context,
		in *generatedPersonalities.GetUsernameByUserIDRequest) (*generatedPersonalities.GetUsernameByUserIDResponse, error)
	GetProfile(ctx context.Context,
		in *generatedPersonalities.GetProfileRequest) (*generatedPersonalities.GetProfileResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//type Response struct {
//	Matches []models.PersonCard `json:"matches"`
//}

type Handler struct {
	communicationsClient generatedCommunications.CommunicationsClient
	sessionClient        generatedAuth.AuthClient
	personalitiesClient  generatedPersonalities.PersonalitiesClient
	imageService         ImageService
	logger               *zap.Logger
}

func NewHandler(communicationsClient generatedCommunications.CommunicationsClient, sessionClient generatedAuth.AuthClient,
	personalitiesClient generatedPersonalities.PersonalitiesClient, imageService ImageService, logger *zap.Logger) *Handler {
	return &Handler{
		communicationsClient: communicationsClient,
		sessionClient:        sessionClient,
		personalitiesClient:  personalitiesClient,
		imageService:         imageService,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	getMatchListRequest := &generatedCommunications.GetMatchListRequest{UserID: userId.UserId}
	authors, err := h.communicationsClient.GetMatchList(ctx, getMatchListRequest)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting authors ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	var matches []models.PersonCard
	for _, author := range authors.Authors {
		var matchedUser models.PersonCard
		getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: author}
		profile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
		matchedUser.Profile = models.Profile{
			ID:        int(profile.Profile.ID),
			FirstName: profile.Profile.FirstName,
			LastName:  profile.Profile.LastName,
			Age:       int(profile.Profile.Age),
			Gender:    profile.Profile.Gender,
			Target:    profile.Profile.Target,
			About:     profile.Profile.About,
		}
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "bad get profile", http.StatusInternalServerError)
			return
		}
		var links []models.Image
		links, err = h.imageService.GetImageLinksByUserId(ctx, int(author))
		if err != nil {
			h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchedUser.Images = links
		matchedUser.UserId = int(author)
		getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: author}
		username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
		matchedUser.Username = username.Username
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
