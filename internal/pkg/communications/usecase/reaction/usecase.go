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
	//if firstname == "" {
	//	authors, err = u.repo.GetMatchesByUsername(ctx, userID, username)
	//	if err != nil {
	//		u.logger.Error("UseCase GetMatchesBySearch: failed to GetMatchesByUsername", zap.Error(err))
	//		return nil, fmt.Errorf("failed to GetMatchesByUsername: %w", err)
	//	}
	//} else {
	//	authors, err = u.repo.GetMatchesByFirstName(ctx, userID, firstname)
	//	if err != nil {
	//		u.logger.Error("UseCase GetMatchesBySearch: failed to GetMatchesByFirstName", zap.Error(err))
	//		return nil, fmt.Errorf("failed to GetMatchesByFirstName: %w", err)
	//	}
	//
	//}
	authors, err = u.repo.GetMatchesByString(ctx, userID, search)
	if err != nil {
		u.logger.Error("UseCase GetMatchesBySearch: failed to GetMatchesBySearch", zap.Error(err))
		return nil, fmt.Errorf("failed to GetMatchesBySearch: %w", err)
	}
	u.logger.Info("UseCase GetMatchesBySearch", zap.Int("users", len(authors)))
	return authors, nil
}
