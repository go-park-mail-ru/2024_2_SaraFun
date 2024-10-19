package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"sparkit/internal/models"
	"testing"
)

func TestGetUserList(t *testing.T) {
	db, mock, err := sqlmock.New()
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
	successRows := sqlmock.NewRows([]string{"id", "username", "age", "gender"}).
		AddRow(user1.ID, user1.Username, user1.Age, user1.Gender).
		AddRow(user2.ID, user2.Username, user2.Age, user2.Gender).
		AddRow(user3.ID, user3.Username, user3.Age, user3.Gender).
		AddRow(user4.ID, user4.Username, user4.Age, user4.Gender).
		AddRow(user5.ID, user5.Username, user5.Age, user5.Gender)
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
	}{
		{
			name:             "successful test",
			resultQueryList:  successRows,
			resultQueryError: nil,
			wantList:         users,
			wantErr:          nil,
		},
		{
			name:             "bad test",
			resultQueryList:  nil,
			resultQueryError: errors.New("test"),
			wantList:         []models.User{},
			wantErr:          fmt.Errorf("GetUserList err: %v", errors.New("test")),
		},
		{
			name:             "bad scanning",
			resultQueryList:  badRows,
			resultQueryError: nil,
			wantList:         []models.User{},
			wantErr:          fmt.Errorf("GetUserList err during scanning"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db}
			if tt.resultQueryList != nil {
				mock.ExpectQuery("SELECT id, username, age, gender FROM users").WillReturnRows(tt.resultQueryList)
			} else {
				mock.ExpectQuery("SELECT id, username, age, gender FROM users").WillReturnError(tt.resultQueryError)
			}

			list, err := storage.GetUserList(context.Background())
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

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	user1 := models.User{ID: 1, Username: "kirill", Password: "124", Age: 20, Gender: "male"}

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
			storage := Storage{db}
			mock.ExpectExec("INSERT INTO users").
				WithArgs(tt.user.Username, tt.user.Password, tt.user.Age, tt.user.Gender).
				WillReturnError(tt.queryErr)
			//WillReturnResult(sqlmock.NewResult(1, 1))
			err := storage.AddUser(context.Background(), tt.user)
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
			storage := Storage{db}
			mock.ExpectExec("DELETE FROM users").WithArgs(tt.username).WillReturnError(tt.execResult)
			err := storage.DeleteUser(context.Background(), tt.username)
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

func TestGetUserByUsername(t *testing.T) {
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
			storage := Storage{db}
			if tt.queryResult != nil {
				mock.ExpectQuery("SELECT id, username, password FROM users").WillReturnRows(tt.queryResult)
			} else {
				mock.ExpectQuery("SELECT id, username, password FROM users").WillReturnError(tt.queryError)
			}

			user, err := storage.GetUserByUsername(context.Background(), tt.username)

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
