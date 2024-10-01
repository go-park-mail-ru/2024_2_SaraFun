package session

import (
	"context"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sync"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	sessions map[string]int
}

func New() *InMemoryStorage {
	return &InMemoryStorage{mu: sync.RWMutex{}, sessions: make(map[string]int)}
}

func (repo *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.sessions[session.SessionID] = session.UserID
	return nil
}

func (repo *InMemoryStorage) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if val, ok := repo.sessions[sessionID]; !ok {
		return 0, sparkiterrors.ErrInvalidSession
	} else {
		return val, nil
	}
}

func (repo *InMemoryStorage) CheckSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return sparkiterrors.ErrInvalidSession
	}
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if _, ok := repo.sessions[sessionID]; !ok {
		return sparkiterrors.ErrInvalidSession
	} else {
		return nil
	}
}

func (repo *InMemoryStorage) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return sparkiterrors.ErrInvalidSession
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	delete(repo.sessions, sessionID)
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
