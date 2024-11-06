package uploadimage_test

import (
	"bytes"
	"errors"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/handlers/uploadimage"
	"sparkit/internal/handlers/uploadimage/mocks"
	"sparkit/internal/utils/consts"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUploadImageHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockImageService := mocks.NewMockImageService(mockCtrl)
	mockSessionService := mocks.NewMockSessionService(mockCtrl)
	logger := zaptest.NewLogger(t)

	handler := uploadimage.NewHandler(mockImageService, mockSessionService, logger)

	tests := []struct {
		name             string
		cookieValue      string
		userId           int
		getUserIDError   error
		saveImageID      int
		saveImageError   error
		expectedStatus   int
		expectedResponse string
		mockFile         bool
	}{
		{
			name:             "successful upload",
			cookieValue:      "valid-session-id",
			userId:           1,
			saveImageID:      12345,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"ImageId":12345}`,
			mockFile:         true,
		},
		{
			name:             "session not found",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
			mockFile:         true,
		},
		{
			name:             "user session error",
			cookieValue:      "invalid-session-id",
			getUserIDError:   errors.New("session service error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "user session err\n",
			mockFile:         true,
		},

		{
			name:             "save image error",
			cookieValue:      "valid-session-id",
			userId:           1,
			saveImageError:   errors.New("image service error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "save image err\n",
			mockFile:         true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			var req *http.Request
			if tt.mockFile {
				var b bytes.Buffer
				w := multipart.NewWriter(&b)
				fileWriter, err := w.CreateFormFile("image", "test.jpg")
				if err != nil {
					t.Fatalf("Не удалось создать form файл: %v", err)
				}
				_, err = fileWriter.Write([]byte("fake image content"))
				if err != nil {
					t.Fatalf("Не удалось записать контент в form файл: %v", err)
				}
				w.Close()
				req = httptest.NewRequest(http.MethodPost, "/upload", &b)
				req.Header.Set("Content-Type", w.FormDataContentType())
			} else {

				req = httptest.NewRequest(http.MethodPost, "/upload", nil)
			}

			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				mockSessionService.EXPECT().
					GetUserIDBySessionID(gomock.Any(), tt.cookieValue).
					Return(tt.userId, tt.getUserIDError).
					Times(1)
			}

			if tt.cookieValue != "" && tt.getUserIDError == nil && tt.mockFile {
				mockImageService.EXPECT().
					SaveImage(gomock.Any(), gomock.Any(), gomock.Any(), tt.userId).
					Return(tt.saveImageID, tt.saveImageError).
					Times(1)
			}

			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
			}
		})
	}
}
