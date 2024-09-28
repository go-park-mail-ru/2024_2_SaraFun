package session

import (
	"context"
	"sparkit/internal/models"
)

type InMemoryStorage struct {
	sessions map[string]int
}

func New() *InMemoryStorage {
	storage := InMemoryStorage{}
	storage.sessions = make(map[string]int)
	return &storage
}

func (repo *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	repo.sessions[session.SessionID] = session.UserID
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
