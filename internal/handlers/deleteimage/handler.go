package deleteimage

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ImageService interface {
	DeleteImage(ctx context.Context, id int) error
}

type Handler struct {
	service ImageService
}

func NewHandler(service ImageService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	imageId, err := strconv.Atoi(mux.Vars(r)["imageId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.service.DeleteImage(ctx, imageId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
