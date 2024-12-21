package user

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

	//sparkiterrors "sparkit/internal/errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) AddUser(ctx context.Context, user models.User) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
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
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	_, err := repo.DB.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("DeleteUser err: %v", err)
	}
	return nil
}

func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var user models.User
	err := repo.DB.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return user, nil
}

func (repo *Storage) GetUserList(ctx context.Context, userId int) ([]models.User, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
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
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var profileId int
	err := repo.DB.QueryRow("SELECT profile FROM users WHERE id=$1", userId).Scan(&profileId)
	if err != nil {
		return -1, fmt.Errorf("GetProfileIdByUserId err: %w", err)
	}
	return profileId, nil
}

func (repo *Storage) GetUsernameByUserId(ctx context.Context, userId int) (string, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var username string
	err := repo.DB.QueryRow("SELECT username FROM users WHERE id=$1", userId).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("GetUserByUsername err: %w", err)
	}
	return username, nil
}

func (repo *Storage) GetFeedList(ctx context.Context, userId int, receivers []int) ([]models.User, error) {

	query := `SELECT u.id, u.username
              FROM users u 
              JOIN profile p ON p.id = u.profile
              WHERE u.id != $1 AND
              u.id NOT IN (SELECT receiver FROM reaction WHERE author = $1) AND
              p.gender != (SELECT gender FROM profile WHERE id = (SELECT profile FROM users WHERE id = $1))`
	rows, err := repo.DB.Query(query, userId)
	if err != nil {
		return []models.User{}, fmt.Errorf("GetFeedList err: %w", err)
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

func (repo *Storage) GetUserIdByUsername(ctx context.Context, username string) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var userId int
	err := repo.DB.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userId)
	if err != nil {
		return -1, fmt.Errorf("GetUserByUsername err: %w", err)
	}
	return userId, nil
}

func (repo *Storage) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		repo.logger.Error("failed to check username unique", zap.Error(err))
		return false, fmt.Errorf("CheckUsernameUnique err: %w", err)
	}
	return exists, nil
}

func (repo *Storage) ChangePassword(ctx context.Context, userID int, password string) error {
	_, err := repo.DB.ExecContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, userID)
	if err != nil {
		repo.logger.Error("failed to change password", zap.Error(err))
		return fmt.Errorf("ChangePassword err: %w", err)
	}
	return nil
}
