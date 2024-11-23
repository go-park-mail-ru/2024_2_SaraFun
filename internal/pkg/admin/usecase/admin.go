package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/admin/repo"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/hashing"
)

type AdminUsecase interface {
	Login(username, password string) (*models.AdminUser, error)
	CreateAdmin(user *models.AdminUser) error
}

type adminUsecase struct {
	adminRepo repo.AdminUserRepo
}

func NewAdminUsecase(ar repo.AdminUserRepo) AdminUsecase {
	return &adminUsecase{
		adminRepo: ar,
	}
}

func (u *adminUsecase) Login(username, password string) (*models.AdminUser, error) {
	admin, err := u.adminRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !hashing.CheckPasswordHash(password, admin.Password) {
		return nil, errors.New("invalid username or password")
	}

	return admin, nil
}

func (u *adminUsecase) CreateAdmin(user *models.AdminUser) error {
	hashedPassword, err := hashing.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return u.adminRepo.Create(user)
}
