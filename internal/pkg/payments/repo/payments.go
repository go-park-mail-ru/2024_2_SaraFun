package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{
		DB:     db,
		logger: logger,
	}
}

func (repo *Storage) AddBalance(ctx context.Context, userID int, amount int) error {
	query := `INSERT INTO balance (user_id, balance) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	return nil
}

func (repo *Storage) AddDailyLikeCount(ctx context.Context, userID int, amount int) error {
	query := `INSERT INTO daily_likes (user_id, likes_count) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	return nil
}

func (repo *Storage) AddPurchasedLikeCount(ctx context.Context, userID int, amount int) error {
	query := `INSERT INTO purchased_likes (user_id, likes_count) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	return nil
}

func (repo *Storage) ChangeBalance(ctx context.Context, userID int, amount int) error {
	query := `UPDATE balance SET balance = balance + $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, amount, userID)
	if err != nil {
		repo.logger.Error("ChangeBalance db exec error", zap.Error(err))
		return fmt.Errorf("ChangeBalance db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) ChangeDailyLikeCount(ctx context.Context, userID int, amount int) error {
	query := `UPDATE daily_likes SET likes_count = likes_count + $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, amount, userID)
	if err != nil {
		repo.logger.Error("ChangeDailyLikeCount db exec error", zap.Error(err))
		return fmt.Errorf("ChangeDailyLikeCount db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) ChangePurchasedLikeCount(ctx context.Context, userID int, amount int) error {
	query := `UPDATE purchased_likes SET likes_count = likes_count + $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, amount, userID)
	if err != nil {
		repo.logger.Error("ChangePurchasedLikeCount db exec error", zap.Error(err))
		return fmt.Errorf("ChangePurchasedLikeCount db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) SetBalance(ctx context.Context, userID int, balance int) error {
	query := `UPDATE balance SET balance = $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, balance, userID)
	if err != nil {
		repo.logger.Error("ChangeBalance db exec error", zap.Error(err))
		return fmt.Errorf("ChangeBalance db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) SetDailyLikesCountToAll(ctx context.Context, balance int) error {
	query := `UPDATE daily_likes SET likes_count = $1`
	_, err := repo.DB.ExecContext(ctx, query, balance)
	if err != nil {
		repo.logger.Error("ChangeBalance db exec error", zap.Error(err))
		return fmt.Errorf("ChangeBalance db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) SetDailyLikesCount(ctx context.Context, userID int, balance int) error {
	query := `UPDATE daily_likes SET likes_count = $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, balance, userID)
	if err != nil {
		repo.logger.Error("ChangeBalance db exec error", zap.Error(err))
		return fmt.Errorf("ChangeBalance db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) SetPurchasedLikesCount(ctx context.Context, userID int, balance int) error {
	query := `UPDATE purchased_likes SET likes_count = $1 WHERE userID = $2`
	_, err := repo.DB.ExecContext(ctx, query, balance, userID)
	if err != nil {
		repo.logger.Error("ChangeBalance db exec error", zap.Error(err))
		return fmt.Errorf("ChangeBalance db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) GetBalance(ctx context.Context, userID int) (int, error) {
	query := `SELECT balance FROM balance WHERE userID = $1`
	var amount int
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&amount)
	if err != nil {
		repo.logger.Error("GetBalance db query error", zap.Error(err))
		return -1, fmt.Errorf("GetBalance db query error: %w", err)
	}
	return amount, nil
}

func (repo *Storage) GetDailyLikesCount(ctx context.Context, userID int) (int, error) {
	query := `SELECT likes_count FROM daily_likes WHERE userID = $1`
	var amount int
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&amount)
	if err != nil {
		repo.logger.Error("GetBalance db query error", zap.Error(err))
		return -1, fmt.Errorf("GetBalance db query error: %w", err)
	}
	return amount, nil
}

func (repo *Storage) GetPurchasedLikesCount(ctx context.Context, userID int) (int, error) {
	query := `SELECT likes_count FROM purchased_likes WHERE userID = $1`
	var amount int
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&amount)
	if err != nil {
		repo.logger.Error("GetBalance db query error", zap.Error(err))
		return -1, fmt.Errorf("GetBalance db query error: %w", err)
	}
	return amount, nil
}

func (repo *Storage) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	query := `INSERT INTO product (title, description, imagelink, price) VALUES ($1, $2, $3, $4)`
	var id int
	err := repo.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.ImageLink, product.Price).Scan(&id)
	if err != nil {
		repo.logger.Error("CreateProduct db query error", zap.Error(err))
		return -1, fmt.Errorf("CreateProduct db query error: %w", err)
	}
	return id, nil
}

func (repo *Storage) GetProduct(ctx context.Context, title string) (models.Product, error) {
	query := `SELECT title, description, imagelink, price FROM product WHERE title = $1`
	var product models.Product
	err := repo.DB.QueryRowContext(ctx, query, title).Scan(&product.Title, &product.Description, &product.ImageLink, &product.Price)
	if err != nil {
		repo.logger.Error("GetProduct db query error", zap.Error(err))
		return models.Product{}, fmt.Errorf("GetProduct db query error: %w", err)
	}
	return product, nil
}

func (repo *Storage) UpdateProduct(ctx context.Context, title string, product models.Product) error {
	query := `UPDATE product SET title = $1, description = $2, imagelink = $3, price = $4 WHERE title = $5`
	_, err := repo.DB.ExecContext(ctx, query, product.Title, product.Description, product.ImageLink, product.Price, title)
	if err != nil {
		repo.logger.Error("UpdateProduct db exec error", zap.Error(err))
		return fmt.Errorf("UpdateProduct db exec error: %w", err)
	}
	return nil
}
