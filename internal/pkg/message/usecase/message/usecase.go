package message

import (
	"context"
	"errors"
	"fmt"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddMessage(ctx context.Context, message *models.Message) (int, error)
	GetLastMessage(ctx context.Context, authorID int, receiverID int) (models.Message, error)
	GetChatMessages(ctx context.Context, firstUserID int, secondUserID int) ([]models.Message, error)
	GetMessagesBySearch(ctx context.Context, userID int, page int, search string) ([]models.Message, error)
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (s *UseCase) AddMessage(ctx context.Context, message *models.Message) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("AddMessage usecase", zap.String("req_id", req_id))
	messageId, err := s.repo.AddMessage(ctx, message)
	if err != nil {
		s.logger.Error("AddMessage error", zap.Error(err))
		return -1, fmt.Errorf("Usecase AddMessage error: %w", err)
	}
	return messageId, nil
}

func (s *UseCase) GetLastMessage(ctx context.Context, authorID int, receiverID int) (models.Message, bool, error) {
	self := true
	msg, err := s.repo.GetLastMessage(ctx, authorID, receiverID)
	if err != nil {
		if err == sparkiterrors.ErrNoResult {
			return models.Message{}, false, nil
		}
		s.logger.Error("GetLastMessage error", zap.Error(err))
		return models.Message{}, false, fmt.Errorf("Usecase GetLastMessage error: %w", err)
	}
	if msg.Author != authorID {
		self = false
	}
	return msg, self, nil
}

func (s *UseCase) GetChatMessages(ctx context.Context, firstUserID int, secondUserID int) ([]models.Message, error) {
	msgs, err := s.repo.GetChatMessages(ctx, firstUserID, secondUserID)
	if err != nil {
		s.logger.Error("GetChatMessages error", zap.Error(err))
		return nil, fmt.Errorf("Usecase GetChatMessages error: %w", err)
	}
	return msgs, nil
}

func (s *UseCase) GetMessagesBySearch(ctx context.Context, userID int, page int, search string) ([]models.Message, error) {
	pageErr := checkPage(page)
	if pageErr != nil {
		s.logger.Error("GetMessagesBySearch error", zap.Error(pageErr))
		return nil, fmt.Errorf("GetMessagesBySearch error: %w", pageErr)
	}
	searchErr := checkSearch(search)
	if searchErr != nil {
		s.logger.Error("GetMessagesBySearch error", zap.Error(searchErr))
		return nil, fmt.Errorf("GetMessagesBySearch error: %w", searchErr)
	}
	userIDErr := checkUserID(userID)
	if userIDErr != nil {
		s.logger.Error("GetMessagesBySearch error", zap.Error(userIDErr))
		return nil, fmt.Errorf("GetMessagesBySearch error: %w", userIDErr)
	}

	msgs, err := s.repo.GetMessagesBySearch(ctx, userID, page, search)
	if err != nil {
		s.logger.Error("GetMessagesBySearch error", zap.Error(err))
		return nil, fmt.Errorf("Usecase GetMessagesBySearch error: %w", err)
	}
	s.logger.Info("msgs", zap.Any("msgs", msgs))
	return msgs, nil
}

func checkPage(page int) error {
	if page < 1 {
		return errors.New("page must be at least 1")
	}
	return nil
}

func checkSearch(search string) error {
	if search == "" {
		return errors.New("search cannot be empty")
	}
	return nil
}

func checkUserID(userID int) error {
	if userID < 1 {
		return errors.New("userId must be at least 1")
	}
	return nil
}
