package profile

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/usecase/profile/mocks"
	"sparkit/internal/utils/consts"
	"testing"
	"time"
)

func TestCreateProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	defer logger.Sync()

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
			profile:                 models.Profile{Age: 20},
			repoCreateProfileResult: 2,
			repoCreateProfileError:  nil,
			repoCreateProfileCount:  1,
			logger:                  logger,
			wantId:                  2,
		},
		{
			name:                    "bad create profile",
			profile:                 models.Profile{Age: 15},
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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().CreateProfile(ctx, tt.profile).Return(tt.repoCreateProfileResult, tt.repoCreateProfileError).
				Times(tt.repoCreateProfileCount)
			s := New(repo, logger)
			id, err := s.CreateProfile(ctx, tt.profile)
			require.ErrorIs(t, err, tt.repoCreateProfileError)
			if id != tt.wantId {
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
	defer logger.Sync()

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
			profile:            models.Profile{Age: 20},
			updateProfileErr:   nil,
			updateProfileCount: 1,
			logger:             logger,
		},
		{
			name:               "bad update profile",
			id:                 1,
			profile:            models.Profile{Age: 15},
			updateProfileErr:   errors.New("failed to update profile with age: 15"),
			updateProfileCount: 1,
			logger:             logger,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().UpdateProfile(ctx, tt.id, tt.profile).Return(tt.updateProfileErr).
				Times(tt.updateProfileCount)
			s := New(repo, logger)
			err := s.UpdateProfile(ctx, tt.id, tt.profile)
			require.ErrorIs(t, err, tt.updateProfileErr)
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	defer logger.Sync()

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
			returnProfile: models.Profile{Age: 20},
			returnError:   nil,
			callCount:     1,
			logger:        logger,
			wantProfile:   models.Profile{Age: 20},
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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetProfile(ctx, tt.id).Return(tt.returnProfile, tt.returnError).
				Times(tt.callCount)
			s := New(repo, logger)

			profile, err := s.GetProfile(ctx, tt.id)

			require.ErrorIs(t, err, tt.returnError)
			if profile != tt.wantProfile {
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
	defer logger.Sync()

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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().DeleteProfile(ctx, tt.id).Return(tt.returnError).
				Times(tt.callCount)

			s := New(repo, logger)
			err := s.DeleteProfile(ctx, tt.id)
			require.ErrorIs(t, err, tt.returnError)
		})
	}
}
