package profile

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sparkit/internal/models"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile models.Profile) (int64, error)
	UpdateProfile(ctx context.Context, id int64, profile models.Profile) error
	GetProfile(ctx context.Context, id int64) (models.Profile, error)
	DeleteProfile(ctx context.Context, id int) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) CreateProfile(ctx context.Context, profile models.Profile) (int64, error) {
	res, err := u.repo.CreateProfile(ctx, profile)
	if err != nil {
		u.logger.Error("create profile err", zap.Error(err))
		return 0, fmt.Errorf("create profile err: %w", err)
	}
	return res, nil
}

func (u *UseCase) UpdateProfile(ctx context.Context, id int64, profile models.Profile) error {
	if err := u.repo.UpdateProfile(ctx, id, profile); err != nil {
		u.logger.Error("update profile err", zap.Error(err))
		return fmt.Errorf("update profile err: %w", err)
	}
	return nil
}

func (u *UseCase) GetProfile(ctx context.Context, id int64) (models.Profile, error) {
	res, err := u.repo.GetProfile(ctx, id)
	if err != nil {
		u.logger.Error("get profile err", zap.Error(err))
		return models.Profile{}, fmt.Errorf("get profile err: %w", err)
	}
	return res, nil
}

func (u *UseCase) DeleteProfile(ctx context.Context, id int) error {
	if err := u.repo.DeleteProfile(ctx, id); err != nil {
		u.logger.Error("delete profile err", zap.Error(err))
		return fmt.Errorf("delete profile err: %w", err)
	}
	return nil
}
