package user

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"sparkit/internal/utils/hashing"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddUser(ctx context.Context, user models.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserList(ctx context.Context, userId int) ([]models.User, error)
	GetProfileIdByUserId(ctx context.Context, userId int) (int, error)
	GetUsernameByUserId(ctx context.Context, userId int) (string, error)
	GetFeedList(ctx context.Context, userId int, receivers []int) ([]models.User, error)
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) GetFeedList(ctx context.Context, userId int, receivers []int) ([]models.User, error) {
	users, err := u.repo.GetFeedList(ctx, userId, receivers)
	if err != nil {
		u.logger.Error("bad getuserlist", zap.Error(err))
		return nil, fmt.Errorf("failed to get user list: %w", err)
	}
	return users, nil
}

func (u *UseCase) RegisterUser(ctx context.Context, user models.User) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	id, err := u.repo.AddUser(ctx, user)
	if err != nil {
		u.logger.Error("bad adduser", zap.Error(err))
		return -1, sparkiterrors.ErrRegistrationUser
	}
	return id, nil
}

func (u *UseCase) CheckPassword(ctx context.Context, username string, password string) (models.User, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
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

func (u *UseCase) GetUserList(ctx context.Context, userId int) ([]models.User, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	users, err := u.repo.GetUserList(ctx, userId)
	if err != nil {
		u.logger.Error("bad getuserlist", zap.Error(err))
		return []models.User{}, errors.New("failed to get user list")
	}
	return users, nil
}

func (u *UseCase) GetProfileIdByUserId(ctx context.Context, userId int) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	profileId, err := u.repo.GetProfileIdByUserId(ctx, userId)
	if err != nil {
		u.logger.Error("failed to get profile id", zap.Int("user_id", userId), zap.Error(err))
		return -1, fmt.Errorf("failed to get profile id by user id: %w", err)
	}
	return profileId, nil
}

func (u *UseCase) GetUsernameByUserId(ctx context.Context, userId int) (string, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	username, err := u.repo.GetUsernameByUserId(ctx, userId)
	if err != nil {
		u.logger.Error("failed to get username", zap.Int("user_id", userId), zap.Error(err))
		return "", sparkiterrors.ErrWrongCredentials
	}
	return username, nil
}
