package user

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/utils/hashing"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddUser(ctx context.Context, user models.User) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserList(ctx context.Context) ([]models.User, error)
	GetProfileIdByUserId(ctx context.Context, userId int) (int64, error)
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) RegisterUser(ctx context.Context, user models.User) error {
	err := u.repo.AddUser(ctx, user)
	if err != nil {
		u.logger.Error("bad adduser", zap.Error(err))
		return sparkiterrors.ErrRegistrationUser
	}
	return nil
}

func (u *UseCase) CheckPassword(ctx context.Context, username string, password string) (models.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		u.logger.Error("bad getuserbyusername", zap.Error(err))
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
	if hashing.CheckPasswordHash(password, user.Password) {
		u.logger.Debug("check password successfully", zap.String("username", user.Username), zap.String("password", user.Password))
		return user, nil
	} else {
		u.logger.Info("bad check password", zap.String("username", user.Username), zap.String("password", user.Password))
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
}

func (u *UseCase) GetUserList(ctx context.Context) ([]models.User, error) {
	users, err := u.repo.GetUserList(ctx)
	if err != nil {
		u.logger.Error("bad getuserlist", zap.Error(err))
		return []models.User{}, errors.New("failed to get user list")
	}
	return users, nil
}

func (u *UseCase) GetProfileIdByUserId(ctx context.Context, userId int) (int64, error) {
	profileId, err := u.repo.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		u.logger.Error("failed to get profile id", zap.Int("user_id", userId), zap.Error(err))
		return -1, fmt.Errorf("failed to get profile id by user id: %w", err)
	}
	return profileId, nil
}
