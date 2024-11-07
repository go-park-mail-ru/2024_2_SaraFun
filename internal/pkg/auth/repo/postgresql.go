package repo

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) AddUser(ctx context.Context, user models.User) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var id int
	err := repo.DB.QueryRow("INSERT INTO users (username, password, profile) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Profile).Scan(&id)
	if err != nil {
		repo.logger.Error("failed to insert user", zap.Error(err))
		return -1, fmt.Errorf("AddUser err : %v: ", err)
	}
	return id, nil
}

func (repo *Storage) DeleteUser(ctx context.Context, username string) error {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	_, err := repo.DB.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("DeleteUser err: %v", err)
	}
	return nil
}
