package user

import (
	"context"
	"database/sql/driver"
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

func TestGetUserList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	db, mock, err := sqlmock.New()
	logger := zap.NewNop()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	user1 := models.User{ID: 1}
	user2 := models.User{ID: 2}
	user3 := models.User{ID: 3}
	user4 := models.User{ID: 4}
	user5 := models.User{ID: 5}
	users := []models.User{user1, user2, user3, user4, user5}
	successRows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(user1.ID, user1.Username).
		AddRow(user2.ID, user2.Username).
		AddRow(user3.ID, user3.Username).
		AddRow(user4.ID, user4.Username).
		AddRow(user5.ID, user5.Username)
	badRows := sqlmock.NewRows([]string{"random"}).
		AddRow("1").
		AddRow("2").
		AddRow("3")
	tests := []struct {
		name             string
		resultQueryList  *sqlmock.Rows
		resultQueryError error
		wantList         []models.User
		wantErr          error
		logger           *zap.Logger
	}{
		{
			name:             "successful test",
			resultQueryList:  successRows,
			resultQueryError: nil,
			wantList:         users,
			wantErr:          nil,
			logger:           logger,
		},
		{
			name:             "bad test",
			resultQueryList:  nil,
			resultQueryError: errors.New("test"),
			wantList:         []models.User{},
			wantErr:          fmt.Errorf("GetUserList err: %v", errors.New("test")),
			logger:           logger,
		},
		{
			name:             "bad scanning",
			resultQueryList:  badRows,
			resultQueryError: nil,
			wantList:         []models.User{},
			wantErr:          fmt.Errorf("GetUserList err during scanning"),
			logger:           logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			if tt.resultQueryList != nil {
				mock.ExpectQuery("SELECT id, username FROM users").WillReturnRows(tt.resultQueryList)
			} else {
				mock.ExpectQuery("SELECT id, username FROM users").WillReturnError(tt.resultQueryError)
			}

			list, err := storage.GetUserList(ctx, 1)
			if err != nil && tt.wantErr != nil && (err.Error() != tt.wantErr.Error()) {
				//t.Errorf("GetUserList() error = %v, wantErr %v,", err, tt.wantErr)
				t.Errorf("GetUserList() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if (err != nil && tt.wantErr == nil) && (err == nil && tt.wantErr != nil) {
					t.Errorf("GetUserList() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if len(list) != len(tt.wantList) {
				t.Errorf("GetUserList() got = %v, want %v", list, tt.wantList)
			}
			for i, user := range list {
				if user != tt.wantList[i] {
					t.Errorf("GetUserList() got = %v, want %v", list, tt.wantList)
				}
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
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
		name        string
		username    string
		queryResult *sqlmock.Rows
		queryError  error
		wantUser    models.User
		wantErr     error
	}{
		{
			name:        "successful test",
			username:    "kirill",
			queryResult: sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", "kirill", "124"),
			queryError:  nil,
			wantUser:    models.User{ID: 1, Username: "kirill", Password: "124"},
			wantErr:     nil,
		},
		{
			name:        "bad test",
			username:    "andrey",
			queryResult: nil,
			queryError:  errors.New("bad username"),
			wantUser:    models.User{},
			wantErr:     fmt.Errorf("GetUserByUsername err: bad username"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			if tt.queryResult != nil {
				mock.ExpectQuery("SELECT id, username, password FROM users").WillReturnRows(tt.queryResult)
			} else {
				mock.ExpectQuery("SELECT id, username, password FROM users").WillReturnError(tt.queryError)
			}

			user, err := storage.GetUserByUsername(ctx, tt.username)

			if err != nil && tt.wantErr != nil && (err.Error() != tt.wantErr.Error()) {
				//t.Errorf("GetUserList() error = %v, wantErr %v,", err, tt.wantErr)
				t.Errorf("GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if (err != nil && tt.wantErr == nil) && (err == nil && tt.wantErr != nil) {
					t.Errorf("GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if user != tt.wantUser {
				t.Errorf("GetUserByUsername() got = %v, want %v", user, tt.wantUser)
			}
		})
	}
}

func TestProfileIdByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)

	tests := []struct {
		name              string
		userId            int
		queryRows         *sqlmock.Rows
		queryError        error
		expectedProfileID int
	}{
		{
			name:              "successfull test",
			userId:            1,
			queryRows:         sqlmock.NewRows([]string{"id"}).AddRow(1),
			queryError:        nil,
			expectedProfileID: 1,
		},
		{
			name:              "bad test",
			userId:            2,
			queryRows:         nil,
			queryError:        errors.New("test"),
			expectedProfileID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT profile FROM users").
					WithArgs(tt.userId).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT profile FROM users").
					WithArgs(tt.userId).
					WillReturnError(tt.queryError)
			}

			id, err := repo.GetProfileIdByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedProfileID, id)
		})
	}
}

func TestGetUsernameByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := New(db, logger)
	tests := []struct {
		name             string
		userId           int
		queryRows        *sqlmock.Rows
		queryError       error
		expectedUsername string
	}{
		{
			name:             "successfull test",
			userId:           1,
			queryRows:        sqlmock.NewRows([]string{"username"}).AddRow("Kirill"),
			queryError:       nil,
			expectedUsername: "Kirill",
		},
		{
			name:             "bad test",
			userId:           2,
			queryRows:        nil,
			queryError:       errors.New("test"),
			expectedUsername: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT username FROM users").
					WithArgs(tt.userId).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT username FROM users").
					WithArgs(tt.userId).
					WillReturnError(tt.queryError)
			}
			username, err := repo.GetUsernameByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedUsername, username)
		})
	}
}

func TestGetFeedList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name         string
		userID       int
		receivers    []int
		queryRows    *sqlmock.Rows
		queryError   error
		expectedList []models.User
	}{
		{
			name:      "successfull test",
			userID:    1,
			receivers: []int{1, 2, 3},
			queryRows: sqlmock.NewRows([]string{"id", "username"}).
				AddRow(1, "Король").
				AddRow(2, "Кирилл").
				AddRow(3, "Крепыш"),
			queryError: nil,
			expectedList: []models.User{
				{
					ID:       1,
					Username: "Король",
				},
				{
					ID:       2,
					Username: "Кирилл",
				},
				{
					ID:       3,
					Username: "Крепыш",
				},
			},
		},
		{
			name:         "bad test",
			userID:       2,
			receivers:    []int{1, 2, 3},
			queryRows:    nil,
			queryError:   errors.New("error"),
			expectedList: []models.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.userID).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.userID).
					WillReturnError(tt.queryError)
			}
			list, err := repo.GetFeedList(ctx, tt.userID, tt.receivers)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedList, list)
		})
	}
}

func TestGetUserIdByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name       string
		username   string
		queryRows  *sqlmock.Rows
		queryError error
		expectedID int
	}{
		{
			name:       "successfull test",
			username:   "Kirill",
			queryRows:  sqlmock.NewRows([]string{"id"}).AddRow(1),
			queryError: nil,
			expectedID: 1,
		},
		{
			name:       "bad test",
			username:   "",
			queryRows:  nil,
			queryError: errors.New("test"),
			expectedID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs(tt.username).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs(tt.username).
					WillReturnError(tt.queryError)
			}

			id, err := repo.GetUserIdByUsername(ctx, tt.username)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestCheckUsernameExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name           string
		username       string
		queryRows      *sqlmock.Rows
		queryError     error
		expectedExists bool
	}{
		{
			name:     "successfull test",
			username: "Kirill",
			queryRows: sqlmock.NewRows([]string{"exists"}).
				AddRow(true),
			queryError:     nil,
			expectedExists: true,
		},
		{
			name:           "bad test",
			username:       "",
			queryRows:      nil,
			queryError:     errors.New("test"),
			expectedExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT EXISTS").
					WithArgs(tt.username).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT EXISTS").
					WithArgs(tt.username).
					WillReturnError(tt.queryError)
			}
			exists, err := repo.CheckUsernameExists(ctx, tt.username)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestChangePassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name        string
		userId      int
		password    string
		queryResult driver.Result
		queryError  error
	}{
		{
			name:        "successfull test",
			userId:      1,
			password:    "123456",
			queryResult: sqlmock.NewResult(1, 1),
			queryError:  nil,
		},
		{
			name:        "bad test",
			userId:      1,
			password:    "123-456",
			queryResult: nil,
			queryError:  errors.New("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectExec("UPDATE").
					WithArgs(tt.password, tt.userId).
					WillReturnResult(tt.queryResult)
			} else {
				mock.ExpectExec("UPDATE").
					WithArgs(tt.password, tt.userId).
					WillReturnError(tt.queryError)
			}
			err := repo.ChangePassword(ctx, tt.userId, tt.password)
			require.ErrorIs(t, err, tt.queryError)
		})
	}
}
