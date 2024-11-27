package repo

import (
	"bytes"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
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
