package logout

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	logout_mocks "sparkit/internal/handlers/logout/mocks"
	"sparkit/internal/utils/consts"
)

func TestLogoutHandler(t *testing.T) {
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
	}{
		{
			name:             "wrong method",
			method:           http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedResponse: "method is not allowed\n",
			expectCookie:     false,
		},
		{
			name:             "session cookie not found",
			method:           http.MethodGet,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
			expectCookie:     false,
		},
		{
			name:               "failed to delete session",
			method:             http.MethodGet,
			cookieValue:        "invalid-session-id",
			expectedStatus:     http.StatusInternalServerError,
			expectedResponse:   "failed to logout\n",
			deleteSessionError: errors.New("database error"),
			expectCookie:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := logout_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(service)

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
