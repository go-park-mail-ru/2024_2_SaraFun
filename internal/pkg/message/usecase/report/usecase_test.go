package report

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/usecase/report/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddReport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)

	tests := []struct {
		name         string
		report       models.Report
		repoReportID int
		repoError    error
		repoCount    int
		expectedID   int
	}{
		{
			name: "successfull test",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "abuse",
				Body:     "Он меня обозвал",
			},
			repoReportID: 1,
			repoError:    nil,
			repoCount:    1,
			expectedID:   1,
		},
		{
			name: "bad test",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "",
				Body:     "",
			},
			repoReportID: -1,
			repoError:    errors.New("errors"),
			repoCount:    1,
			expectedID:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().AddReport(ctx, tt.report).Return(tt.repoReportID, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			id, err := usecase.AddReport(ctx, tt.report)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestGetReportIfExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		firstUserID    int
		secondUserID   int
		repoReport     models.Report
		repoError      error
		repoCount      int
		expectedReport models.Report
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			repoReport: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "abuse",
				Body:     "он обзывается",
			},
			repoError: nil,
			repoCount: 1,
			expectedReport: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "abuse",
				Body:     "он обзывается",
			},
		},
		{
			name:           "bad test",
			firstUserID:    1,
			secondUserID:   2,
			repoReport:     models.Report{},
			repoError:      errors.New("error"),
			repoCount:      1,
			expectedReport: models.Report{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetReportIfExists(ctx, tt.firstUserID, tt.secondUserID).Return(tt.repoReport, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			rpeort, err := usecase.GetReportIfExists(ctx, tt.firstUserID, tt.secondUserID)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedReport, rpeort)
		})
	}
}

func TestCheckUsersBlockNotExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		firstUserID    int
		secondUserID   int
		repoReport     models.Report
		repoError      error
		repoCount      int
		expectedString string
		expectedError  error
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			repoReport: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "abuse",
				Body:     "он злой",
			},
			repoError:      nil,
			repoCount:      1,
			expectedString: "Вы заблокировали данного пользователя",
			expectedError:  errors.New("block exists"),
		},
		{
			name:           "bad test",
			firstUserID:    1,
			secondUserID:   2,
			repoReport:     models.Report{},
			repoError:      errors.New("error"),
			repoCount:      1,
			expectedString: "",
			expectedError:  fmt.Errorf("GetReport error: %w", errors.New("error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetReportIfExists(ctx, tt.firstUserID, tt.secondUserID).Return(tt.repoReport, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			result, err := usecase.CheckUsersBlockNotExists(ctx, tt.firstUserID, tt.secondUserID)
			require.Error(t, err, tt.expectedError)
			require.Equal(t, tt.expectedString, result)
		})
	}
}
