package checkauth

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/pkg/auth/delivery/checkauth/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"sparkit/internal/utils/consts"
)

func TestCheckAuthHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name              string
		method            string
		cookieValue       string
		checkSessionError error
		expectedStatus    int
		expectedResponse  string
		logger            *zap.Logger
	}{
		{
			name:              "successful session check",
			method:            http.MethodGet,
			cookieValue:       "valid-session-id",
			expectedStatus:    http.StatusOK,
			expectedResponse:  "ok",
			checkSessionError: nil,
			logger:            logger,
		},
		{
			name:             "wrong method",
			method:           http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedResponse: "method is not allowed\n",
			logger:           logger,
		},
		{
			name:             "session cookie not found",
			method:           http.MethodGet,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
			logger:           logger,
		},
		{
			name:              "invalid session",
			method:            http.MethodGet,
			cookieValue:       "invalid-session-id",
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  "session is not valid\n",
			checkSessionError: errors.New("invalid session"),
			logger:            logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := checkauth_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(service, tt.logger)

			if tt.method == http.MethodGet && tt.cookieValue != "" {
				req := httptest.NewRequest(tt.method, "/checkauth", nil)
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				w := httptest.NewRecorder()

				if tt.checkSessionError != nil {
					service.EXPECT().CheckSession(gomock.Any(), tt.cookieValue).Return(tt.checkSessionError).Times(1)
				} else {
					service.EXPECT().CheckSession(gomock.Any(), tt.cookieValue).Return(nil).Times(1)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel() // Отменяем контекст после завершения работы
				ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
				req = req.WithContext(ctx)
				handler.Handle(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
				}

				if w.Body.String() != tt.expectedResponse {
					t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
				}
			} else {
				req := httptest.NewRequest(tt.method, "/checkauth", nil)
				w := httptest.NewRecorder()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel() // Отменяем контекст после завершения работы
				ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
				req = req.WithContext(ctx)
				handler.Handle(w, req)

				if w.Code != tt.expectedStatus {
					t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
				}

				if w.Body.String() != tt.expectedResponse {
					t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
				}
			}
		})
	}
}
