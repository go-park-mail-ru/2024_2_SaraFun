package logout

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	logout_mocks "sparkit/internal/handlers/logout/mocks"
	"sparkit/internal/utils/consts"
)

func TestLogoutHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name               string
		method             string
		cookieValue        string
		deleteSessionError error
		expectedStatus     int
		expectedResponse   string
		expectCookie       bool
		logger             *zap.Logger
	}{
		{
			name:             "wrong method",
			method:           http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedResponse: "method is not allowed\n",
			expectCookie:     false,
			logger:           logger,
		},
		{
			name:             "session cookie not found",
			method:           http.MethodGet,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
			expectCookie:     false,
			logger:           logger,
		},
		{
			name:               "failed to delete session",
			method:             http.MethodGet,
			cookieValue:        "invalid-session-id",
			expectedStatus:     http.StatusInternalServerError,
			expectedResponse:   "failed to logout\n",
			deleteSessionError: errors.New("database error"),
			expectCookie:       false,
			logger:             logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := logout_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(service, tt.logger)

			var req *http.Request
			if tt.cookieValue != "" {
				req = httptest.NewRequest(tt.method, "/logout", nil)
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				service.EXPECT().DeleteSession(gomock.Any(), tt.cookieValue).Return(tt.deleteSessionError).Times(1)
			} else {
				req = httptest.NewRequest(tt.method, "/logout", nil)
				if tt.method == http.MethodGet {
					service.EXPECT().DeleteSession(gomock.Any(), "").Times(0)
				}
			}

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

			cookies := w.Result().Cookies()

			if tt.expectCookie {
				if len(cookies) == 0 || cookies[0].Name != consts.SessionCookie || cookies[0].Value != "" || cookies[0].Expires.Before(time.Now()) {
					t.Errorf("cookie not set correctly")
				}
			} else {
				if len(cookies) > 0 {
					t.Errorf("unexpected cookie set")
				}
			}
		})
	}
}
