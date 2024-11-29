package user

import (
	"context"
	"errors"
	"fmt"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/usecase/user/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/hashing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func TestGetFeed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	defer logger.Sync()

	tests := []struct {
		name        string
		userId      int
		receivers   []int
		returnUsers []models.User
		returnError error
		returnCount int
		wantUsers   []models.User
		logger      *zap.Logger
	}{
		{
			name:        "successfull test",
			userId:      1,
			receivers:   []int{2, 3},
			returnUsers: []models.User{{ID: 1, Username: "Kirill", Password: "123456"}, {ID: 2, Username: "Andrey", Password: "222222"}},
			returnError: nil,
			returnCount: 1,
			wantUsers:   []models.User{{ID: 1, Username: "Kirill", Password: "123456"}, {ID: 2, Username: "Andrey", Password: "222222"}},
			logger:      logger,
		},
		{
			name:        "bad test",
			userId:      1,
			receivers:   []int{2, 3},
			returnUsers: nil,
			returnError: errors.New("test error"),
			returnCount: 1,
			wantUsers:   nil,
			logger:      logger,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetFeedList(ctx, tt.userId, tt.receivers).Return(tt.returnUsers, tt.returnError).Times(tt.returnCount)

			s := New(repo, logger)

			list, err := s.GetFeedList(ctx, tt.userId, tt.receivers)
			t.Log(err)
			t.Log(tt.returnError)
			require.ErrorIs(t, err, tt.returnError)
			for i, v := range list {
				if v != tt.wantUsers[i] {
					t.Errorf("GetFeedList() test error got = %v, want %v", list, tt.wantUsers[i])
				}
			}
		})
	}
}

func TestGetUserIdByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger := zap.NewNop()

	tests := []struct {
		name        string
		username    string
		returnId    int
		returnError error
		wantId      int
		wantErr     error
	}{
		{
			name:        "Successful test",
			username:    "testuser",
			returnId:    1,
			returnError: nil,
			wantId:      1,
			wantErr:     nil,
		},
		{
			name:        "Error case",
			username:    "nonexistent",
			returnId:    -1,
			returnError: errors.New("not found"),
			wantId:      -1,
			wantErr:     sparkiterrors.ErrWrongCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetUserIdByUsername(ctx, tt.username).Return(tt.returnId, tt.returnError).Times(1)

			u := New(repo, logger)
			userId, err := u.GetUserIdByUsername(ctx, tt.username)
			require.Equal(t, tt.wantId, userId)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCheckUsernameExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger := zap.NewNop()

	tests := []struct {
		name        string
		username    string
		returnExist bool
		returnError error
		wantExist   bool
		wantErr     error
	}{
		{
			name:        "Successful test - Exists",
			username:    "testuser",
			returnExist: true,
			returnError: nil,
			wantExist:   true,
			wantErr:     nil,
		},
		{
			name:        "Successful test - Does Not Exist",
			username:    "newuser",
			returnExist: false,
			returnError: nil,
			wantExist:   false,
			wantErr:     nil,
		},
		{
			name:        "Error case",
			username:    "erroruser",
			returnExist: false,
			returnError: errors.New("database error"),
			wantExist:   false,
			wantErr:     sparkiterrors.ErrWrongCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().CheckUsernameExists(ctx, tt.username).Return(tt.returnExist, tt.returnError).Times(1)

			u := New(repo, logger)
			exists, err := u.CheckUsernameExists(ctx, tt.username)
			require.Equal(t, tt.wantExist, exists)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetProfileIdByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger := zap.NewNop()

	tests := []struct {
		name          string
		userID        int
		returnProfile int
		returnError   error
		wantProfile   int
		wantError     bool
	}{
		{
			name:          "Successful test",
			userID:        1,
			returnProfile: 123,
			returnError:   nil,
			wantProfile:   123,
			wantError:     false,
		},
		{
			name:          "Error case",
			userID:        2,
			returnProfile: -1,
			returnError:   fmt.Errorf("database error"),
			wantProfile:   -1,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetProfileIdByUserId(ctx, tt.userID).Return(tt.returnProfile, tt.returnError).Times(1)

			u := New(repo, logger)

			profileID, err := u.GetProfileIdByUserId(ctx, tt.userID)
			if tt.wantError {
				require.Error(t, err, "Expected error, but got none")
			} else {
				require.NoError(t, err, "Did not expect an error, but got one")
			}
			require.Equal(t, tt.wantProfile, profileID, "Profile ID mismatch")
		})
	}
}
