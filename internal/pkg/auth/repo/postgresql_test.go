package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"testing"
	"time"
)

func TestAddUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	user1 := models.User{ID: 1, Username: "kirill", Password: "124", Profile: 1}

	tests := []struct {
		name     string
		user     models.User
		queryErr error
		wantErr  error
	}{
		{
			name:     "successful test",
			user:     user1,
			wantErr:  nil,
			queryErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			mock.ExpectExec("INSERT INTO users").
				WithArgs(tt.user.Username, tt.user.Password, tt.user.Profile).
				WillReturnError(tt.queryErr)
			//WillReturnResult(sqlmock.NewResult(1, 1))
			_, err := storage.AddUser(ctx, tt.user)
			if err != nil && tt.wantErr != nil && (err.Error() != tt.wantErr.Error()) {

				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if (err != nil && tt.wantErr == nil) && (err == nil && tt.wantErr != nil) {
					t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		username   string
		execResult error
		wantErr    error
	}{
		{
			name:       "successful test",
			username:   "kirill",
			execResult: nil,
			wantErr:    nil,
		},
		{
			name:       "bad test",
			username:   "andrey",
			execResult: errors.New("bad username"),
			wantErr:    fmt.Errorf("DeleteUser err: bad username"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			mock.ExpectExec("DELETE FROM users").WithArgs(tt.username).WillReturnError(tt.execResult)
			err := storage.DeleteUser(ctx, tt.username)
			//if err != tt.wantErr {
			//	t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			//}
			if err != nil && tt.wantErr != nil && (err.Error() != tt.wantErr.Error()) {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if (err != nil && tt.wantErr == nil) && (err == nil && tt.wantErr != nil) {
					t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
