package user

import (
	"context"
	_ "database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStorage_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		inputUser models.User
		mockQuery func()
		expected  int
		wantErr   bool
	}{
		{
			name: "Successful AddUser",
			inputUser: models.User{
				Username: "testuser",
				Password: "password123",
				Profile:  1,
			},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("testuser", "password123", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name: "Error in AddUser",
			inputUser: models.User{
				Username: "testuser",
				Password: "password123",
				Profile:  1,
			},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("testuser", "password123", 1).
					WillReturnError(errors.New("insert error"))
			},
			expected: -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			id, err := repo.AddUser(ctx, tt.inputUser)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, id)
		})
	}
}

func TestStorage_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name     string
		username string
		mockExec func()
		wantErr  bool
	}{
		{
			name:     "Successful DeleteUser",
			username: "testuser",
			mockExec: func() {
				mock.ExpectExec("DELETE FROM users").
					WithArgs("testuser").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:     "Error in DeleteUser",
			username: "testuser",
			mockExec: func() {
				mock.ExpectExec("DELETE FROM users").
					WithArgs("testuser").
					WillReturnError(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExec()
			err := repo.DeleteUser(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStorage_GetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		username  string
		mockQuery func()
		expected  models.User
		wantErr   bool
	}{
		{
			name:     "Successful GetUserByUsername",
			username: "testuser",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, username, password FROM users").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
						AddRow(1, "testuser", "password123"))
			},
			expected: models.User{
				ID:       1,
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name:     "Error in GetUserByUsername",
			username: "testuser",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, username, password FROM users").
					WithArgs("testuser").
					WillReturnError(errors.New("query error"))
			},
			expected: models.User{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			user, err := repo.GetUserByUsername(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, user)
		})
	}
}

func TestStorage_GetProfileIdByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		userId    int
		mockQuery func()
		expected  int
		wantErr   bool
	}{
		{
			name:   "Successful GetProfileIdByUserId",
			userId: 1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT profile FROM users WHERE id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"profile"}).
						AddRow(101))
			},
			expected: 101,
			wantErr:  false,
		},
		{
			name:   "Error in GetProfileIdByUserId",
			userId: 1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT profile FROM users WHERE id").
					WithArgs(1).
					WillReturnError(errors.New("query error"))
			},
			expected: -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			profileId, err := repo.GetProfileIdByUserId(ctx, tt.userId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, profileId)
		})
	}
}

func TestStorage_CheckUsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		username  string
		mockQuery func()
		expected  bool
		wantErr   bool
	}{
		{
			name:     "Successful CheckUsernameExists",
			username: "user1",
			mockQuery: func() {
				mock.ExpectQuery("SELECT EXISTS").
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).
						AddRow(true))
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Error in CheckUsernameExists",
			username: "user1",
			mockQuery: func() {
				mock.ExpectQuery("SELECT EXISTS").
					WithArgs("user1").
					WillReturnError(errors.New("query error"))
			},
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			exists, err := repo.CheckUsernameExists(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, exists)
		})
	}
}

func TestStorage_ChangePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name     string
		userId   int
		password string
		mockExec func()
		wantErr  bool
	}{
		{
			name:     "Successful ChangePassword",
			userId:   1,
			password: "newpassword123",
			mockExec: func() {
				mock.ExpectExec("UPDATE users SET password").
					WithArgs("newpassword123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "Error in ChangePassword",
			userId:   1,
			password: "newpassword123",
			mockExec: func() {
				mock.ExpectExec("UPDATE users SET password").
					WithArgs("newpassword123", 1).
					WillReturnError(errors.New("update error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExec()
			err := repo.ChangePassword(ctx, tt.userId, tt.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStorage_GetUserList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		userId    int
		mockQuery func()
		expected  []models.User
		wantErr   bool
	}{
		{
			name:   "Successful GetUserList",
			userId: 1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, username FROM users WHERE id !=").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(2, "user2").
						AddRow(3, "user3"))
			},
			expected: []models.User{
				{ID: 2, Username: "user2"},
				{ID: 3, Username: "user3"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			users, err := repo.GetUserList(ctx, tt.userId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, users)
		})
	}
}

func TestStorage_GetUsernameByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		userId    int
		mockQuery func()
		expected  string
		wantErr   bool
	}{
		{
			name:   "Successful GetUsernameByUserId",
			userId: 1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT username FROM users WHERE id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).
						AddRow("testuser"))
			},
			expected: "testuser",
			wantErr:  false,
		},
		{
			name:   "Error in GetUsernameByUserId",
			userId: 1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT username FROM users WHERE id").
					WithArgs(1).
					WillReturnError(errors.New("query error"))
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			username, err := repo.GetUsernameByUserId(ctx, tt.userId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, username)
		})
	}
}

func TestStorage_GetFeedList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	tests := []struct {
		name      string
		userId    int
		receivers []int
		mockQuery func()
		expected  []models.User
		wantErr   bool
	}{
		{
			name:      "Successful GetFeedList",
			userId:    1,
			receivers: []int{2, 3},
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, username FROM users WHERE id !=").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(4, "user4").
						AddRow(5, "user5"))
			},
			expected: []models.User{
				{ID: 4, Username: "user4"},
				{ID: 5, Username: "user5"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			users, err := repo.GetFeedList(ctx, tt.userId, tt.receivers)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, users)
		})
	}
}
