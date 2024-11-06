package image

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"sparkit/internal/models"
	"testing"
)

//SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int64, error)
//
//GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
//
//DeleteImage(ctx context.Context, id int) error

func TestSaveImage(t *testing.T) {
	logger := zap.NewNop()
	testFile, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	successRow := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	//badRows := sqlmock.NewRows([]string{"random"}).
	//	AddRow(1)

	tests := []struct {
		name        string
		file        multipart.File
		fileExt     string
		userId      int
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

			id, err := storage.SaveImage(context.Background(), tt.file, tt.fileExt, tt.userId)
			require.ErrorIs(t, err, tt.queryErr)
			if id != tt.wantId {
				t.Errorf("SaveImage() id = %v, want %v", id, tt.wantId)
			}
		})
	}
}

func TestGetImageLinksByUserId(t *testing.T) {
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	good_rows := sqlmock.NewRows([]string{"id", "link"}).
		AddRow(1, "link1").
		AddRow(2, "link2")

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
			wantImages:  []models.Image{{Id: 1, Link: "link1"}, {Id: 2, Link: "link2"}},
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

			images, err := storage.GetImageLinksByUserId(context.Background(), tt.userId)
			require.ErrorIs(t, err, tt.queryErr)
			for i, image := range images {
				if image != tt.wantImages[i] {
					t.Errorf("GetImageLinksByUserId() images[%d] = %v, want %v", i, image, tt.wantImages[i])
				}
			}
		})
	}
}
