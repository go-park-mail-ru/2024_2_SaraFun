package changepassword

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

func TestHandler_Simplified(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")

	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authClient := authmocks.NewMockAuthClient(mockCtrl)
	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)

	handler := NewHandler(authClient, personalitiesClient, logger)

	validRequest := Request{
		CurrentPassword: "oldpass",
		NewPassword:     "newpass",
	}
	validBody, _ := json.Marshal(validRequest)

	tests := []struct {
		name                     string
		cookieValue              string
		body                     []byte
		expectedStatus           int
		expectedResponseContains string
	}{
		{
			name:                     "bad json",
			cookieValue:              "valid_session",
			body:                     []byte(`{bad json`),
			expectedStatus:           http.StatusBadRequest,
			expectedResponseContains: "Нам не нравится ваш запрос :(",
		},
		{
			name:                     "no cookie",
			cookieValue:              "",
			body:                     validBody,
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "Вы не авторизованы!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/change_password", bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: tt.cookieValue,
				})
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
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && string(s[0:len(substr)]) == substr) ||
		(len(s) > len(substr) && string(s[len(s)-len(substr):]) == substr) ||
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
