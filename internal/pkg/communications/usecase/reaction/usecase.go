package reaction

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddReaction(ctx context.Context, reaction models.Reaction) error
	GetMatchList(ctx context.Context, userId int) ([]int, error)
	GetReactionList(ctx context.Context, userId int) ([]int, error)
	GetMatchTime(ctx context.Context, firstUser int, secondUser int) (string, error)
	GetMatchesByUsername(ctx context.Context, userID int, username string) ([]int, error)
	GetMatchesByFirstName(ctx context.Context, userID int, firstname string) ([]int, error)
	GetMatchesByString(ctx context.Context, userID int, search string) ([]int, error)
	UpdateOrCreateReaction(ctx context.Context, reaction models.Reaction) error
	CheckMatchExists(ctx context.Context, firstUser int, secondUser int) (bool, error)
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) AddReaction(ctx context.Context, reaction models.Reaction) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	err := u.repo.AddReaction(ctx, reaction)
	if err != nil {
		u.logger.Error("UseCase AddReaction: failed to add reaction", zap.Error(err))
		return fmt.Errorf("failed to AddReaction: %w", err)
	}
	return nil
}

func (u *UseCase) GetMatchList(ctx context.Context, userId int) ([]int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	authors, err := u.repo.GetMatchList(ctx, userId)
	u.logger.Info("matches", zap.Any("authors", authors))
	if err != nil {
		u.logger.Error("UseCase GetMatchList: failed to GetMatchList", zap.Error(err))
		return nil, fmt.Errorf("failed to GetMatchList: %w", err)
	}
	return authors, nil
}

func (u *UseCase) GetReactionList(ctx context.Context, userId int) ([]int, error) {
	receivers, err := u.repo.GetReactionList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to GetReactionList: %w", err)
	}
	return receivers, nil
}

func (u *UseCase) GetMatchTime(ctx context.Context, firstUser int, secondUser int) (string, error) {
	time, err := u.repo.GetMatchTime(ctx, firstUser, secondUser)
	if err != nil {
		u.logger.Error("UseCase GetMatchTime: failed to GetMatchTime", zap.Error(err))
		return "", fmt.Errorf("failed to GetMatchTime: %w", err)
	}
	return time, nil
}

func (u *UseCase) GetMatchesBySearch(ctx context.Context, userID int, search string) ([]int, error) {
	var authors []int
	var err error
	authors, err = u.repo.GetMatchesByString(ctx, userID, search)
	if err != nil {
		u.logger.Error("UseCase GetMatchesBySearch: failed to GetMatchesBySearch", zap.Error(err))
		return nil, fmt.Errorf("failed to GetMatchesBySearch: %w", err)
	}
	u.logger.Info("UseCase GetMatchesBySearch", zap.Int("users", len(authors)))
	return authors, nil
}

func (u *UseCase) UpdateOrCreateReaction(ctx context.Context, reaction models.Reaction) error {
	u.logger.Info("reaction", zap.Any("reaction", reaction))
	err := u.repo.UpdateOrCreateReaction(ctx, reaction)
	if err != nil {
		u.logger.Error("UseCase UpdateOrCreateReaction: failed to UpdateOrCreateReaction", zap.Error(err))
		return fmt.Errorf("failed to UpdateOrCreateReaction: %w", err)
	}
	return nil
}

func (u *UseCase) CheckMatchExists(ctx context.Context, firstUser int, secondUser int) (bool, error) {
	u.logger.Info("check match exists usecase start")
	exists, err := u.repo.CheckMatchExists(ctx, firstUser, secondUser)
	if err != nil {
		u.logger.Error("UseCase CheckMatchExists: failed to CheckMatchExists", zap.Error(err))
		return false, fmt.Errorf("failed to CheckMatchExists: %w", err)
	}
	return exists, nil
}
