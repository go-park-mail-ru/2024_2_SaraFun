package report

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStorage_AddReport(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		report      models.Report
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedID  int
		expectedErr error
	}{
		{
			name: "Successful AddReport",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "Spam",
				Body:     "This is a spam report.",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO report \(author, receiver, reason, body\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(1, 2, "Spam", "This is a spam report.").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Database Error",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Reason:   "Spam",
				Body:     "This is a spam report.",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO report \(author, receiver, reason, body\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(1, 2, "Spam", "This is a spam report.").
					WillReturnError(errors.New("database error"))
			},
			expectedID:  -1,
			expectedErr: errors.New("AddReport insert report: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			id, err := repo.AddReport(context.Background(), tt.report)

			require.Equal(t, tt.expectedID, id)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetReportIfExists(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		firstUserID    int
		secondUserID   int
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedReport models.Report
		expectedErr    error
	}{
		{
			name:         "Report Exists",
			firstUserID:  1,
			secondUserID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\), author, receiver, body FROM report`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count", "author", "receiver", "body"}).
						AddRow(1, 1, 2, "This is a spam report."))
			},
			expectedReport: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "This is a spam report.",
			},
			expectedErr: nil,
		},
		{
			name:         "Report Does Not Exist",
			firstUserID:  1,
			secondUserID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\), author, receiver, body FROM report`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count", "author", "receiver", "body"}))
			},
			expectedReport: models.Report{},
			expectedErr:    errors.New("this report dont exists"),
		},
		{
			name:         "Database Error",
			firstUserID:  1,
			secondUserID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\), author, receiver, body FROM report`).
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expectedReport: models.Report{},
			expectedErr:    errors.New("CheckReportExists select report: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			report, err := repo.GetReportIfExists(context.Background(), tt.firstUserID, tt.secondUserID)

			require.Equal(t, tt.expectedReport, report)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
