package user

import (
	"context"
	"database/sql"
	"sparkit/internal/models"
)

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (repo *Storage) AddUser(ctx context.Context, user models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

func (repo *Storage) DeleteUser(ctx context.Context, username string) error {
	_, err := repo.DB.Exec("DELETE FROM users WHERE username=$1", username)
	return err
}

func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := repo.DB.QueryRow("SELECT * FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func (repo *Storage) GetUserList(ctx context.Context) ([]models.User, error) {
	var users []models.User
	rows, err := repo.DB.Query("SELECT * FROM users")
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return []models.User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}
