package repo

import (
	"bytes"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"os"
	"testing"
	"time"
)

type mockFile struct {
	*bytes.Reader
}

func (f *mockFile) Close() error {
	return nil
}

func TestSaveImage(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")

	// Создаем содержимое файла
	fileContent := []byte("test file content")
	testFile := &mockFile{bytes.NewReader(fileContent)}

	// Создаем необходимую директорию
	err := os.MkdirAll("C:/home/reufee/imagedata/", os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	defer os.RemoveAll("C:/home/reufee/") // Удаляем директорию после теста

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	successRow := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	tests := []struct {
		name        string
		file        *mockFile
		fileExt     string
		userId      int
		queryErr    error
		queryResult *sqlmock.Rows
		wantId      int
		wantErr     error
	}{
		{
			name:        "successful test",
			file:        testFile,
			fileExt:     "png",
			userId:      1,
			queryErr:    nil,
			queryResult: successRow,
			wantId:      1,
			wantErr:     nil,
		},
		{
			name:     "error test",
			file:     testFile,
			fileExt:  "txt",
			userId:   1,
			queryErr: errors.New("test error"),
			wantId:   -1,
			wantErr:  errors.New("saveImage err: test error"),
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			if tt.queryErr != nil {
				mock.ExpectQuery("INSERT INTO photo").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("INSERT INTO photo").WillReturnRows(tt.queryResult)
			}

			// Сбрасываем указатель файла
			tt.file.Seek(0, io.SeekStart)
			id, err := storage.SaveImage(ctx, tt.file, tt.fileExt, tt.userId)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			if id != tt.wantId {
				t.Errorf("SaveImage() id = %v, want %v", id, tt.wantId)
			}
		})
	}
}
