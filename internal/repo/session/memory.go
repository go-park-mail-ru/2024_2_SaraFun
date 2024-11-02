package session

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"strconv"
	"time"
)

type InMemoryStorage struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

func New(client *redis.Client, logger *zap.Logger) *InMemoryStorage {
	return &InMemoryStorage{redisClient: client, logger: logger}
}

func (repo *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	if err := repo.redisClient.Set(ctx, session.SessionID, session.UserID, time.Hour*24).Err(); err != nil {
		repo.logger.Error("failed to add session", zap.Error(err))
		return fmt.Errorf("add session failed: %w", err)
	}
	repo.logger.Info("added session", zap.String("sessionID", session.SessionID))
	return nil
}

func (repo *InMemoryStorage) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	val, err := repo.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		repo.logger.Error("get sessionId err", zap.Error(err))
		return 0, sparkiterrors.ErrInvalidSession
	}
	userId, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("convert session id %s to int: %w", val, err)
	}
	repo.logger.Info("got session id", zap.String("sessionID", sessionID))
	return userId, nil
}

func (repo *InMemoryStorage) CheckSession(ctx context.Context, sessionID string) error {
	val, err := repo.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		return sparkiterrors.ErrInvalidSession
	}
	if _, err := strconv.Atoi(val); err != nil {
		repo.logger.Error("invalid session id", zap.String("session_id", sessionID))
		return fmt.Errorf("convert session id %s to int: %w", val, err)
	}
	repo.logger.Info("checked session id", zap.String("sessionID", sessionID))
	return nil
}

func (repo *InMemoryStorage) DeleteSession(ctx context.Context, sessionID string) error {
	if err := repo.redisClient.Del(ctx, sessionID).Err(); err != nil {
		repo.logger.Error("delete session failed", zap.Error(err))
		return fmt.Errorf("delete session failed: %w", err)
	}
	repo.logger.Info("deleted session", zap.String("sessionID", sessionID))
	return nil
}
