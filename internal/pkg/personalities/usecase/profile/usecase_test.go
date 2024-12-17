package profile

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/usecase/profile/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestUseCase_CreateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()

	tests := []struct {
		name          string
		profile       models.Profile
		mockSetup     func()
		wantErr       bool
		wantErrMsg    string
		wantProfileID int
	}{
		{
			name:    "success",
			profile: models.Profile{BirthdayDate: "2000-01-01"},
			mockSetup: func() {
				repo.EXPECT().CreateProfile(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, p models.Profile) (int, error) {
					if p.Age != time.Now().Year()-2000 {
						return 0, fmt.Errorf("age not calculated correctly")
					}
					return 1, nil
				})
			},
			wantErr:       false,
			wantProfileID: 1,
		},
		{
			name:       "birth date format error",
			profile:    models.Profile{BirthdayDate: "invalid-date"},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: "get age error:",
		},
		{
			name:       "small age error",
			profile:    models.Profile{BirthdayDate: time.Now().AddDate(-17, 0, 0).Format("2006-01-02")},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: sparkiterrors.ErrSmallAge.Error(),
		},
		{
			name:       "big age error",
			profile:    models.Profile{BirthdayDate: "1900-01-01"},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: sparkiterrors.ErrBigAge.Error(),
		},
		{
			name:    "repo error",
			profile: models.Profile{BirthdayDate: "1990-01-01"},
			mockSetup: func() {
				repo.EXPECT().CreateProfile(ctx, gomock.Any()).Return(0, errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "create profile err: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			id, err := uc.CreateProfile(ctx, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.wantErrMsg != "" {
				if err.Error() == tt.wantErrMsg {
					// exact match
				} else if !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error message mismatch: got %v, want contains %v", err.Error(), tt.wantErrMsg)
				}
			}
			if !tt.wantErr && id != tt.wantProfileID {
				t.Errorf("id mismatch: got %v, want %v", id, tt.wantProfileID)
			}
		})
	}
}

func TestUseCase_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()

	tests := []struct {
		name       string
		id         int
		profile    models.Profile
		mockSetup  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "success",
			id:      1,
			profile: models.Profile{BirthdayDate: "2000-01-01"},
			mockSetup: func() {
				repo.EXPECT().UpdateProfile(ctx, 1, gomock.Any()).Return(nil)
			},
		},
		{
			name:       "birth date format error",
			id:         1,
			profile:    models.Profile{BirthdayDate: "invalid"},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: "get age error:",
		},
		{
			name: "small age error",
			id:   1,
			profile: models.Profile{
				BirthdayDate: time.Now().AddDate(-17, 0, 0).Format("2006-01-02"),
			},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: sparkiterrors.ErrSmallAge.Error(),
		},
		{
			name: "big age error",
			id:   1,
			profile: models.Profile{
				BirthdayDate: "1900-01-01",
			},
			mockSetup:  func() {},
			wantErr:    true,
			wantErrMsg: sparkiterrors.ErrBigAge.Error(),
		},
		{
			name:    "repo error",
			id:      1,
			profile: models.Profile{BirthdayDate: "1990-01-01"},
			mockSetup: func() {
				repo.EXPECT().UpdateProfile(ctx, 1, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "update profile err: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := uc.UpdateProfile(ctx, tt.id, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message mismatch: got %v, want contains %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestUseCase_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()

	tests := []struct {
		name       string
		id         int
		mockSetup  func()
		wantErr    bool
		wantErrMsg string
		wantAge    int
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				repo.EXPECT().GetProfile(ctx, 1).Return(models.Profile{BirthdayDate: "1990-01-01"}, nil)
			},
			wantAge: time.Now().Year() - 1990,
		},
		{
			name: "repo error",
			id:   1,
			mockSetup: func() {
				repo.EXPECT().GetProfile(ctx, 1).Return(models.Profile{}, errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "get profile err: db error",
		},
		{
			name: "birth date format error",
			id:   1,
			mockSetup: func() {
				repo.EXPECT().GetProfile(ctx, 1).Return(models.Profile{BirthdayDate: "invalid"}, nil)
			},
			wantErr:    true,
			wantErrMsg: "get profile err:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			res, err := uc.GetProfile(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message mismatch: got %v, want contains %v", err.Error(), tt.wantErrMsg)
			}
			if !tt.wantErr {
				if res.Age != tt.wantAge {
					t.Errorf("age mismatch: got %v, want %v", res.Age, tt.wantAge)
				}
			}
		})
	}
}

func TestUseCase_DeleteProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()

	tests := []struct {
		name       string
		id         int
		mockSetup  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				repo.EXPECT().DeleteProfile(ctx, 1).Return(nil)
			},
		},
		{
			name: "error",
			id:   1,
			mockSetup: func() {
				repo.EXPECT().DeleteProfile(ctx, 1).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "delete profile err: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := uc.DeleteProfile(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("error message mismatch: got %v, want %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestGetAge(t *testing.T) {
	tests := []struct {
		name       string
		birthday   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "valid date",
			birthday: "2000-01-01",
		},
		{
			name:       "invalid format",
			birthday:   "invalid",
			wantErr:    true,
			wantErrMsg: "birth date format error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age, err := GetAge(tt.birthday)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message mismatch: got %v, want contains %v", err.Error(), tt.wantErrMsg)
			}
			if !tt.wantErr {
				expectedAge := time.Now().Year() - 2000
				if time.Now().YearDay() < time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC).YearDay() {
					expectedAge--
				}
				if age != expectedAge {
					t.Errorf("age mismatch: got %v, want %v", age, expectedAge)
				}
			}
		})
	}
}

func TestCheckAge(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		wantErr error
	}{
		{
			name:    "age ok",
			age:     30,
			wantErr: nil,
		},
		{
			name:    "small age",
			age:     17,
			wantErr: sparkiterrors.ErrSmallAge,
		},
		{
			name:    "big age",
			age:     121,
			wantErr: sparkiterrors.ErrBigAge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkAge(tt.age)
			if err != tt.wantErr {
				t.Errorf("error mismatch: got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || (len(str) > len(substr) && (searchSubstring(str, substr))))
}

func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
