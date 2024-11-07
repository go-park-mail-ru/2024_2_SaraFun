package usecase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int, error)
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
	DeleteImage(ctx context.Context, id int) error
}

type UseCase struct {
	imageRepo Repository
	logger    *zap.Logger
}

func New(imageRepo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{
		imageRepo: imageRepo,
		logger:    logger,
	}
}

func (u *UseCase) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	log.Print("before repo save image")
	id, err := u.imageRepo.SaveImage(ctx, file, fileExt, userId)
	if err != nil {
		u.logger.Error("save image failed", zap.Error(err))
		return -1, fmt.Errorf("UseCase SaveImage err: %w", err)
	}
	log.Print("after repo save image")
	return id, nil
}

func (u *UseCase) GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	links, err := u.imageRepo.GetImageLinksByUserId(ctx, id)
	if err != nil {
		u.logger.Error("UseCase GetImageLink err", zap.Error(err))
		return nil, fmt.Errorf("UseCase GetImageLink err: %w", err)
	}
	return links, nil
}

func (u *UseCase) DeleteImage(ctx context.Context, id int) error {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	if err := u.imageRepo.DeleteImage(ctx, id); err != nil {
		u.logger.Error("UseCase DeleteImage err", zap.Error(err))
		return fmt.Errorf("UseCase DeleteImage err: %w", err)
	}
	return nil
}
