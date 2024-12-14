package getprofile

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService

type PersonalitiesClient interface {
	GetProfile(ctx context.Context,
		in *generatedPersonalities.GetProfileRequest) (*generatedPersonalities.GetProfileResponse, error)
	GetProfileIdByUserId(ctx context.Context,
		in *generatedPersonalities.GetProfileIDByUserIDRequest) (*generatedPersonalities.GetProfileIDByUserIDResponse, error)
	GetUserIdByUsername(ctx context.Context,
		in *generatedPersonalities.GetUserIDByUsernameRequest) (*generatedPersonalities.GetUserIDByUsernameResponse, error)
}

type Response struct {
	Profile models.Profile `json:"profile"`
	Images  []models.Image `json:"images"`
}

//easyjson:skip
type Handler struct {
	imageService        ImageService
	personalitiesClient generatedPersonalities.PersonalitiesClient
	logger              *zap.Logger
}

func NewHandler(imageService ImageService, personalitiesClient generatedPersonalities.PersonalitiesClient, logger *zap.Logger) *Handler {
	return &Handler{imageService: imageService, personalitiesClient: personalitiesClient, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	vars := mux.Vars(r)
	//userId, _ := strconv.Atoi(vars["userId"])
	username := vars["username"]
	getUserIdRequest := &generatedPersonalities.GetUserIDByUsernameRequest{Username: username}
	userId, err := h.personalitiesClient.GetUserIDByUsername(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("Error getting user id by username", zap.String("username", username), zap.Error(err))
		http.Error(w, "don`t get user by username", http.StatusInternalServerError)
		return
	}
	getProfileIdRequest := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: userId.UserID}
	profileId, err := h.personalitiesClient.GetProfileIDByUserID(ctx, getProfileIdRequest)
	if err != nil {
		h.logger.Error("getprofileidbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var links []models.Image
	links, err = h.imageService.GetImageLinksByUserId(ctx, int(userId.UserID))
	if err != nil {
		h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: profileId.ProfileID}
	profile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
	if err != nil {
		h.logger.Error("getprofileerror", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profileResponse := models.Profile{
		ID:           int(profile.Profile.ID),
		FirstName:    profile.Profile.FirstName,
		LastName:     profile.Profile.LastName,
		Age:          int(profile.Profile.Age),
		Gender:       profile.Profile.Gender,
		Target:       profile.Profile.Target,
		About:        profile.Profile.About,
		BirthdayDate: profile.Profile.BirthDate,
	}
	response := Response{
		Profile: profileResponse,
		Images:  links,
	}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("json marshal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write jsonData error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("getprofile success")
}
