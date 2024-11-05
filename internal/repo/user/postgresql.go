package user

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

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

func (repo *Storage) AddUser(ctx context.Context, user models.User) (int64, error) {
	var id int64
	err := repo.DB.QueryRow("INSERT INTO users (username, password, profile) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Profile).Scan(&id)
	if err != nil {
		repo.logger.Error("failed to insert user", zap.Error(err))
		return -1, fmt.Errorf("AddUser err : %v: ", err)
	}
	return id, nil
}

func (repo *Storage) DeleteUser(ctx context.Context, username string) error {
	_, err := repo.DB.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("DeleteUser err: %v", err)
	}
	return nil
}

func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := repo.DB.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return user, nil
}

func (repo *Storage) GetUserList(ctx context.Context) ([]models.User, error) {
	var users []models.User
	rows, err := repo.DB.Query("SELECT id, username FROM users")
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

func (repo *Storage) GetProfileIdByUserId(ctx context.Context, userId int) (int64, error) {
	var profileId int64
	err := repo.DB.QueryRow("SELECT profile FROM users WHERE id=$1", userId).Scan(&profileId)
	if err != nil {
		return -1, fmt.Errorf("GetProfileIdByUserId err: %v", err)
	}
	return profileId, nil
}

func (repo *Storage) GetUsernameByUserId(ctx context.Context, userId int) (string, error) {
	var username string
	err := repo.DB.QueryRow("SELECT username FROM users WHERE id=$1", userId).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return username, nil
}
