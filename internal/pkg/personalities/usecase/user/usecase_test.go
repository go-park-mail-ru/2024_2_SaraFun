package user

import (
	"context"
	"errors"
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
	//defer logger.Sync()
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
	//defer logger.Sync()
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
	//defer logger.Sync()

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

func TestGetUserList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name          string
		userId        int
		repoReturns   []models.User
		repoError     error
		repoCount     int
		expectedUsers []models.User
	}{
		{
			name:          "successfull test",
			userId:        1,
			repoReturns:   []models.User{{ID: 1, Username: "Kirill", Password: "123456"}},
			repoError:     nil,
			repoCount:     1,
			expectedUsers: []models.User{{ID: 1, Username: "Kirill", Password: "123456"}},
		},
		{
			name:          "bad test",
			userId:        1,
			repoReturns:   nil,
			repoError:     errors.New("test error"),
			repoCount:     1,
			expectedUsers: []models.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetUserList(ctx, tt.userId).Return(tt.repoReturns, tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			list, err := s.GetUserList(ctx, tt.userId)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedUsers, list)
		})
	}
}

func TestGetProfileIdBYUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name       string
		userId     int
		repoReturn int
		repoError  error
		repoCount  int
		expectedID int
	}{
		{
			name:       "successfull test",
			userId:     1,
			repoReturn: 1,
			repoError:  nil,
			repoCount:  1,
			expectedID: 1,
		},
		{
			name:       "bad test",
			userId:     1,
			repoReturn: -1,
			repoError:  errors.New("test error"),
			repoCount:  1,
			expectedID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetProfileIdByUserId(ctx, tt.userId).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			id, err := s.GetProfileIdByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestGetUsernameBYUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name             string
		userId           int
		repoReturn       string
		repoError        error
		repoCount        int
		expectedUsername string
		expectedError    error
	}{
		{
			name:             "successfull test",
			userId:           1,
			repoReturn:       "Kirill",
			repoError:        nil,
			repoCount:        1,
			expectedUsername: "Kirill",
			expectedError:    nil,
		},
		{
			name:             "bad test",
			userId:           1,
			repoReturn:       "",
			repoError:        errors.New("test error"),
			repoCount:        1,
			expectedUsername: "",
			expectedError:    sparkiterrors.ErrWrongCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetUsernameByUserId(ctx, tt.userId).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			username, err := s.GetUsernameByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedUsername, username)
		})
	}
}

func TestGetUserIDByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		username       string
		repoReturn     int
		repoError      error
		repoCount      int
		expectedUserID int
		expectedError  error
	}{
		{
			name:           "successfull test",
			username:       "Kirill",
			repoReturn:     1,
			repoError:      nil,
			repoCount:      1,
			expectedUserID: 1,
			expectedError:  nil,
		},
		{
			name:           "bad test",
			username:       "",
			repoReturn:     -1,
			repoError:      errors.New("test error"),
			repoCount:      1,
			expectedUserID: -1,
			expectedError:  sparkiterrors.ErrWrongCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetUserIdByUsername(ctx, tt.username).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			id, err := s.GetUserIdByUsername(ctx, tt.username)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedUserID, id)
		})
	}
}

func TestCheckUsernameExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		username       string
		repoReturn     bool
		repoError      error
		repoCount      int
		expectedExists bool
		expectedError  error
	}{
		{
			name:           "successfull test",
			username:       "Kirill",
			repoReturn:     true,
			repoError:      nil,
			repoCount:      1,
			expectedExists: true,
			expectedError:  nil,
		},
		{
			name:           "bad test",
			username:       "",
			repoReturn:     true,
			repoError:      errors.New("test error"),
			repoCount:      1,
			expectedExists: false,
			expectedError:  sparkiterrors.ErrWrongCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().CheckUsernameExists(ctx, tt.username).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			exists, err := s.CheckUsernameExists(ctx, tt.username)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestChangePassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name      string
		userId    int
		password  string
		repoError error
		repoCount int
	}{
		{
			name:      "successfull test",
			userId:    1,
			password:  "123456",
			repoError: nil,
			repoCount: 1,
		},
		{
			name:      "bad test",
			userId:    1,
			password:  "",
			repoError: errors.New("test error"),
			repoCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().ChangePassword(ctx, tt.userId, gomock.Any()).Return(tt.repoError).Times(tt.repoCount)
			s := New(repo, logger)
			err := s.ChangePassword(ctx, tt.userId, tt.password)
			require.ErrorIs(t, err, tt.repoError)
		})
	}
}
