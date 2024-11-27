package report

import (
	"context"
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddReport(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()

	storage := New(db, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		report      models.Report
		mockSetup   func()
		wantID      int
		wantErr     error
		ctx         context.Context
		expectPanic bool
	}{
		{
			name: "Successful insert",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "Inappropriate content",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO report").
					WithArgs(1, 2, "Inappropriate content").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
			},
			wantID:      100,
			wantErr:     nil,
			ctx:         ctx,
			expectPanic: false,
		},
		{
			name: "Database error on insert",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "Spam",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO report").
					WithArgs(1, 2, "Spam").
					WillReturnError(errors.New("insert failed"))
			},
			wantID:      -1,
			wantErr:     fmt.Errorf("AddReport insert report: insert failed"),
			ctx:         ctx,
			expectPanic: false,
		},
		{
			name: "Missing RequestIDKey in context",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "Harassment",
			},
			mockSetup: func() {
				// Мок не настраивается, так как ошибка происходит до обращения к базе
			},
			wantID:      -1,
			wantErr:     fmt.Errorf("invalid or missing request ID"),
			ctx:         context.Background(),
			expectPanic: true, // Ожидается панику
		},
		{
			name: "Invalid RequestIDKey type in context",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "Harassment",
			},
			mockSetup: func() {
				// Мок не настраивается, так как ошибка происходит до обращения к базе
			},
			wantID:      -1,
			wantErr:     fmt.Errorf("invalid or missing request ID"),
			ctx:         context.WithValue(context.Background(), consts.RequestIDKey, 12345), // Неправильный тип
			expectPanic: true,                                                                // Ожидается панику
		},
		{
			name: "Empty Body in report",
			report: models.Report{
				Author:   1,
				Receiver: 2,
				Body:     "",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO report").
					WithArgs(1, 2, "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(101))
			},
			wantID:      101,
			wantErr:     nil,
			ctx:         ctx,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			if tt.expectPanic {
				require.Panics(t, func() {
					storage.AddReport(tt.ctx, tt.report)
				}, "The code did not panic as expected")
			} else {
				gotID, err := storage.AddReport(tt.ctx, tt.report)

				if tt.wantErr != nil {
					require.Error(t, err)
					require.EqualError(t, err, tt.wantErr.Error())
					require.Equal(t, tt.wantID, gotID)
				} else {
					require.NoError(t, err)
					require.Equal(t, tt.wantID, gotID)
				}

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}
		})
	}
}
