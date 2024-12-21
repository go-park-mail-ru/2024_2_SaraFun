package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCreateProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		profile  models.Profile
		queryErr error
		wantID   int
		wantErr  error
	}{
		{
			name:     "successful create",
			profile:  models.Profile{FirstName: "John", LastName: "Doe", BirthdayDate: "2000-01-01", Gender: "male", Target: "friends", About: "Hello!"},
			queryErr: nil,
			wantID:   1,
			wantErr:  nil,
		},
		{
			name:     "error on create",
			profile:  models.Profile{FirstName: "Jane", LastName: "Doe", BirthdayDate: "1995-01-01", Gender: "female", Target: "friends", About: "Hi!"},
			queryErr: errors.New("some insert error"),
			wantID:   -1,
			wantErr:  fmt.Errorf("CreateProfile err: %v", errors.New("some insert error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectQuery("INSERT INTO profile").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("INSERT INTO profile").
					WithArgs(tt.profile.FirstName, tt.profile.LastName, tt.profile.BirthdayDate, tt.profile.Gender, tt.profile.Target, tt.profile.About).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantID))
			}

			id, err := storage.CreateProfile(ctx, tt.profile)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Equal(t, tt.wantID, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		id       int
		profile  models.Profile
		queryErr error
		wantErr  error
	}{
		{
			name:     "successful update",
			id:       1,
			profile:  models.Profile{FirstName: "John", LastName: "Doe", BirthdayDate: "2000-01-01", Gender: "male", Target: "friends", About: "Updated!"},
			queryErr: nil,
			wantErr:  nil,
		},
		{
			name:     "error on update",
			id:       1,
			profile:  models.Profile{FirstName: "Jane", LastName: "Doe", BirthdayDate: "1990-01-01", Gender: "female", Target: "friends", About: "Updated again!"},
			queryErr: errors.New("some update error"),
			wantErr:  fmt.Errorf("UpdateProfile err: %v", errors.New("some update error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectExec("UPDATE profile").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectExec("UPDATE profile").
					WithArgs(tt.profile.FirstName, tt.profile.LastName, tt.profile.BirthdayDate, tt.profile.Gender, tt.profile.Target, tt.profile.About, tt.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := storage.UpdateProfile(ctx, tt.id, tt.profile)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name        string
		id          int
		mockRow     *sqlmock.Rows
		queryErr    error
		wantProfile models.Profile
		wantErr     error
	}{
		{
			name:        "successful get",
			id:          1,
			mockRow:     sqlmock.NewRows([]string{"id", "firstname", "lastname", "birthday_date", "gender", "target", "about"}).AddRow(1, "John", "Doe", "2000-01-01", "male", "friends", "Hello!"),
			queryErr:    nil,
			wantProfile: models.Profile{ID: 1, FirstName: "John", LastName: "Doe", BirthdayDate: "2000-01-01", Gender: "male", Target: "friends", About: "Hello!"},
			wantErr:     nil,
		},
		{
			name:        "error on get",
			id:          2,
			mockRow:     nil,
			queryErr:    errors.New("some select error"),
			wantProfile: models.Profile{},
			wantErr:     fmt.Errorf("GetProfile err: %v", errors.New("some select error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT id, firstname, lastname, birthday_date, gender, target, about FROM profile WHERE").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT id, firstname, lastname, birthday_date, gender, target, about FROM profile WHERE").
					WithArgs(tt.id).
					WillReturnRows(tt.mockRow)
			}

			profile, err := storage.GetProfile(ctx, tt.id)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Equal(t, tt.wantProfile, profile)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantProfile, profile)
			}
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		id       int
		queryErr error
		wantErr  error
	}{
		{
			name:     "successful delete",
			id:       1,
			queryErr: nil,
			wantErr:  nil,
		},
		{
			name:     "error on delete",
			id:       2,
			queryErr: errors.New("some delete error"),
			wantErr:  fmt.Errorf("DeleteProfile err: %v", errors.New("some delete error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectExec("DELETE FROM profile").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectExec("DELETE FROM profile").WithArgs(tt.id).WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := storage.DeleteProfile(ctx, tt.id)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
