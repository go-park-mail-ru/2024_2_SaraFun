package message

import (
	"context"
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

func TestAddMessage(t *testing.T) {
	// Создаём мок-драйвер базы данных
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
		message     *models.Message
		mockSetup   func()
		wantID      int
		wantErr     error
		ctx         context.Context
		expectPanic bool
	}{
		{
			name: "Successful insert",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "Hello, World!",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO message").
					WithArgs(1, 2, "Hello, World!").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
			},
			wantID:      100,
			wantErr:     nil,
			ctx:         ctx,
			expectPanic: false,
		},
		{
			name: "Database error on insert",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "Hello, World!",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO message").
					WithArgs(1, 2, "Hello, World!").
					WillReturnError(errors.New("insert failed"))
			},
			wantID:      -1,
			wantErr:     fmt.Errorf("AddMessage error: insert failed"),
			ctx:         ctx,
			expectPanic: false,
		},
		{
			name: "Missing RequestIDKey in context",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "Hello, World!",
			},
			mockSetup: func() {

			},
			wantID:      -1,
			wantErr:     fmt.Errorf("invalid or missing request ID"),
			ctx:         context.Background(),
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной цикла
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			if tt.expectPanic {
				require.Panics(t, func() {
					storage.AddMessage(tt.ctx, tt.message)
				}, "The code did not panic as expected")
			} else {
				gotID, err := storage.AddMessage(tt.ctx, tt.message)

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
