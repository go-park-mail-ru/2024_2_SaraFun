package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/usecase/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"testing"
	"time"
)

func TestSaveImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	testFile, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name                   string
		file                   multipart.File
		fileExt                string
		userId                 int
		ordNumber              int
		expectedSaveImageId    int
		expectedSaveImageError error
		expectedSaveImageCount int
		logger                 *zap.Logger
		wantId                 int
	}{
		{
			name:                   "successful test",
			file:                   testFile,
			fileExt:                ".png",
			userId:                 1,
			ordNumber:              1,
			expectedSaveImageId:    1,
			expectedSaveImageError: nil,
			expectedSaveImageCount: 1,
			logger:                 logger,
			wantId:                 1,
		},
		{
			name:                   "bad test",
			file:                   testFile,
			fileExt:                ".txt",
			userId:                 1,
			ordNumber:              1,
			expectedSaveImageId:    0,
			expectedSaveImageError: errors.New("error"),
			expectedSaveImageCount: 1,
			logger:                 logger,
			wantId:                 -1,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().SaveImage(ctx, gomock.Any(), tt.fileExt, tt.userId, tt.ordNumber).
				Return(tt.expectedSaveImageId, tt.expectedSaveImageError).
				Times(tt.expectedSaveImageCount)

			u := New(repo, logger)
			id, err := u.SaveImage(ctx, tt.file, tt.fileExt, tt.userId, tt.ordNumber)
			require.ErrorIs(t, err, tt.expectedSaveImageError)
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
	images := []models.Image{{Id: 1, Link: "link1"},
		{Id: 2, Link: "link2"},
	}
	tests := []struct {
		name           string
		userId         int
		expectedImages []models.Image
		expectedError  error
		expectedCount  int
		logger         *zap.Logger
		wantImages     []models.Image
	}{
		{
			name:           "successful test",
			userId:         1,
			expectedImages: images,
			expectedError:  nil,
			expectedCount:  1,
			logger:         logger,
			wantImages:     images,
		},
		{
			name:           "bad test",
			userId:         1,
			expectedImages: nil,
			expectedError:  errors.New("error"),
			expectedCount:  1,
			logger:         logger,
			wantImages:     nil,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetImageLinksByUserId(ctx, tt.userId).Return(tt.expectedImages, tt.expectedError).Times(tt.expectedCount)
			u := New(repo, logger)
			imgs, err := u.GetImageLinksByUserId(ctx, tt.userId)
			require.ErrorIs(t, err, tt.expectedError)
			for i, img := range imgs {
				if img != tt.wantImages[i] {
					t.Errorf("GetImageLinksByUserId() img = %v, want %v", img, tt.wantImages[i])
				}
			}
		})
	}
}

func TestDeleteImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()

	tests := []struct {
		name          string
		userId        int
		expectedError error
		expectedCount int
		logger        *zap.Logger
	}{
		{
			name:          "successful test",
			userId:        1,
			expectedError: nil,
			expectedCount: 1,
			logger:        logger,
		},
		{
			name:          "bad test",
			userId:        1,
			expectedError: errors.New("error"),
			expectedCount: 1,
			logger:        logger,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().DeleteImage(ctx, tt.userId).Return(tt.expectedError).Times(tt.expectedCount)

			u := New(repo, logger)
			err := u.DeleteImage(ctx, tt.userId)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
