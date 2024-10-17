package checkauth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	checkauth_mocks "sparkit/internal/handlers/checkauth/mocks"
	"sparkit/internal/utils/consts"
)

func TestCheckAuthHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name              string
		method            string
		cookieValue       string
		checkSessionError error
		expectedStatus    int
		expectedResponse  string
	}{
		{
			name:              "successful session check",
			method:            http.MethodGet,
			cookieValue:       "valid-session-id",
			expectedStatus:    http.StatusOK,
			expectedResponse:  "ok",
			checkSessionError: nil,
		},
		{
			name:             "wrong method",
			method:           http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedResponse: "method is not allowed\n",
		},
		{
			name:             "session cookie not found",
			method:           http.MethodGet,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
		},
		{
			name:              "invalid session",
			method:            http.MethodGet,
			cookieValue:       "invalid-session-id",
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  "session is not valid\n",
			checkSessionError: errors.New("invalid session"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := checkauth_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(service)

			if tt.method == http.MethodGet && tt.cookieValue != "" {
				req := httptest.NewRequest(tt.method, "/checkauth", nil)
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				w := httptest.NewRecorder()
				if tt.checkSessionError != nil {
					service.EXPECT().CheckSession(gomock.Any(), tt.cookieValue).Return(tt.checkSessionError).Times(1)
				} else {
					service.EXPECT().CheckSession(gomock.Any(), tt.cookieValue).Return(nil).Times(1)
				}
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
