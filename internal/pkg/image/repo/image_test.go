package repo

import (
	"bytes"
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

type mockFile struct {
	*bytes.Reader
}

func (f *mockFile) Close() error {
	return nil
}

func TestDeleteImage(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "req-12345")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		imageId  int
		queryErr error
		wantErr  error
	}{
		{
			name:     "successful deletion",
			imageId:  1,
			queryErr: nil,
			wantErr:  nil,
		},
		{
			name:     "query error",
			imageId:  2,
			queryErr: errors.New("delete failed"),
			wantErr:  errors.New("deleteImage err: delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}

			if tt.queryErr != nil {
				mock.ExpectExec("DELETE FROM photo").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectExec("DELETE FROM photo").WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := storage.DeleteImage(ctx, tt.imageId)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetImageLinksByUserId(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "req-12345")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name        string
		userId      int
		mockRows    *sqlmock.Rows
		queryErr    error
		expected    []models.Image
		expectedErr error
	}{
		{
			name:   "successful fetch image links",
			userId: 1,
			mockRows: sqlmock.NewRows([]string{"id", "link", "number"}).
				AddRow(1, "/path/to/image1.png", 1).
				AddRow(2, "/path/to/image2.png", 2),
			queryErr: nil,
			expected: []models.Image{
				{Id: 1, Link: "/path/to/image1.png", Number: 1},
				{Id: 2, Link: "/path/to/image2.png", Number: 2},
			},
			expectedErr: nil,
		},
		{
			name:        "database error on fetch image links",
			userId:      1,
			mockRows:    nil,
			queryErr:    errors.New("database error"),
			expected:    nil,
			expectedErr: errors.New("GetImageLink err: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := New(db, logger)

			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT id, link, number FROM photo").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT id, link, number FROM photo").
					WithArgs(tt.userId).
					WillReturnRows(tt.mockRows)
			}

			links, err := storage.GetImageLinksByUserId(ctx, tt.userId)

			require.Equal(t, tt.expected, links)
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateOrdNumbers(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "req-12345")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name        string
		numbers     []models.Image
		mockQuery   func()
		expectedErr error
	}{
		{
			name: "successful update order numbers",
			numbers: []models.Image{
				{Id: 1, Number: 2},
				{Id: 2, Number: 1},
			},
			mockQuery: func() {
				mock.ExpectExec("UPDATE photo SET number").
					WithArgs(2, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE photo SET number").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "database error during update",
			numbers: []models.Image{
				{Id: 1, Number: 2},
			},
			mockQuery: func() {
				mock.ExpectExec("UPDATE photo SET number").
					WithArgs(2, 1).
					WillReturnError(errors.New("database error"))
			},
			expectedErr: errors.New("updateOrdNumbers err: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := New(db, logger)
			tt.mockQuery()

			err := storage.UpdateOrdNumbers(ctx, tt.numbers)

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
