package profile

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

func (repo *Storage) CreateProfile(ctx context.Context, profile models.Profile) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var res int
	err := repo.DB.QueryRow("INSERT INTO profile (firstname, lastname, birthday_date, gender, target, about) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		profile.FirstName, profile.LastName, profile.BirthdayDate, profile.Gender, profile.Target, profile.About).Scan(&res)
	if err != nil {
		repo.logger.Error("error inserting profile", zap.Error(err))
		return -1, fmt.Errorf("CreateProfile err: %v", err)
	}
	id := res
	return id, nil
}
func (repo *Storage) UpdateProfile(ctx context.Context, id int, profile models.Profile) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	repo.logger.Info("id is", zap.Int("id", id))
	repo.logger.Info("profile is", zap.Any("profile", profile))
	_, err := repo.DB.Exec(`UPDATE profile SET firstname= $1,
                   lastname= $2,
                   birthday_date = $3,
                   gender = $4,
                   target = $5,
                   about = $6
                   WHERE id = $7`,
		profile.FirstName, profile.LastName, profile.BirthdayDate, profile.Gender, profile.Target, profile.About, id)
	if err != nil {
		repo.logger.Error("error updating profile", zap.Error(err))
		return fmt.Errorf("UpdateProfile err: %v", err)
	}
	return nil
}

func (repo *Storage) GetProfile(ctx context.Context, id int) (models.Profile, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var profile models.Profile
	err := repo.DB.QueryRow("SELECT id, firstname, lastname, birthday_date, gender, target, about FROM profile WHERE (id) = $1", id).Scan(&profile.ID,
		&profile.FirstName, &profile.LastName, &profile.BirthdayDate, &profile.Gender, &profile.Target, &profile.About)
	if err != nil {
		repo.logger.Error("error getting profile", zap.Error(err))
		return models.Profile{}, fmt.Errorf("GetProfile err: %v", err)
	}
	return profile, nil
}

func (repo *Storage) DeleteProfile(ctx context.Context, id int) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	_, err := repo.DB.Exec("DELETE FROM profile WHERE (id) = $1", id)
	if err != nil {
		repo.logger.Error("error deleting profile", zap.Error(err))
		return fmt.Errorf("DeleteProfile err: %v", err)
	}
	return nil
}
