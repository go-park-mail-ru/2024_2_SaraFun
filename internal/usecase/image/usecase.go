package image

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
)

type imageRepository interface {
	SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error
	GetImageLinksByUserId(ctx context.Context, id int) ([]string, error)
	DeleteImage(ctx context.Context, id int) error
}

type UseCase struct {
	imageRepo imageRepository
	logger    *zap.Logger
}

func New(imageRepo imageRepository, logger *zap.Logger) *UseCase {
	return &UseCase{
		imageRepo: imageRepo,
		logger:    logger,
	}
}

func (u *UseCase) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error {
	log.Print("before repo save image")
	err := u.imageRepo.SaveImage(ctx, file, fileExt, userId)
	if err != nil {
		u.logger.Error("save image failed", zap.Error(err))
		return fmt.Errorf("UseCase SaveImage err: %w", err)
	}
	log.Print("after repo save image")
	return nil
}

func (u *UseCase) GetImageLinksByUserId(ctx context.Context, id int) ([]string, error) {
	links, err := u.imageRepo.GetImageLinksByUserId(ctx, id)
	if err != nil {
		u.logger.Error("UseCase GetImageLink err", zap.Error(err))
		return nil, fmt.Errorf("UseCase GetImageLink err: %w", err)
	}
	return links, nil
}

func (u *UseCase) DeleteImage(ctx context.Context, id int) error {
	if err := u.imageRepo.DeleteImage(ctx, id); err != nil {
		u.logger.Error("UseCase DeleteImage err", zap.Error(err))
		return fmt.Errorf("UseCase DeleteImage err: %w", err)
	}
	return nil
}
