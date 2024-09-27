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
