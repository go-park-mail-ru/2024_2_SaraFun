package repo

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"os/user"
	"testing"
	"time"
)

//SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int64, error)
//
//GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
//
//DeleteImage(ctx context.Context, id int) error

func TestSaveImage(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	testFile, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	os_user, err := user.Current()
	if err != nil {
		t.Error(err)
	}
	os.Setenv("OS_USER", os_user.Username)
	successRow := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	//badRows := sqlmock.NewRows([]string{"random"}).
	//	AddRow(1)

	tests := []struct {
		name        string
		file        multipart.File
		fileExt     string
		userId      int
		ordNumber   int
		queryErr    error
		queryResult *sqlmock.Rows
		wantId      int
		wantErr     error
	}{
		{
			name:        "successfull test",
			file:        testFile,
			fileExt:     "png",
			userId:      1,
			ordNumber:   1,
			queryErr:    nil,
			queryResult: successRow,
			wantId:      1,
			wantErr:     nil,
		},
		{
			name:     "bad test",
			file:     testFile,
			fileExt:  "txt",
			userId:   1,
			queryErr: errors.New("test error"),
			wantId:   -1,
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

			id, err := storage.SaveImage(ctx, tt.file, tt.fileExt, tt.userId, tt.ordNumber)
			require.ErrorIs(t, err, tt.queryErr)
			if id != tt.wantId {
				t.Errorf("SaveImage() id = %v, want %v", id, tt.wantId)
			}
		})
	}
}

func TestGetImageLinksByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	good_rows := sqlmock.NewRows([]string{"id", "link", "number"}).
		AddRow(1, "link1", 1).
		AddRow(2, "link2", 2)

	tests := []struct {
		name        string
		userId      int
		queryErr    error
		queryResult *sqlmock.Rows
		wantImages  []models.Image
	}{
		{
			name:        "successfull test",
			userId:      1,
			queryErr:    nil,
			queryResult: good_rows,
			wantImages:  []models.Image{{Id: 1, Link: "link1", Number: 1}, {Id: 2, Link: "link2", Number: 2}},
		},
		{
			name:        "bad test",
			userId:      1,
			queryErr:    errors.New("test error"),
			queryResult: nil,
			wantImages:  []models.Image{},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{db, logger}
			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT").WillReturnRows(tt.queryResult)
			}

			images, err := storage.GetImageLinksByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.queryErr)
			for i, image := range images {
				if image != tt.wantImages[i] {
					t.Errorf("GetImageLinksByUserId() images[%d] = %v, want %v", i, image, tt.wantImages[i])
				}
			}
		})
	}
}

func TestDeleteImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
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
			name:     "successful delete test",
			imageId:  1,
			queryErr: nil,
			wantErr:  nil,
		},
		{
			name:     "error delete test",
			imageId:  1,
			queryErr: errors.New("delete error"),
			wantErr:  errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := New(db, logger)
			if tt.queryErr != nil {

				mock.ExpectExec("DELETE FROM photo WHERE id = \\$1").WithArgs(tt.imageId).WillReturnError(tt.queryErr)
			} else {
				mock.ExpectExec("DELETE FROM photo WHERE id = \\$1").WithArgs(tt.imageId).WillReturnResult(sqlmock.NewResult(1, 1))
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
