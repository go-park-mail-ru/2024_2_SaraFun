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
	query := `INSERT INTO balance (userID, balance) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	return nil
}

func (repo *Storage) AddDailyLikeCount(ctx context.Context, userID int, amount int) error {
	query := `INSERT INTO daily_likes (userID, likes_count) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	return nil
}

func (repo *Storage) AddPurchasedLikeCount(ctx context.Context, userID int, amount int) error {
	query := `INSERT INTO purchased_likes (userID, likes_count) VALUES ($1, $2)`
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
	repo.logger.Info("userID", zap.Int("userID", userID))
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
	query := `INSERT INTO product (title, description, imagelink, price, product_count) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := repo.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.ImageLink,
		product.Price, product.Count).Scan(&id)
	if err != nil {
		repo.logger.Error("CreateProduct db query error", zap.Error(err))
		return -1, fmt.Errorf("CreateProduct db query error: %w", err)
	}
	return id, nil
}

func (repo *Storage) GetProduct(ctx context.Context, title string) (models.Product, error) {
	query := `SELECT title, description, imagelink, price, product_count FROM product WHERE title = $1`
	var product models.Product
	repo.logger.Info("title", zap.String("title", title))
	err := repo.DB.QueryRowContext(ctx, query, title).Scan(&product.Title,
		&product.Description, &product.ImageLink, &product.Price, &product.Count)
	if err != nil {
		repo.logger.Error("GetProduct db query error", zap.Error(err))
		return models.Product{}, fmt.Errorf("GetProduct db query error: %w", err)
	}
	return product, nil
}

func (repo *Storage) UpdateProduct(ctx context.Context, title string, product models.Product) error {
	query := `UPDATE product SET title = $1, description = $2, imagelink = $3, price = $4, product_count = $5 WHERE title = $6`
	_, err := repo.DB.ExecContext(ctx, query, product.Title, product.Description, product.ImageLink,
		product.Price, product.Count, title)
	if err != nil {
		repo.logger.Error("UpdateProduct db exec error", zap.Error(err))
		return fmt.Errorf("UpdateProduct db exec error: %w", err)
	}
	return nil
}

func (repo *Storage) GetProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT title, description, imagelink, price, product_count FROM product`
	var products []models.Product
	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("GetProducts db query error", zap.Error(err))
		return products, fmt.Errorf("GetProducts db query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Title, &product.Description, &product.ImageLink, &product.Price, &product.Count)
		if err != nil {
			repo.logger.Error("GetProducts row scan error", zap.Error(err))
			return products, fmt.Errorf("GetProducts row scan error: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (repo *Storage) AddAward(ctx context.Context, award models.Award) error {
	query := `INSERT INTO award (day_number, award_type, award_count) VALUES ($1, $2, $3)`
	_, err := repo.DB.ExecContext(ctx, query, award.DayNumber, award.Type, award.Count)
	if err != nil {
		repo.logger.Error("AddAward db query error", zap.Error(err))
		return fmt.Errorf("AddAward db query error: %w", err)
	}
	return nil
}

func (repo *Storage) GetAwards(ctx context.Context) ([]models.Award, error) {
	query := `SELECT day_number, award_type, award_count FROM award`
	var awards []models.Award
	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		repo.logger.Error("GetAwards db query error", zap.Error(err))
		return awards, fmt.Errorf("GetAwards db query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var award models.Award
		err := rows.Scan(&award.DayNumber, &award.Type, &award.Count)
		if err != nil {
			repo.logger.Error("GetAwards row scan error", zap.Error(err))
			return awards, fmt.Errorf("GetAwards row scan error: %w", err)
		}
		awards = append(awards, award)
	}
	return awards, nil
}

func (repo *Storage) GetAwardByDayNumber(ctx context.Context, dayNumber int) (models.Award, error) {
	query := `SELECT day_number, award_type, award_count FROM award WHERE day_number = $1`
	var award models.Award
	err := repo.DB.QueryRowContext(ctx, query, dayNumber).Scan(&award.DayNumber, &award.Type, &award.Count)
	if err != nil {
		repo.logger.Error("GetAwardByDayNumber db query error", zap.Error(err))
		return models.Award{}, fmt.Errorf("GetAwardByDayNumber db query error: %w", err)
	}
	return award, nil
}

func (repo *Storage) AddActivity(ctx context.Context, activity models.Activity) error {
	query := `INSERT INTO user_activity (user_id, last_login, consecutive_days) VALUES ($1, $2, $3)`
	_, err := repo.DB.ExecContext(ctx, query, activity.UserID, activity.Last_Login, activity.Consecutive_days)
	if err != nil {
		repo.logger.Error("AddActivity db query error", zap.Error(err))
		return fmt.Errorf("AddActivity db query error: %w", err)
	}
	return nil
}

func (repo *Storage) GetActivity(ctx context.Context, userID int) (models.Activity, error) {
	query := `SELECT user_id, last_login, consecutive_days FROM user_activity WHERE user_id = $1`
	var activity models.Activity
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&activity.UserID, &activity.Last_Login, &activity.Consecutive_days)
	if err != nil {
		repo.logger.Error("GetActivity db query error", zap.Error(err))
		return models.Activity{}, fmt.Errorf("GetActivity db query error: %w", err)
	}
	return activity, nil
}

func (repo *Storage) GetActivityDay(ctx context.Context, userID int) (int, error) {
	query := `SELECT consecutive_days FROM user_activity WHERE user_id = $1`
	var day int
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&day)
	if err != nil {
		repo.logger.Error("GetActivity db query error", zap.Error(err))
		return -1, fmt.Errorf("GetActivity db query error: %w", err)
	}
	return day, nil
}

func (repo *Storage) UpdateActivity(ctx context.Context, userID int, activity models.Activity) error {
	query := `UPDATE user_activity SET last_login = $1, consecutive_days = $2 WHERE user_id = $3`
	_, err := repo.DB.ExecContext(ctx, query, activity.Last_Login, activity.Consecutive_days, userID)
	if err != nil {
		repo.logger.Error("UpdateActivity db query error", zap.Error(err))
		return fmt.Errorf("UpdateActivity db query error: %w", err)
	}
	return nil
}
