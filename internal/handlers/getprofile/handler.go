package getprofile

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"strconv"
)

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService
type ProfileService interface {
	GetProfile(ctx context.Context, id int) (models.Profile, error)
}

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService
type UserService interface {
	GetProfileIdByUserId(ctx context.Context, userId int) (int, error)
}

type Response struct {
	Profile models.Profile `json:"profile"`
	Images  []models.Image `json:"images"`
}

type Handler struct {
	imageService   ImageService
	profileService ProfileService
	userService    UserService
	logger         *zap.Logger
}

func NewHandler(imageService ImageService, profileService ProfileService, userService UserService, logger *zap.Logger) *Handler {
	return &Handler{imageService: imageService, profileService: profileService, userService: userService, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["userId"])

	profileId, err := h.userService.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		h.logger.Error("getprofileidbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var links []models.Image
	links, err = h.imageService.GetImageLinksByUserId(ctx, userId)
	if err != nil {
		h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile models.Profile
	profile, err = h.profileService.GetProfile(ctx, profileId)
	if err != nil {
		h.logger.Error("getprofileerror", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Profile: profile,
		Images:  links,
	}
	jsonData, err := json.Marshal(response)
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
