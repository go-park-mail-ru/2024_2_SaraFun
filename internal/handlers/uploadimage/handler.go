package uploadimage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sparkit/internal/utils/consts"
)

type ImageService interface {
	SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error
}

type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type Handler struct {
	imageService   ImageService
	sessionService SessionService
}

func NewHandler(imageService ImageService, sessionService SessionService) *Handler {
	return &Handler{imageService: imageService, sessionService: sessionService}
}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "bad ParseMultipartForm", http.StatusInternalServerError)
	}
	log.Print("good parse multipart form")
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "bad image file", http.StatusBadRequest)
	}
	if header == nil {
		http.Error(w, "bad image file", http.StatusBadRequest)
	}
	log.Print("before file.Close()")
	defer file.Close()
	fileExt := filepath.Ext(header.Filename)
	log.Print("good image file")

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	log.Print("good session cookie")
	userId, err := h.sessionService.GetUserIDBySessionID(ctx, cookie.Value)
	if err != nil {
		http.Error(w, "user session err", http.StatusInternalServerError)
		return
	}
	err = h.imageService.SaveImage(ctx, file, fileExt, userId)
	if err != nil {
		http.Error(w, "save image err", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Image saved successfully")
}
