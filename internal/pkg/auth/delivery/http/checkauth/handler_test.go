package checkauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	checkauth_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/checkauth/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=checkauth_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		cookie         *http.Cookie
		mockBehavior   func(mockAuthClient *checkauth_mocks.MockAuthClient)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful Check",
			method: http.MethodGet,
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session-id",
			},
			mockBehavior: func(mockAuthClient *checkauth_mocks.MockAuthClient) {
				mockAuthClient.EXPECT().
					CheckSession(gomock.Any(), &generatedAuth.CheckSessionRequest{SessionID: "valid-session-id"}).
					Return(&generatedAuth.CheckSessionResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:           "Method Not Allowed",
			method:         http.MethodPost,
			cookie:         nil,
			mockBehavior:   func(mockAuthClient *checkauth_mocks.MockAuthClient) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method is not allowed",
		},
		{
			name:           "Missing Cookie",
			method:         http.MethodGet,
			cookie:         nil,
			mockBehavior:   func(mockAuthClient *checkauth_mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "session not found",
		},
		{
			name:   "Invalid Session",
			method: http.MethodGet,
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "invalid-session-id",
			},
			mockBehavior: func(mockAuthClient *checkauth_mocks.MockAuthClient) {
				mockAuthClient.EXPECT().
					CheckSession(gomock.Any(), &generatedAuth.CheckSessionRequest{SessionID: "invalid-session-id"}).
					Return(nil, errors.New("session invalid"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "bad session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthClient := checkauth_mocks.NewMockAuthClient(ctrl)
			logger := zap.NewNop()
			handler := NewHandler(mockAuthClient, logger)

			tt.mockBehavior(mockAuthClient)

			req := httptest.NewRequest(tt.method, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.Handle(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
