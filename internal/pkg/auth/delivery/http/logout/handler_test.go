package logout

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

func TestHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := authmocks.NewMockAuthClient(mockCtrl)
	handler := NewHandler(client, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")

	tests := []struct {
		name                     string
		method                   string
		cookieValue              string
		deleteSessionError       error
		expectedStatus           int
		expectedResponseContains string
	}{
		{
			name:                     "not GET method",
			method:                   http.MethodPost,
			expectedStatus:           http.StatusMethodNotAllowed,
			expectedResponseContains: "method is not allowed",
		},
		{
			name:                     "no cookie",
			method:                   http.MethodGet,
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "session not found",
		},
		{
			name:                     "delete session error",
			method:                   http.MethodGet,
			cookieValue:              "bad_session",
			deleteSessionError:       errors.New("delete session error"),
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "failed to logout",
		},
		{
			name:                     "good test",
			method:                   http.MethodGet,
			cookieValue:              "valid_session",
			expectedStatus:           http.StatusOK,
			expectedResponseContains: "Вы успешно вышли из учетной записи",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cookieValue != "" {
				delReq := &generatedAuth.DeleteSessionRequest{SessionID: tt.cookieValue}
				if tt.deleteSessionError == nil {
					client.EXPECT().DeleteSession(gomock.Any(), delReq).
						Return(&generatedAuth.DeleteSessionResponse{}, nil).Times(1)
				} else {
					client.EXPECT().DeleteSession(gomock.Any(), delReq).
						Return(nil, tt.deleteSessionError).Times(1)
				}
			}

			req := httptest.NewRequest(tt.method, "/logout", nil)
			// Устанавливаем ctx с request_id
			req = req.WithContext(ctx)
			if tt.cookieValue != "" {
				cookie := &http.Cookie{
					Name:  consts.SessionCookie,
					Value: tt.cookieValue,
				}
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v, want %v", tt.name, w.Code, tt.expectedStatus)
			}
			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
				t.Errorf("%s: handler returned unexpected body: got %v, want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
			(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
			(len(substr) > 0 && len(s) > len(substr) && findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
