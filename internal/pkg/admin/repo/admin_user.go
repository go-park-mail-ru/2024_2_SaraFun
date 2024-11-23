package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
)

type AdminUserRepo interface {
	GetByUsername(username string) (*models.AdminUser, error)
	Create(user *models.AdminUser) error
}

type adminUserRepo struct {
	db *sql.DB
}

func NewAdminUserRepo(db *sql.DB) AdminUserRepo {
	return &adminUserRepo{db}
}

func (r *adminUserRepo) GetByUsername(username string) (*models.AdminUser, error) {
	var user models.AdminUser
	query := `SELECT id, username, password, role FROM admin_users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepo) Create(user *models.AdminUser) error {
	query := `INSERT INTO admin_users (username, password, role) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Username, user.Password, user.Role)
	return err
}
