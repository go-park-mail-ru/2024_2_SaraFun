package user

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/usecase/user/mocks"
	"sparkit/internal/utils/consts"
	"sparkit/internal/utils/hashing"
	"testing"
	"time"
)

func TestRegisterUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	defer logger.Sync()
	user1 := models.User{ID: 1}
	user2 := models.User{ID: 2}
	tests := []struct {
		name   string
		user   models.User
		logger *zap.Logger
		want   error
	}{
		{
			name:   "successfull test",
			user:   user1,
			logger: logger,
			want:   nil,
		},
		{
			name:   "bad test",
			user:   user2,
			logger: logger,
			want:   sparkiterrors.ErrRegistrationUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().AddUser(ctx, tt.user).Return(tt.user.ID, tt.want).Times(1)
			u := New(repo, tt.logger)
			_, res := u.RegisterUser(ctx, tt.user)
			if res != tt.want {
				t.Errorf("RegisterUser() = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	defer logger.Sync()
	password1, _ := hashing.HashPassword("123456")
	password2, _ := hashing.HashPassword("222222")
	user1 := models.User{ID: 1, Username: "Kirill", Password: password1}
	user2 := models.User{ID: 2, Username: "Andrey", Password: password2}

	tests := []struct {
		name             string
		user             models.User
		password         string
		getUserError     error
		getUserWant      models.User
		getUserCallCount int
		logger           *zap.Logger
		wantUser         models.User
		wantErr          error
	}{
		{
			name:             "successfull test",
			user:             user1,
			password:         "123456",
			getUserError:     nil,
			getUserWant:      user1,
			getUserCallCount: 1,
			logger:           logger,
			wantUser:         user1,
			wantErr:          nil,
		},
		{
			name:             "bad test",
			user:             user2,
			password:         "333333",
			getUserError:     nil,
			getUserWant:      user2,
			getUserCallCount: 1,
			logger:           logger,
			wantUser:         models.User{},
			wantErr:          sparkiterrors.ErrWrongCredentials,
		},
		{
			name:             "bad username test",
			user:             models.User{Username: "Alexey"},
			password:         "123456",
			getUserError:     sparkiterrors.ErrBadUsername,
			getUserWant:      models.User{},
			getUserCallCount: 1,
			logger:           logger,
			wantUser:         models.User{},
			wantErr:          sparkiterrors.ErrWrongCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetUserByUsername(ctx, tt.user.Username).Return(tt.getUserWant, tt.getUserError).Times(tt.getUserCallCount)
			u := New(repo, tt.logger)
			res, err := u.CheckPassword(ctx, tt.user.Username, tt.password)
			if err != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res != tt.wantUser {
				t.Errorf("CheckPassword() = %v, want %v", res, tt.wantUser)
			}
		})
	}
}
