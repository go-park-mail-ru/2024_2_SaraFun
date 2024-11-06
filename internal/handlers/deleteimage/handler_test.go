package deleteimage_test

import (
	"errors"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/handlers/deleteimage"
	"sparkit/internal/handlers/deleteimage/mocks"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestDeleteImageHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockImageService := mocks.NewMockImageService(mockCtrl)

	logger := zaptest.NewLogger(t)

	handler := deleteimage.NewHandler(mockImageService, logger)

	router := mux.NewRouter()
	router.HandleFunc("/deleteimage/{imageId}", handler.Handle).Methods("DELETE")
	router.HandleFunc("/deleteimage/", handler.Handle).Methods("DELETE") // Регистрация маршрута без imageId

	tests := []struct {
		name             string
		imageId          string
		deleteImageError error
		expectedStatus   int
		expectedResponse string
	}{

		{
			name:             "invalid image ID",
			imageId:          "abc",
			deleteImageError: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "strconv.Atoi: parsing \"abc\": invalid syntax\n",
		},
		{
			name:             "error deleting image",
			imageId:          "456",
			deleteImageError: errors.New("database error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "database error\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var url string
			if tt.imageId != "" {
				url = "/deleteimage/" + tt.imageId
			} else {
				url = "/deleteimage/"
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)

			if tt.imageId != "" {
				imageIdInt, err := strconv.Atoi(tt.imageId)
				if err == nil {
					mockImageService.EXPECT().
						DeleteImage(gomock.Any(), imageIdInt).
						Return(tt.deleteImageError).
						Times(1)
				}
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
			}
		})
	}
}