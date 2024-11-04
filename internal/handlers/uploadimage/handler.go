package uploadimage

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sparkit/internal/utils/consts"
)

type ImageService interface {
	SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int64, error)
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Response struct {
	ImageId int64
}

type Handler struct {
	imageService   ImageService
	sessionService SessionService
	logger         *zap.Logger
}

func NewHandler(imageService ImageService, sessionService SessionService, logger *zap.Logger) *Handler {
	return &Handler{imageService: imageService, sessionService: sessionService, logger: logger}
}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error("parse multipart form", zap.Error(err))
		http.Error(w, "bad ParseMultipartForm", http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("failed to parse multipart form", zap.Error(err))
		http.Error(w, "bad image file", http.StatusBadRequest)
		return
	}
	if header == nil {
		h.logger.Error("failed to parse multipart form")
		http.Error(w, "bad image file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileExt := filepath.Ext(header.Filename)
	log.Print("good image file")

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("failed to get session cookie", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	log.Print("good session cookie")
	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		h.logger.Error("failed to get user id", zap.Error(err))
		http.Error(w, "user session err", http.StatusInternalServerError)
		return
	}
	id, err := h.imageService.SaveImage(ctx, file, fileExt, userId)
	if err != nil {
		h.logger.Error("failed to save image", zap.Error(err))
		http.Error(w, "save image err", http.StatusInternalServerError)
		return
	}

	response := Response{ImageId: id}
	jsonData, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("failed to marshal image", zap.Error(err))
		http.Error(w, "save image err", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("failed to write response", zap.Error(err))
		http.Error(w, "save image err", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.logger.Info("image saved successfully")
}
