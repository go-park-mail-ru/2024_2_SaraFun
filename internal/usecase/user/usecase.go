package user

import (
	"context"
	"errors"
	"fmt"
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
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) RegisterUser(ctx context.Context, user models.User) error {
	err := u.repo.AddUser(ctx, user)
	if err != nil {
		return sparkiterrors.ErrRegistrationUser
	}
	return nil
}

func (u *UseCase) CheckPassword(ctx context.Context, username string, password string) (models.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
	if hashing.CheckPasswordHash(password, user.Password) {
		return user, nil
	} else {
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
}

func (u *UseCase) GetUserList(ctx context.Context) ([]models.User, error) {
	users, err := u.repo.GetUserList(ctx)
	if err != nil {
		return []models.User{}, errors.New("failed to get user list")
	}
	return users, nil
}

func (u *UseCase) GetProfileIdByUserId(ctx context.Context, userId int) (int64, error) {
	profileId, err := u.repo.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		return -1, fmt.Errorf("failed to get profile id by user id: %w", err)
	}
	return profileId, nil
}
