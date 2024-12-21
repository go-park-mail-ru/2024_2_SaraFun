package personalitiesgrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcPersonalitiesHandler_GetFeedList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)

	ctx := context.Background()
	req := &generatedPersonalities.GetFeedListRequest{UserID: 10, Receivers: []int32{20, 30}}
	expectedUsers := []models.User{{ID: 1, Username: "user1"}, {ID: 2, Username: "user2"}}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().GetFeedList(ctx, 10, []int{20, 30}).Return(expectedUsers, nil)
		resp, err := handler.GetFeedList(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if len(resp.Users) != 2 {
			t.Errorf("expected 2 users, got %d", len(resp.Users))
		}
	})

	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().GetFeedList(ctx, 10, []int{20, 30}).Return(nil, errors.New("some error"))
		_, err := handler.GetFeedList(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get feed list error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()

	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()
	req := &generatedPersonalities.RegisterUserRequest{User: &generatedPersonalities.User{ID: 1, Username: "u", Email: "e", Password: "p", Profile: 10}}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().RegisterUser(ctx, models.User{ID: 1, Username: "u", Email: "e", Password: "p", Profile: 10}).Return(100, nil)
		resp, err := handler.RegisterUser(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.UserId != 100 {
			t.Errorf("expected userId=100 got %d", resp.UserId)
		}
	})

	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().RegisterUser(ctx, gomock.Any()).Return(0, errors.New("fail"))
		_, err := handler.RegisterUser(ctx, req)
		if err == nil || !contains(err.Error(), "grpc register user error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_CheckPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)

	ctx := context.Background()
	req := &generatedPersonalities.CheckPasswordRequest{Username: "u", Password: "p"}
	t.Run("success", func(t *testing.T) {
		u := models.User{ID: 10, Username: "u", Email: "e", Password: "p", Profile: 20}
		userUC.EXPECT().CheckPassword(ctx, "u", "p").Return(u, nil)
		resp, err := handler.CheckPassword(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.User.ID != 10 {
			t.Errorf("expected id=10 got %d", resp.User.ID)
		}
	})
	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().CheckPassword(ctx, "u", "p").Return(models.User{}, errors.New("invalid"))
		_, err := handler.CheckPassword(ctx, req)
		st, _ := status.FromError(err)
		if st.Code() != codes.InvalidArgument {
			t.Errorf("expected invalid arg got %v", st.Code())
		}
	})
}

func TestGrpcPersonalitiesHandler_GetProfileIDByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)

	ctx := context.Background()
	req := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: 10}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().GetProfileIdByUserId(ctx, 10).Return(101, nil)
		resp, err := handler.GetProfileIDByUserID(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.ProfileID != 101 {
			t.Errorf("expected 101 got %d", resp.ProfileID)
		}
	})
	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().GetProfileIdByUserId(ctx, 10).Return(0, errors.New("fail"))
		_, err := handler.GetProfileIDByUserID(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get profile id by user id error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_GetUsernameByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)

	ctx := context.Background()
	req := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: 20}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().GetUsernameByUserId(ctx, 20).Return("usernameX", nil)
		resp, err := handler.GetUsernameByUserID(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.Username != "usernameX" {
			t.Errorf("expected usernameX got %s", resp.Username)
		}
	})
	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().GetUsernameByUserId(ctx, 20).Return("", errors.New("fail"))
		_, err := handler.GetUsernameByUserID(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get username by user id error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_GetUserIDByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	req := &generatedPersonalities.GetUserIDByUsernameRequest{Username: "testuser"}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().GetUserIdByUsername(ctx, "testuser").Return(33, nil)
		resp, err := handler.GetUserIDByUsername(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.UserID != 33 {
			t.Errorf("expected 33 got %d", resp.UserID)
		}
	})

	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().GetUserIdByUsername(ctx, "testuser").Return(0, errors.New("fail"))
		_, err := handler.GetUserIDByUsername(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get user id by username error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_CheckUsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	req := &generatedPersonalities.CheckUsernameExistsRequest{Username: "usr"}

	t.Run("exists", func(t *testing.T) {
		userUC.EXPECT().CheckUsernameExists(ctx, "usr").Return(true, nil)
		resp, err := handler.CheckUsernameExists(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.Exists != true {
			t.Errorf("expected true got %v", resp.Exists)
		}
	})
	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().CheckUsernameExists(ctx, "usr").Return(false, errors.New("fail"))
		_, err := handler.CheckUsernameExists(ctx, req)
		if err == nil || !contains(err.Error(), "grpc check username exists error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_CreateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	prof := &generatedPersonalities.Profile{ID: 1, FirstName: "F", LastName: "L", Age: 30, Gender: "M", Target: "T", About: "A", BirthDate: "2000-01-01"}
	req := &generatedPersonalities.CreateProfileRequest{Profile: prof}

	t.Run("success", func(t *testing.T) {
		p := models.Profile{ID: 1, FirstName: "F", LastName: "L", Age: 30, Gender: "M", Target: "T", About: "A", BirthdayDate: "2000-01-01"}
		profileUC.EXPECT().CreateProfile(ctx, p).Return(99, nil)
		resp, err := handler.CreateProfile(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.ProfileId != 99 {
			t.Errorf("expected 99 got %d", resp.ProfileId)
		}
	})

	t.Run("error", func(t *testing.T) {
		profileUC.EXPECT().CreateProfile(ctx, gomock.Any()).Return(0, errors.New("fail"))
		_, err := handler.CreateProfile(ctx, req)
		if err == nil || !contains(err.Error(), "grpc create profile error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	inProf := &generatedPersonalities.Profile{ID: 10, FirstName: "F", LastName: "L", Age: 20, Gender: "M", Target: "T", About: "A", BirthDate: "1990-01-01"}
	req := &generatedPersonalities.UpdateProfileRequest{Profile: inProf}

	t.Run("success", func(t *testing.T) {
		p := models.Profile{ID: 10, FirstName: "F", LastName: "L", Age: 20, Gender: "M", Target: "T", About: "A", BirthdayDate: "1990-01-01"}
		profileUC.EXPECT().UpdateProfile(ctx, 10, p).Return(nil)
		_, err := handler.UpdateProfile(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		profileUC.EXPECT().UpdateProfile(ctx, 10, gomock.Any()).Return(errors.New("fail"))
		_, err := handler.UpdateProfile(ctx, req)
		if err == nil || !contains(err.Error(), "grpc update profile error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	req := &generatedPersonalities.GetProfileRequest{Id: 50}

	t.Run("success", func(t *testing.T) {
		p := models.Profile{ID: 50, FirstName: "F"}
		profileUC.EXPECT().GetProfile(ctx, 50).Return(p, nil)
		resp, err := handler.GetProfile(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		if resp.Profile.ID != 50 {
			t.Errorf("expected 50 got %d", resp.Profile.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		profileUC.EXPECT().GetProfile(ctx, 50).Return(models.Profile{}, errors.New("fail"))
		_, err := handler.GetProfile(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get profile error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_DeleteProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	req := &generatedPersonalities.DeleteProfileRequest{Id: 77}

	t.Run("success", func(t *testing.T) {
		profileUC.EXPECT().DeleteProfile(ctx, 77).Return(nil)
		_, err := handler.DeleteProfile(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		profileUC.EXPECT().DeleteProfile(ctx, 77).Return(errors.New("fail"))
		_, err := handler.DeleteProfile(ctx, req)
		if err == nil || !contains(err.Error(), "grpc delete profile error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGrpcPersonalitiesHandler_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mocks.NewMockUserUsecase(ctrl)
	profileUC := mocks.NewMockProfileUsecase(ctrl)
	logger := zap.NewNop()
	handler := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
	ctx := context.Background()

	req := &generatedPersonalities.ChangePasswordRequest{UserID: 99, Password: "newpass"}

	t.Run("success", func(t *testing.T) {
		userUC.EXPECT().ChangePassword(ctx, 99, "newpass").Return(nil)
		_, err := handler.ChangePassword(ctx, req)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		userUC.EXPECT().ChangePassword(ctx, 99, "newpass").Return(errors.New("fail"))
		_, err := handler.ChangePassword(ctx, req)
		if err == nil || !contains(err.Error(), "Grpc change password error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}
func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
