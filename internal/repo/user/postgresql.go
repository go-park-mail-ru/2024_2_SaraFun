package user

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"sparkit/internal/models"
)

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (repo *Storage) AddUser(ctx context.Context, user models.User) error {
	//_, err := repo.DB.Exec("INSERT INTO users (username, password, age, gender) VALUES ($1, $2, $3, $4)", user.Username, user.Password, user.Age, user.Gender)
	sqlReq, args, err := sq.Insert("users").
		Columns("username", "password", "age", "gender").
		Values(user.Username, user.Password, user.Age, user.Gender).ToSql()
	_, err = repo.DB.Exec(sqlReq, args)
	if err != nil {
		return fmt.Errorf("AddUser err : %v: ", err)
	}
	return nil
}

func (repo *Storage) DeleteUser(ctx context.Context, username string) error {
	//_, err := repo.DB.Exec("DELETE FROM users WHERE username=$1", username)
	sqlReq, args, err := sq.Delete("users").
		Where(sq.Eq{"username": username}).ToSql()
	_, err = repo.DB.Exec(sqlReq, args)
	if err != nil {
		return fmt.Errorf("DeleteUser err: %v", err)
	}
	return nil
}

func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	//err := repo.DB.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	sqlReq, args, err := sq.Select("users").
		Columns("id", "username", "password").
		Where(sq.Eq{"username": username}).
		ToSql()
	err = repo.DB.QueryRow(sqlReq, args).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByUsername err: %v", err)
	}
	return user, nil
}

func (repo *Storage) GetUserList(ctx context.Context) ([]models.User, error) {
	var users []models.User
	sqlReq, args, sqErr := sq.Select("users").
		Columns("id", "username", "age", "gender").
		ToSql()
	if sqErr != nil {
		return []models.User{}, fmt.Errorf("GetUserList err: %v", sqErr)
	}
	//err = repo.DB.QueryRow(sqlReq, args).Scan(&user.ID, &user.Username, &user.Password)
	rows, err := repo.DB.Query(sqlReq, args)
	if err != nil {
		return []models.User{}, fmt.Errorf("GetUserList err: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Age, &user.Gender); err != nil {
			return []models.User{}, fmt.Errorf("GetUserList err during scanning")
		}
		users = append(users, user)
	}
	return users, nil
}
