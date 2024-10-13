package user

import (
	"context"
	"github.com/golang/mock/gomock"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/models"
	"sparkit/internal/usecase/user/mocks"
	"sparkit/internal/utils/hashing"
	"testing"
)

func TestRegisterUser(t *testing.T) {

	user1 := models.User{ID: 1}
	user2 := models.User{ID: 2}
	tests := []struct {
		name string
		user models.User
		want error
	}{
		{
			name: "successfull test",
			user: user1,
			want: nil,
		},
		{
			name: "bad test",
			user: user2,
			want: sparkiterrors.ErrRegistrationUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().AddUser(gomock.Any(), tt.user).Return(tt.want).Times(1)
			u := New(repo)
			res := u.RegisterUser(ctx, tt.user)
			if res != tt.want {
				t.Errorf("RegisterUser() = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// func (repo *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error)
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
			wantUser:         models.User{},
			wantErr:          sparkiterrors.ErrWrongCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetUserByUsername(gomock.Any(), tt.user.Username).Return(tt.getUserWant, tt.getUserError).Times(tt.getUserCallCount)
			u := New(repo)
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

//func TestGetUserList(t *testing.T) {
//	// func (repo *Storage) GetUserList(ctx context.Context) ([]models.User, error)
//	user1 := models.User{ID: 1, Username: "Kirill", Password: "123456"}
//	user2 := models.User{ID: 2, Username: "Andrey", Password: "123456"}
//	user3 := models.User{ID: 3, Username: "Kirill", Password: "123456"}
//	userList := []models.User{user1, user2, user3}
//	tests := []struct{
//		name string
//		wantList []models.User
//		wantErr error
//	}{
//		{
//			name: "successfull test",
//			wantList: userList,
//			wantErr: nil,
//		},
//		{
//			name: "bad test",
//			wantList: []models.User{},
//			wantErr: sparkiterrors.ErrWrongCredentials,
//		},
//	}
//}
