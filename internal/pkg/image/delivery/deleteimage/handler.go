package deleteimage

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_image_service.go -package=mocks sparkit/internal/handlers/deleteimage ImageService

type ImageService interface {
	DeleteImage(ctx context.Context, id int) error
}

type Handler struct {
	service ImageService
	logger  *zap.Logger
}

func NewHandler(service ImageService, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	imageId, err := strconv.Atoi(mux.Vars(r)["imageId"])
	if err != nil {
		h.logger.Error("error getting image id", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.service.DeleteImage(ctx, imageId); err != nil {
		h.logger.Error("error deleting image", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("image delete good operation")
}
