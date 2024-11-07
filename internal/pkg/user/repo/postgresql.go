package repo

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"sparkit/internal/utils/consts"

	//sparkiterrors "sparkit/internal/errors"
	"fmt"
	"sparkit/internal/models"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var user models.User
	err := repo.DB.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return user, nil
}

func (repo *Storage) GetUserList(ctx context.Context, userId int) ([]models.User, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var users []models.User
	rows, err := repo.DB.Query("SELECT id, username FROM users WHERE id != $1", userId)
	if err != nil {
		return []models.User{}, fmt.Errorf("GetUserList err: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return []models.User{}, fmt.Errorf("GetUserList err during scanning")
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *Storage) GetProfileIdByUserId(ctx context.Context, userId int) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var profileId int
	err := repo.DB.QueryRow("SELECT profile FROM users WHERE id=$1", userId).Scan(&profileId)
	if err != nil {
		return -1, fmt.Errorf("GetProfileIdByUserId err: %v", err)
	}
	return profileId, nil
}

func (repo *Storage) GetUsernameByUserId(ctx context.Context, userId int) (string, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var username string
	err := repo.DB.QueryRow("SELECT username FROM users WHERE id=$1", userId).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return username, nil
}

func (repo *Storage) GetFeedList(ctx context.Context, userId int, receivers []int) ([]models.User, error) {
	rows, err := repo.DB.Query("SELECT id, username FROM users WHERE id != $1 AND id NOT IN (SELECT receiver FROM reaction WHERE author = $2)", userId, userId)
	if err != nil {
		return []models.User{}, fmt.Errorf("GetFeedList err: %v", err)
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return []models.User{}, fmt.Errorf("GetFeedList err during scanning")
		}
		users = append(users, user)
	}
	return users, nil
}
