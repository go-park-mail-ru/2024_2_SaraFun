package image

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
)

type imageRepository interface {
	SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error
	//DeleteImage(ctx context.Context, link string) error
}

type UseCase struct {
	imageRepo imageRepository
}

func New(imageRepo imageRepository) *UseCase {
	return &UseCase{
		imageRepo: imageRepo,
	}
}

func (u *UseCase) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error {
	log.Print("before repo save image")
	err := u.imageRepo.SaveImage(ctx, file, fileExt, userId)
	if err != nil {
		return fmt.Errorf("UseCase SaveImage err: %w", err)
	}
	log.Print("after repo save image")
	return nil
}

//func (u *UseCase) DeleteImage(ctx context.Context, link string) error {
//
//}
