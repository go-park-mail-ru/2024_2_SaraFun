package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
	"strconv"
	"time"
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
	AddAward(ctx context.Context, award models.Award) error
	GetAwards(ctx context.Context) ([]models.Award, error)
	GetAwardByDayNumber(ctx context.Context, dayNumber int) (models.Award, error)
	AddActivity(ctx context.Context, activity models.Activity) error
	GetActivityDay(ctx context.Context, userID int) (int, error)
	GetActivity(ctx context.Context, userID int) (models.Activity, error)
	UpdateActivity(ctx context.Context, userID int, activity models.Activity) error
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
	product.ImageLink = "likes" + strconv.Itoa(product.Count) + ".png"
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

func (u *UseCase) GetAwards(ctx context.Context) ([]models.Award, error) {
	awards, err := u.repo.GetAwards(ctx)
	if err != nil {
		u.logger.Error("bad get awards", zap.Error(err))
		return []models.Award{}, fmt.Errorf("bad get awards: %w", err)
	}
	return awards, nil
}

func (u *UseCase) AddAward(ctx context.Context, award models.Award) error {
	award.Sanitize()
	if award.Count < 0 {
		award.Count = 0
	}
	if award.DayNumber < 0 {
		award.DayNumber = 0
	}
	err := u.repo.AddAward(ctx, award)
	if err != nil {
		u.logger.Error("bad add award", zap.Error(err))
		return fmt.Errorf("bad add award: %w", err)
	}
	return nil
}

func (u *UseCase) GetAwardByDayNumber(ctx context.Context, dayNumber int) (models.Award, error) {
	award, err := u.repo.GetAwardByDayNumber(ctx, dayNumber)
	if err != nil {
		u.logger.Error("bad get award by day number", zap.Error(err))
		return models.Award{}, fmt.Errorf("bad get award by day number: %w", err)
	}
	return award, nil
}

func (u *UseCase) AddActivity(ctx context.Context, userID int) error {
	activity := models.Activity{
		Last_Login:       time.Now().Format(time.DateTime),
		Consecutive_days: 1,
		UserID:           userID,
	}
	u.logger.Info("time", zap.Any("activity", activity.Last_Login))
	err := u.repo.AddActivity(ctx, activity)
	if err != nil {
		u.logger.Error("bad add activity", zap.Error(err))
		return fmt.Errorf("bad add activity: %w", err)
	}
	return nil
}

func (u *UseCase) GetActivity(ctx context.Context, userID int) (models.Activity, error) {
	activity, err := u.repo.GetActivity(ctx, userID)
	if err != nil {
		u.logger.Error("bad get activity", zap.Error(err))
		return models.Activity{}, fmt.Errorf("bad get activity: %w", err)
	}
	return activity, nil
}

func (u *UseCase) UpdateActivity(ctx context.Context, userID int) (string, error) {
	var answer string
	activity, err := u.repo.GetActivity(ctx, userID)
	if err != nil {
		u.logger.Error("bad get activity", zap.Error(err))
		return "", fmt.Errorf("bad get activity: %w", err)
	}
	currentConDay := activity.Consecutive_days
	u.logger.Info("day", zap.Any("day", currentConDay))
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	activity_date, err := time.Parse(time.RFC3339, activity.Last_Login)
	if err != nil {
		u.logger.Error("bad get activity last login", zap.Error(err))
		return "", fmt.Errorf("bad get activity last login: %w", err)
	}
	if activity_date.Truncate(24 * time.Hour).Equal(yesterday) {
		activity.Consecutive_days++
	} else if activity_date.Truncate(24 * time.Hour).Before(yesterday) {
		activity.Consecutive_days = 1
	}
	u.logger.Info("after check cons days", zap.Any("day", activity.Consecutive_days))
	activity.Last_Login = now.Format(time.RFC3339)
	err = u.repo.UpdateActivity(ctx, userID, activity)
	if err != nil {
		u.logger.Error("bad update activity", zap.Error(err))
		return "", fmt.Errorf("bad update activity: %w", err)
	}

	if activity.Consecutive_days > currentConDay {
		answer, err = u.Reward(ctx, userID)
		u.logger.Info("get reward answer", zap.Any("answer", answer))
		u.logger.Info("reward success", zap.String("answer", answer))
		if err != nil {
			u.logger.Error("reward error", zap.Error(err))
			return "", fmt.Errorf("reward error: %w", err)
		}
	} else {
		answer = "Если зайдете завтра, то получите подарок!"
	}
	return answer, nil
}

func (u *UseCase) GetActivityDay(ctx context.Context, userID int) (int, error) {
	day, err := u.repo.GetActivityDay(ctx, userID)
	if err != nil {
		u.logger.Error("bad get activity day", zap.Error(err))
		return -1, fmt.Errorf("bad get activity day: %w", err)
	}
	dayNumber := day % 7
	return dayNumber, nil
}

func (u *UseCase) Reward(ctx context.Context, userID int) (string, error) {
	var answer string
	// должны определить, какой день подряд заходит, взять от него остаток от деления на 7
	dayNumber, err := u.GetActivityDay(ctx, userID)
	if err != nil {
		u.logger.Error("bad get activity day", zap.Error(err))
		return "", fmt.Errorf("bad get activity day: %w", err)
	}

	// по этому остатку определить тип и количество награды
	award, err := u.repo.GetAwardByDayNumber(ctx, dayNumber)
	if err != nil {
		u.logger.Error("bad get award by day number", zap.Error(err))
		return "", fmt.Errorf("bad get award by day number: %w", err)
	}

	// в зависимости от типа награды изменить баланс в нужной таблице (например платные лайки)
	if award.Type == "likes" {
		err = u.ChangePurchasedLikeCount(ctx, userID, award.Count)
		if err != nil {
			u.logger.Error("bad change like count", zap.Error(err))
			return "", fmt.Errorf("bad change like count: %w", err)
		}
		answer = "Вы получили в качестве подарка " + strconv.Itoa(award.Count) + " лайков за активное пользование нашим сайтом!"
	}
	return answer, nil
}
