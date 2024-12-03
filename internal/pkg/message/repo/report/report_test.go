package report

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddReport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() err: %v", err)
	}
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	tests := []struct {
		name       string
		report     models.Report
		queryID    int
		queryError error
		expectedID int
	}{
		{
			name: "successfull test",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "test",
				Body:     "success",
			},
			queryID:    1,
			queryError: nil,
			expectedID: 1,
		},
		{
			name: "bad test",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "test",
				Body:     "bad",
			},
			queryID:    0,
			queryError: errors.New("error"),
			expectedID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("INSERT INTO report").
					WithArgs(tt.report.Author, tt.report.Receiver, tt.report.Reason, tt.report.Body).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.queryID))
			} else {
				mock.ExpectQuery("INSERT INTO report").
					WithArgs(tt.report.Author, tt.report.Receiver, tt.report.Reason, tt.report.Body).
					WillReturnError(tt.queryError)

			}
			id, err := repo.AddReport(ctx, tt.report)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedID, id)
		})
	}

}

func TestGetReportIfExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() err: %v", err)
	}
	defer db.Close()
	logger := zap.NewNop()
	repo := New(db, logger)
	tests := []struct {
		name           string
		firstUserID    int
		secondUserID   int
		queryRows      *sqlmock.Rows
		queryError     error
		expectedReport models.Report
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			queryRows: sqlmock.NewRows([]string{"count", "author", "receiver", "body"}).
				AddRow(1, 1, 2, "success"),
			queryError: nil,
			expectedReport: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "success",
			},
		},
		{
			name:         "bad test",
			firstUserID:  1,
			secondUserID: 2,
			queryRows: sqlmock.NewRows([]string{"count", "author", "receiver", "body"}).
				AddRow(1, 1, 2, "bad"),
			queryError:     errors.New("error"),
			expectedReport: models.Report{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.firstUserID, tt.secondUserID).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT").WithArgs(tt.firstUserID, tt.secondUserID).
					WillReturnError(tt.queryError)
			}

			report, err := repo.GetReportIfExists(ctx, tt.firstUserID, tt.secondUserID)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedReport, report)
		})
	}
}
