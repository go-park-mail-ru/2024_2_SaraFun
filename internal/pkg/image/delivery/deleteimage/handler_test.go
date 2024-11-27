package deleteimage_test

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/deleteimage"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/deleteimage/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel() // Отменяем контекст после завершения работы
			ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
			req = req.WithContext(ctx)
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
