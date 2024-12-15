package usecase

//
//import (
//	"context"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/usecase/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/zap"
//	"testing"
//	"time"
//)
//
//func TestAddBalance(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().AddBalance(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.AddBalance(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestAddDailyLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().AddDailyLikeCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.AddDailyLikesCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestAddPurchasedLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().AddPurchasedLikeCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.AddPurchasedLikesCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestChangeBalance(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().ChangeBalance(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.ChangeBalance(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestChangeDailyLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().ChangeDailyLikeCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.ChangeDailyLikeCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestChangePurchasedLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().ChangePurchasedLikeCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.ChangePurchasedLikeCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestSetBalance(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().SetBalance(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.SetBalance(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestSetDailyLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().SetDailyLikesCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.SetDailyLikeCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestSetDailyLikesCountToAll(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().SetDailyLikesCountToAll(ctx, tt.amount).Return(tt.repoErr)
//			err := usecase.SetDailyLikeCountToAll(ctx, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestSetPurchasedLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name    string
//		userID  int
//		amount  int
//		repoErr error
//	}{
//		{
//			name:    "good test",
//			userID:  1,
//			amount:  10,
//			repoErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().SetPurchasedLikesCount(ctx, tt.userID, tt.amount).Return(tt.repoErr)
//			err := usecase.SetPurchasedLikeCount(ctx, tt.userID, tt.amount)
//			require.ErrorIs(t, err, tt.repoErr)
//		})
//	}
//}
//
//func TestGetBalance(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name           string
//		userID         int
//		repoReturn     int
//		repoErr        error
//		repoTimes      int
//		expectedAmount int
//	}{
//		{
//			name:           "good test",
//			userID:         1,
//			repoReturn:     10,
//			repoErr:        nil,
//			repoTimes:      1,
//			expectedAmount: 10,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().GetBalance(ctx, tt.userID).Return(tt.repoReturn, tt.repoErr).Times(tt.repoTimes)
//			amount, err := usecase.GetBalance(ctx, tt.userID)
//			require.ErrorIs(t, err, tt.repoErr)
//			require.Equal(t, tt.expectedAmount, amount)
//		})
//	}
//}
//
//func TestGetDailyLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name           string
//		userID         int
//		repoReturn     int
//		repoErr        error
//		repoTimes      int
//		expectedAmount int
//	}{
//		{
//			name:           "good test",
//			userID:         1,
//			repoReturn:     10,
//			repoErr:        nil,
//			repoTimes:      1,
//			expectedAmount: 10,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().GetDailyLikesCount(ctx, tt.userID).Return(tt.repoReturn, tt.repoErr).Times(tt.repoTimes)
//			amount, err := usecase.GetDailyLikesCount(ctx, tt.userID)
//			require.ErrorIs(t, err, tt.repoErr)
//			require.Equal(t, tt.expectedAmount, amount)
//		})
//	}
//}
//
//func TestGetPurchasedLikesCount(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	repo := mocks.NewMockRepository(mockCtrl)
//	usecase := New(repo, logger)
//
//	tests := []struct {
//		name           string
//		userID         int
//		repoReturn     int
//		repoErr        error
//		repoTimes      int
//		expectedAmount int
//	}{
//		{
//			name:           "good test",
//			userID:         1,
//			repoReturn:     10,
//			repoErr:        nil,
//			repoTimes:      1,
//			expectedAmount: 10,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo.EXPECT().GetPurchasedLikesCount(ctx, tt.userID).Return(tt.repoReturn, tt.repoErr).Times(tt.repoTimes)
//			amount, err := usecase.GetPurchasedLikesCount(ctx, tt.userID)
//			require.ErrorIs(t, err, tt.repoErr)
//			require.Equal(t, tt.expectedAmount, amount)
//		})
//	}
//}
