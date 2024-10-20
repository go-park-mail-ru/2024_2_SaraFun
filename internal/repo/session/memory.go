package session

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"strconv"
	"time"
)

type InMemoryStorage struct {
	redisClient *redis.Client
}

func New(client *redis.Client) *InMemoryStorage {
	return &InMemoryStorage{redisClient: client}
}

func (repo *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	//repo.sessions[session.SessionID] = session.UserID
	if err := repo.redisClient.Set(ctx, session.SessionID, session.UserID, time.Hour*24).Err(); err != nil {
		return fmt.Errorf("add session failed: %w", err)
	}
	return nil
}

func (repo *InMemoryStorage) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	//if val, ok := repo.sessions[sessionID]; !ok {
	//	return 0, sparkiterrors.ErrInvalidSession
	//} else {
	//	return val, nil
	//}
	val, err := repo.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		return 0, sparkiterrors.ErrInvalidSession
	}
	userId, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("convert session id %s to int: %w", val, err)
	}
	return userId, nil
}

func (repo *InMemoryStorage) CheckSession(ctx context.Context, sessionID string) error {
	//if sessionID == "" {
	//	return sparkiterrors.ErrInvalidSession
	//}
	//if _, ok := repo.sessions[sessionID]; !ok {
	//	return sparkiterrors.ErrInvalidSession
	//} else {
	//	return nil
	//}
	val, err := repo.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		return sparkiterrors.ErrInvalidSession
	}
	if _, err := strconv.Atoi(val); err != nil {
		return fmt.Errorf("convert session id %s to int: %w", val, err)
	}
	return nil
}

func (repo *InMemoryStorage) DeleteSession(ctx context.Context, sessionID string) error {
	//if sessionID == "" {
	//	return sparkiterrors.ErrInvalidSession
	//}
	//delete(repo.sessions, sessionID)
	//return nil
	if err := repo.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return fmt.Errorf("delete session failed: %w", err)
	}
	return nil
}

//func (repo *InMemoryStorage) DeleteSessionByUserID(ctx context.Context, userID int) error {
//	//for i, u := range repo.sessions {
//	//	if u.UserID == userID {
//	//		repo.sessions = append(repo.sessions[:i], repo.sessions[i+1:]...)
//	//		return nil
//	//	}
//	//}
//	//return errors.New("user not found")
//	delete(repo.sessions, session)
//}

//func (repo *InMemoryStorage) GetSessionByUserID(ctx context.Context, userID int) (models.Session, error) {
//	//for _, u := range repo.sessions {
//	//	if u.UserID == userID {
//	//		return u, nil
//	//	}
//	//}
//	//return models.Session{}, errors.New("session not found")
//
//}
