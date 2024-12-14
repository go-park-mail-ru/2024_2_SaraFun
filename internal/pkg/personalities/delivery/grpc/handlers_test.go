package personalitiesgrpc

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCreateProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	//defer logger.Sync()

	tests := []struct {
		name                    string
		profile                 models.Profile
		repoCreateProfileResult int
		repoCreateProfileError  error
		repoCreateProfileCount  int
		logger                  *zap.Logger
		wantId                  int
	}{
		{
			name:                    "succesful create profile",
			profile:                 models.Profile{Age: 24, BirthdayDate: "2000-01-01"},
			repoCreateProfileResult: 2,
			repoCreateProfileError:  nil,
			repoCreateProfileCount:  1,
			logger:                  logger,
			wantId:                  2,
		},
		{
			name:                    "bad create profile",
			profile:                 models.Profile{Age: 14, BirthdayDate: "2010-01-01"},
			repoCreateProfileResult: 0,
			repoCreateProfileError:  errors.New("failed to create profile with age: 15"),
			repoCreateProfileCount:  1,
			logger:                  logger,
			wantId:                  0,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUC := mocks.NewMockUserUsecase(mockCtrl)
			profileUC := mocks.NewMockProfileUsecase(mockCtrl)
			profileUC.EXPECT().CreateProfile(ctx, tt.profile).Return(tt.repoCreateProfileResult, tt.repoCreateProfileError).
				Times(tt.repoCreateProfileCount)
			s := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
			profile := &generatedPersonalities.Profile{
				ID:        int32(tt.profile.ID),
				FirstName: tt.profile.FirstName,
				LastName:  tt.profile.LastName,
				Age:       int32(tt.profile.Age),
				Gender:    tt.profile.Gender,
				Target:    tt.profile.Target,
				About:     tt.profile.About,
				BirthDate: tt.profile.BirthdayDate,
			}
			req := &generatedPersonalities.CreateProfileRequest{Profile: profile}
			id, err := s.CreateProfile(ctx, req)
			if err != nil {
				id = &generatedPersonalities.CreateProfileResponse{}
			}
			require.ErrorIs(t, err, tt.repoCreateProfileError)
			if int(id.ProfileId) != tt.wantId {
				t.Errorf("CreateProfile() id = %v, want %v", id, tt.wantId)
			}
		})
	}

}

func TestUpdateProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	//defer logger.Sync()

	tests := []struct {
		name               string
		id                 int
		profile            models.Profile
		updateProfileErr   error
		updateProfileCount int
		logger             *zap.Logger
	}{
		{
			name:               "succesful update profile",
			id:                 1,
			profile:            models.Profile{ID: 1, Age: 24, BirthdayDate: "2000-01-01"},
			updateProfileErr:   nil,
			updateProfileCount: 1,
			logger:             logger,
		},
		{
			name:               "bad update profile",
			id:                 1,
			profile:            models.Profile{ID: 1, Age: 14, BirthdayDate: "2010-01-01"},
			updateProfileErr:   errors.New("failed to update profile with age: 15"),
			updateProfileCount: 1,
			logger:             logger,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUC := mocks.NewMockUserUsecase(mockCtrl)
			profileUC := mocks.NewMockProfileUsecase(mockCtrl)
			profileUC.EXPECT().UpdateProfile(ctx, tt.id, tt.profile).Return(tt.updateProfileErr).
				Times(tt.updateProfileCount)
			s := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
			profile := &generatedPersonalities.Profile{
				ID:        int32(tt.id),
				FirstName: tt.profile.FirstName,
				LastName:  tt.profile.LastName,
				Age:       int32(tt.profile.Age),
				Gender:    tt.profile.Gender,
				Target:    tt.profile.Target,
				About:     tt.profile.About,
				BirthDate: tt.profile.BirthdayDate,
			}
			req := &generatedPersonalities.UpdateProfileRequest{
				Id:      int32(tt.id),
				Profile: profile,
			}
			_, err := s.UpdateProfile(ctx, req)
			require.ErrorIs(t, err, tt.updateProfileErr)
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	//defer logger.Sync()

	tests := []struct {
		name          string
		id            int
		returnProfile models.Profile
		returnError   error
		callCount     int
		logger        *zap.Logger
		wantProfile   models.Profile
	}{
		{
			name:          "successfull get profile",
			id:            1,
			returnProfile: models.Profile{Age: 24, BirthdayDate: "2000-01-01"},
			returnError:   nil,
			callCount:     1,
			logger:        logger,
			wantProfile:   models.Profile{Age: 24, BirthdayDate: "2000-01-01"},
		},
		{
			name:          "bad get profile",
			id:            2,
			returnProfile: models.Profile{},
			returnError:   errors.New("failed to get profile"),
			callCount:     1,
			logger:        logger,
			wantProfile:   models.Profile{},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUC := mocks.NewMockUserUsecase(mockCtrl)
			profileUC := mocks.NewMockProfileUsecase(mockCtrl)
			profileUC.EXPECT().GetProfile(ctx, tt.id).Return(tt.returnProfile, tt.returnError).
				Times(tt.callCount)
			s := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)

			req := &generatedPersonalities.GetProfileRequest{Id: int32(tt.id)}
			profile, err := s.GetProfile(ctx, req)
			if err != nil {
				profile = &generatedPersonalities.GetProfileResponse{Profile: &generatedPersonalities.Profile{}}
			}
			checkProfile := models.Profile{
				ID:           int(profile.Profile.ID),
				FirstName:    profile.Profile.FirstName,
				LastName:     profile.Profile.LastName,
				Age:          int(profile.Profile.Age),
				BirthdayDate: profile.Profile.BirthDate,
				Gender:       profile.Profile.Gender,
				Target:       profile.Profile.Target,
				About:        profile.Profile.About,
			}

			require.ErrorIs(t, err, tt.returnError)
			if checkProfile != tt.wantProfile {
				t.Errorf("GetProfile() profile = %v, want %v", profile, tt.wantProfile)
			}
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	//defer logger.Sync()

	tests := []struct {
		name        string
		id          int
		returnError error
		callCount   int
		logger      *zap.Logger
	}{
		{
			name:        "good delete profile",
			id:          1,
			returnError: nil,
			callCount:   1,
			logger:      logger,
		},
		{
			name:        "bad delete profile",
			id:          2,
			returnError: errors.New("failed to delete profile"),
			callCount:   1,
			logger:      logger,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//repo := mocks.NewMockRepository(mockCtrl)
			userUC := mocks.NewMockUserUsecase(mockCtrl)
			profileUC := mocks.NewMockProfileUsecase(mockCtrl)
			profileUC.EXPECT().DeleteProfile(ctx, tt.id).Return(tt.returnError).
				Times(tt.callCount)

			s := NewGrpcPersonalitiesHandler(userUC, profileUC, logger)
			req := &generatedPersonalities.DeleteProfileRequest{Id: int32(tt.id)}
			_, err := s.DeleteProfile(ctx, req)
			require.ErrorIs(t, err, tt.returnError)
		})
	}
}
