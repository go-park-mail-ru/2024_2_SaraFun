package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddBalance(ctx context.Context, userID int, amount int) error
	AddDailyLikeCount(ctx context.Context, userID int, amount int) error
	AddPurchasedLikeCount(ctx context.Context, userID int, amount int) error
	ChangeBalance(ctx context.Context, userID int, amount int) error
	ChangeDailyLikeCount(ctx context.Context, userID int, amount int) error
	ChangePurchasedLikeCount(ctx context.Context, userID int, amount int) error
	SetBalance(ctx context.Context, userID int, balance int) error
	SetDailyLikesCount(ctx context.Context, userID int, balance int) error
	SetPurchasedLikesCount(ctx context.Context, userID int, balance int) error
	SetDailyLikesCountToAll(ctx context.Context, amount int) error
	GetBalance(ctx context.Context, userID int) (int, error)
	GetDailyLikesCount(ctx context.Context, userID int) (int, error)
	GetPurchasedLikesCount(ctx context.Context, userID int) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	GetProduct(ctx context.Context, title string) (models.Product, error)
	UpdateProduct(ctx context.Context, title string, product models.Product) error
	GetProducts(ctx context.Context) ([]models.Product, error)
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

func (u *UseCase) AddBalance(ctx context.Context, userID int, amount int) error {
	err := u.repo.AddBalance(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change balance", zap.Error(err))
		return fmt.Errorf("failed to change balance: %w", err)
	}
	return nil
}

func (u *UseCase) AddDailyLikesCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.AddDailyLikeCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change balance", zap.Error(err))
		return fmt.Errorf("failed to change balance: %w", err)
	}
	return nil
}

func (u *UseCase) AddPurchasedLikesCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.AddPurchasedLikeCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change balance", zap.Error(err))
		return fmt.Errorf("failed to change balance: %w", err)
	}
	return nil
}

func (u *UseCase) ChangeBalance(ctx context.Context, userID int, amount int) error {
	err := u.repo.ChangeBalance(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change balance", zap.Error(err))
		return fmt.Errorf("failed to change balance: %w", err)
	}
	return nil
}

func (u *UseCase) ChangeDailyLikeCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.ChangeDailyLikeCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change daily like count", zap.Error(err))
		return fmt.Errorf("failed to change daily like count: %w", err)
	}
	return nil
}

func (u *UseCase) ChangePurchasedLikeCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.ChangePurchasedLikeCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change purchased like count", zap.Error(err))
		return fmt.Errorf("failed to change purchased like count: %w", err)
	}
	return nil
}

func (u *UseCase) SetBalance(ctx context.Context, userID int, amount int) error {
	err := u.repo.SetBalance(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change purchased like count", zap.Error(err))
		return fmt.Errorf("failed to change purchased like count: %w", err)
	}
	return nil
}

func (u *UseCase) SetDailyLikeCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.SetDailyLikesCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change purchased like count", zap.Error(err))
		return fmt.Errorf("failed to change purchased like count: %w", err)
	}
	return nil
}

func (u *UseCase) SetDailyLikeCountToAll(ctx context.Context, amount int) error {
	err := u.repo.SetDailyLikesCountToAll(ctx, amount)
	if err != nil {
		u.logger.Warn("failed to change purchased like count", zap.Error(err))
		return fmt.Errorf("failed to change purchased like count: %w", err)
	}
	return nil
}

func (u *UseCase) SetPurchasedLikeCount(ctx context.Context, userID int, amount int) error {
	err := u.repo.SetPurchasedLikesCount(ctx, userID, amount)
	if err != nil {
		u.logger.Warn("failed to change purchased like count", zap.Error(err))
		return fmt.Errorf("failed to change purchased like count: %w", err)
	}
	return nil
}

func (u *UseCase) GetBalance(ctx context.Context, userID int) (int, error) {
	amount, err := u.repo.GetBalance(ctx, userID)
	if err != nil {
		u.logger.Error("usecase get balance error", zap.Error(err))
		return -1, fmt.Errorf("failed to get balance: %w", err)
	}
	return amount, err
}

func (u *UseCase) GetDailyLikesCount(ctx context.Context, userID int) (int, error) {
	amount, err := u.repo.GetDailyLikesCount(ctx, userID)
	if err != nil {
		u.logger.Error("usecase get balance error", zap.Error(err))
		return -1, fmt.Errorf("failed to get balance: %w", err)
	}
	return amount, err
}

func (u *UseCase) GetPurchasedLikesCount(ctx context.Context, userID int) (int, error) {
	amount, err := u.repo.GetPurchasedLikesCount(ctx, userID)
	if err != nil {
		u.logger.Error("usecase get balance error", zap.Error(err))
		return -1, fmt.Errorf("failed to get balance: %w", err)
	}
	return amount, err
}

func (u *UseCase) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	if product.Price < 0 {
		u.logger.Error("usecase create product bad price", zap.Int("price", product.Price))
		return -1, fmt.Errorf("invalid price")
	}
	product.ImageLink = product.Title + ".png"
	id, err := u.repo.CreateProduct(ctx, product)
	if err != nil {
		u.logger.Error("usecase create product error", zap.Error(err))
		return -1, fmt.Errorf("failed to create product: %w", err)
	}
	return id, err
}

func (u *UseCase) GetProduct(ctx context.Context, title string) (models.Product, error) {
	profile, err := u.repo.GetProduct(ctx, title)
	if err != nil {
		u.logger.Error("usecase create product error", zap.Error(err))
		return models.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return profile, err
}

func (u *UseCase) UpdateProduct(ctx context.Context, title string, product models.Product) error {
	if product.Price < 0 {
		u.logger.Error("usecase create product bad price", zap.Int("price", product.Price))
		return fmt.Errorf("invalid price: %v", product.Price)
	}
	err := u.repo.UpdateProduct(ctx, title, product)
	if err != nil {
		u.logger.Error("usecase create product error", zap.Error(err))
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (u *UseCase) CheckBalance(ctx context.Context, userID int, needMoney int) error {
	amount, err := u.repo.GetBalance(ctx, userID)
	if err != nil {
		u.logger.Error("usecase get balance error", zap.Error(err))
		return fmt.Errorf("failed to get balance: %w", err)
	}
	if amount < needMoney {
		u.logger.Error("usecase get balance error", zap.Int("amount", amount))
		return errors.New("Недостаточно средств")
	}
	return err
}

func (u *UseCase) GetProducts(ctx context.Context) ([]models.Product, error) {
	profiles, err := u.repo.GetProducts(ctx)
	if err != nil {
		u.logger.Error("usecase create product error", zap.Error(err))
		return []models.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return profiles, err
}
