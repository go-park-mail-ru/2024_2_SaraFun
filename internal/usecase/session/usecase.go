package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"time"
)

type Repository interface {
	AddSession(ctx context.Context, session models.Session) error
	//DeleteSessionByUserID(ctx context.Context, userID int) error
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (s *UseCase) CreateSession(ctx context.Context, user models.User) (models.Session, error) {
	session := models.Session{
		SessionID: uuid.New().String(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	err := s.repo.AddSession(ctx, session)
	if err != nil {
		return models.Session{}, errors.New("failed to create session")
	}
	return session, nil
}

func (s *UseCase) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	userID, err := s.repo.GetUserIDBySessionID(ctx, sessionID)
	if err != nil {
		return 0, sparkiterrors.ErrInvalidSession
	}
	return userID, nil
}

//func (s *UseCase) DeleteSessionByUserID(ctx context.Context, userID int) error {
//	err := s.repo.DeleteSessionByUserID(ctx, userID)
//	if err != nil {
//		return errors.New("failed to delete session")
//	}
//	return nil
//}

//func (s *UseCase) GetSessionByUserID(ctx context.Context, userID int) (models.Session, error) {
//	session, err := s.repo.GetSessionByUserID(ctx, userID)
//	if err != nil {
//		return models.Session{}, errors.New("failed to get session")
//	}
//	return session, nil
//}
