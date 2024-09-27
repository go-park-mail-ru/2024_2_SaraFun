package pkg

import (
	"context"
	"fmt"
	"internal/models"
	"time"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	GetUsers(ctx context.Context, count int64, offset int64) ([]models.User, error)
}

type inMemoryUserRepo struct {
	users []models.User
}

func (repo *inMemoryUserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	user.ID = fmt.Sprintf("%d", len(repo.users)+1) // генерация ID
	user.CreatedAt = time.Now()
	repo.users = append(repo.users, user)
	return user, nil
}

func (repo *inMemoryUserRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	for _, user := range repo.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("user not found")
}

func (repo *inMemoryUserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	for _, user := range repo.users {
		if user.email == email {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("user not found")
}

func (repo *inMemoryUserRepo) GetUsers(ctx context.Context, count int64, offset int64) ([]models.User, error) {
	start := offset
	end := offset + count
	if start > int64(len(repo.users)) {
		return []models.User{}, nil
	}
	if end > int64(len(repo.users)) {
		end = int64(len(repo.users))
	}
	return repo.users[start:end], nil
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	ListUsers(ctx context.Context, count int64, offset int64) ([]models.User, error)
}

type userUsecase struct {
	userRepo UserRepo
}

func NewUserUsecase(repo UserRepo) UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) RegisterUser(ctx context.Context, user models.User) (models.User, error) {
	return u.userRepo.CreateUser(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id string) (models.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context, count int64, offset int64) ([]models.User, error) {
	return u.userRepo.GetUsers(ctx, count, offset)
}

//-----------------------------------------------------
//работа с сессией

type SessionRepo interface {
	CreateSession(ctx context.Context, session models.Session) (models.Session, error)
	GetSessionByID(ctx context.Context, id string) (models.Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}
   
type inMemorySessionRepo struct {
s	essions []models.Session
}

func (repo *inMemorySessionRepo) CreateSession(ctx context.Context, session models.Session) (models.Session, error) {
	session.ID = fmt.Sprintf("%d", len(repo.sessions)+1) // генерация ID
	session.CreatedAt = time.Now()
	repo.sessions = append(repo.sessions, session)
	return session, nil
}

func (repo *inMemorySessionRepo) GetSessionByID(ctx context.Context, id string) (models.Session, error) {
	for _, session := range repo.sessions {
		if session.ID == id {
		return session, nil
		}
	}
	return models.Session{}, fmt.Errorf("session not found")
}

func (repo *inMemorySessionRepo) GetSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error) {
	var userSessions []models.Session
	for _, session := range repo.sessions {
		if session.UserID == userID {
		userSessions = append(userSessions, session)
		}
	}
	return userSessions, nil
}

func (repo *inMemorySessionRepo) DeleteSession(ctx context.Context, id string) error {
	for i, session := range repo.sessions {
		if session.ID == id {
		repo.sessions = append(repo.sessions[:i], repo.sessions[i+1:]...)
		return nil
		}
	}
	return fmt.Errorf("session not found")
}

type SessionUsecase interface {
	RegisterSession(ctx context.Context, session models.Session) (models.Session, error)
	GetSession(ctx context.Context, id string) (models.Session, error)
	ListSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type sessionUsecase struct {
	sessionRepo SessionRepo
}

func NewSessionUsecase(repo SessionRepo) SessionUsecase {
	return &sessionUsecase{sessionRepo: repo}
}

func (u *sessionUsecase) RegisterSession(ctx context.Context, session models.Session) (models.Session, error) {
	return u.sessionRepo.CreateSession(ctx, session)
}

func (u *sessionUsecase) GetSession(ctx context.Context, id string) (models.Session, error) {
	return u.sessionRepo.GetSessionByID(ctx, id)
}

func (u *sessionUsecase) ListSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error) {
	return u.sessionRepo.GetSessionsByUserID(ctx, userID)
}

func (u *sessionUsecase) DeleteSession(ctx context.Context, id string) error {
	return u.sessionRepo.DeleteSession(ctx, id)
}
