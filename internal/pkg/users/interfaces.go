package pkg

import (
	"context"
	"errors"
	"sparkit/cmd/hashing"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"time"
)

//type UserRepo interface {
//	CreateUser(ctx context.Context, user models.User) (models.User, error)
//	GetUserByID(ctx context.Context, id string) (models.User, error)
//	GetUsers(ctx context.Context, count int64, offset int64) ([]models.User, error)
//}
//
//type inMemoryUserRepo struct {
//	users []models.User
//}
//
//func (repo *inMemoryUserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
//	user.ID = fmt.Sprintf("%d", len(repo.users)+1) // генерация ID
//	user.CreatedAt = time.Now()
//	repo.users = append(repo.users, user)
//	return user, nil
//}
//
//func (repo *inMemoryUserRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
//	for _, user := range repo.users {
//		if user.ID == id {
//			return user, nil
//		}
//	}
//	return models.User{}, fmt.Errorf("user not found")
//}
//
//func (repo *inMemoryUserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
//	for _, user := range repo.users {
//		if user.email == email {
//			return user, nil
//		}
//	}
//	return models.User{}, fmt.Errorf("user not found")
//}
//
//func (repo *inMemoryUserRepo) GetUsers(ctx context.Context, count int64, offset int64) ([]models.User, error) {
//	start := offset
//	end := offset + count
//	if start > int64(len(repo.users)) {
//		return []models.User{}, nil
//	}
//	if end > int64(len(repo.users)) {
//		end = int64(len(repo.users))
//	}
//	return repo.users[start:end], nil
//}
//
//type UserService interface {
//	RegisterUser(ctx context.Context, user models.User) (models.User, error)
//	GetUser(ctx context.Context, id string) (models.User, error)
//	ListUsers(ctx context.Context, count int64, offset int64) ([]models.User, error)
//}
//
//type userService struct {
//	userRepo UserRepo
//}
//
//func NewUserUsecase(repo UserRepo) UserService {
//	return &userService{userRepo: repo}
//}
//
//func (u *userService) RegisterUser(ctx context.Context, user models.User) (models.User, error) {
//	return u.userRepo.CreateUser(ctx, user)
//}
//
//func (u *userService) GetUser(ctx context.Context, id string) (models.User, error) {
//	return u.userRepo.GetUserByID(ctx, id)
//}
//
//func (u *userService) ListUsers(ctx context.Context, count int64, offset int64) ([]models.User, error) {
//	return u.userRepo.GetUsers(ctx, count, offset)
//}
//
////-----------------------------------------------------
////работа с сессией
//
//type SessionRepo interface {
//	CreateSession(ctx context.Context, session models.Session) (models.Session, error)
//	GetSessionByID(ctx context.Context, id string) (models.Session, error)
//	GetSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error)
//	DeleteSession(ctx context.Context, id string) error
//}
//
//type inMemorySessionRepo struct {
//	sessions []models.Session
//}
//
//func (repo *inMemorySessionRepo) CreateSession(ctx context.Context, session models.Session) (models.Session, error) {
//	session.ID = fmt.Sprintf("%d", len(repo.sessions)+1) // генерация ID
//	session.CreatedAt = time.Now()
//	repo.sessions = append(repo.sessions, session)
//	return session, nil
//}
//
//func (repo *inMemorySessionRepo) GetSessionByID(ctx context.Context, id string) (models.Session, error) {
//	for _, session := range repo.sessions {
//		if session.ID == id {
//			return session, nil
//		}
//	}
//	return models.Session{}, fmt.Errorf("session not found")
//}
//
//func (repo *inMemorySessionRepo) GetSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error) {
//	var userSessions []models.Session
//	for _, session := range repo.sessions {
//		if session.UserID == userID {
//			userSessions = append(userSessions, session)
//		}
//	}
//	return userSessions, nil
//}
//
//func (repo *inMemorySessionRepo) DeleteSession(ctx context.Context, id string) error {
//	for i, session := range repo.sessions {
//		if session.ID == id {
//			repo.sessions = append(repo.sessions[:i], repo.sessions[i+1:]...)
//			return nil
//		}
//	}
//	return fmt.Errorf("session not found")
//}
//
//type SessionService interface {
//	RegisterSession(ctx context.Context, session models.Session) (models.Session, error)
//	GetSession(ctx context.Context, id string) (models.Session, error)
//	ListSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error)
//	DeleteSession(ctx context.Context, id string) error
//}
//
//type sessionService struct {
//	sessionRepo SessionRepo
//}
//
//func NewSessionUsecase(repo SessionRepo) SessionService {
//	return &sessionService{sessionRepo: repo}
//}
//
//func (u *sessionService) RegisterSession(ctx context.Context, session models.Session) (models.Session, error) {
//	return u.sessionRepo.CreateSession(ctx, session)
//}
//
//func (u *sessionService) GetSession(ctx context.Context, id string) (models.Session, error) {
//	return u.sessionRepo.GetSessionByID(ctx, id)
//}
//
//func (u *sessionService) ListSessionsByUserID(ctx context.Context, userID string) ([]models.Session, error) {
//	return u.sessionRepo.GetSessionsByUserID(ctx, userID)
//}
//
//func (u *sessionService) DeleteSession(ctx context.Context, id string) error {
//	return u.sessionRepo.DeleteSession(ctx, id)
//}

type UserRepository interface {
	AddUser(ctx context.Context, user models.User) error
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

type SessionRepository interface {
	AddSession(ctx context.Context, session models.Session) error
	DeleteSessionByUserID(ctx context.Context, userID int) error
	GetSessionByUserId(ctx context.Context, userID int) (models.Session, error)
}

type InMemoryUserRepository struct {
	users []models.User
}

type InMemorySessionRepository struct {
	sessions []models.Session
}

func (repo *InMemoryUserRepository) AddUser(ctx context.Context, user models.User) error {
	repo.users = append(repo.users, user)
	return nil
}

func (repo *InMemoryUserRepository) DeleteUser(ctx context.Context, username string) error {
	for i, u := range repo.users {
		if u.Username == username {
			repo.users = append(repo.users[:i], repo.users[i+1:]...)
		}
	}
	return nil
}

func (repo *InMemoryUserRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	for _, u := range repo.users {
		if u.Username == username {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (repo *InMemorySessionRepository) AddSession(ctx context.Context, session models.Session) error {
	repo.sessions = append(repo.sessions, session)
	return nil
}

func (repo *InMemorySessionRepository) DeleteSessionByUserID(ctx context.Context, userID int) error {
	for i, u := range repo.sessions {
		if u.UserID == userID {
			repo.sessions = append(repo.sessions[:i], repo.sessions[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (repo *InMemorySessionRepository) GetSessionByUserID(ctx context.Context, userID int) (models.Session, error) {
	for _, u := range repo.sessions {
		if u.UserID == userID {
			return u, nil
		}
	}
	return models.Session{}, errors.New("session not found")
}

type UserService interface {
	RegisterUser(ctx context.Context, user models.User) error
	checkPassword(ctx context.Context, username string, password string) (models.User, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
	DeleteSessionByUserID(ctx context.Context, userID int) error
	GetSessionByUserId(ctx context.Context, userID int) (models.Session, error)
}

type userService struct {
	repo InMemoryUserRepository
}

type sessionService struct {
	repo InMemorySessionRepository
}

func (u *userService) RegisterUser(ctx context.Context, user models.User) error {
	err := u.repo.AddUser(ctx, user)
	if err != nil {
		return errors.New("failed to register user")
	}
	return nil
}

func (u *userService) CheckPassword(ctx context.Context, username string, password string) (models.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
	if hashing.CheckPasswordHash(password, user.Password) {
		return user, nil
	} else {
		return models.User{}, sparkiterrors.ErrWrongCredentials
	}
}

func (s *sessionService) CreateSession(ctx context.Context, user models.User) (models.Session, error) {
	session := models.Session{
		ID:        len(s.repo.sessions) + 1,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	err := s.repo.AddSession(ctx, session)
	if err != nil {
		return models.Session{}, errors.New("failed to create session")
	}
	return session, nil
}

func (s *sessionService) DeleteSessionByUserID(ctx context.Context, userID int) error {
	err := s.repo.DeleteSessionByUserID(ctx, userID)
	if err != nil {
		return errors.New("failed to delete session")
	}
	return nil
}

func (s *sessionService) GetSessionByUserID(ctx context.Context, userID int) (models.Session, error) {
	session, err := s.repo.GetSessionByUserID(ctx, userID)
	if err != nil {
		return models.Session{}, errors.New("failed to get session")
	}
	return session, nil
}

func NewUserService(repo InMemoryUserRepository) *userService {
	return &userService{repo: repo}
}

func NewSessionService(repo InMemorySessionRepository) *sessionService {
	return &sessionService{repo: repo}
}
