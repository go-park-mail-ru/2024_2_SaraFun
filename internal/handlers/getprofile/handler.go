package getprofile

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"sparkit/internal/models"
	"strconv"
)

type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]string, error)
}

type ProfileService interface {
	GetProfile(ctx context.Context, id int64) (models.Profile, error)
}

type UserService interface {
	GetProfileIdByUserId(ctx context.Context, userId int) (int64, error)
}

type Response struct {
	Profile   models.Profile `json:"profile"`
	ImageURLs []string       `json:"imageURLs"`
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
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["userId"])

	profileId, err := h.userService.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		h.logger.Error("getprofileidbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var links []string
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
		Profile:   profile,
		ImageURLs: links,
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
