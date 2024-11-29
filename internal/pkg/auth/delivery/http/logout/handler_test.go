package logout

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	logout_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/logout/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=logout_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		cookie         *http.Cookie
		mockBehavior   func(mockAuthClient *logout_mocks.MockAuthClient)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful Logout",
			method: http.MethodGet,
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session-id",
			},
			mockBehavior: func(mockAuthClient *logout_mocks.MockAuthClient) {
				mockAuthClient.EXPECT().
					DeleteSession(gomock.Any(), &generatedAuth.DeleteSessionRequest{SessionID: "valid-session-id"}).
					Return(&generatedAuth.DeleteSessionResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "log out is complete",
		},
		{
			name:           "Method Not Allowed",
			method:         http.MethodPost,
			cookie:         nil,
			mockBehavior:   func(mockAuthClient *logout_mocks.MockAuthClient) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method is not allowed",
		},
		{
			name:           "Missing Cookie",
			method:         http.MethodGet,
			cookie:         nil,
			mockBehavior:   func(mockAuthClient *logout_mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "session not found",
		},
		{
			name:   "Delete Session Error",
			method: http.MethodGet,
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "invalid-session-id",
			},
			mockBehavior: func(mockAuthClient *logout_mocks.MockAuthClient) {
				mockAuthClient.EXPECT().
					DeleteSession(gomock.Any(), &generatedAuth.DeleteSessionRequest{SessionID: "invalid-session-id"}).
					Return(nil, errors.New("delete session error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to logout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем контроллер gomock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем мок AuthClient
			mockAuthClient := logout_mocks.NewMockAuthClient(ctrl)
			// Настраиваем поведение мока
			tt.mockBehavior(mockAuthClient)

			// Создаем логгер
			logger := zap.NewNop()
			// Создаем обработчик
			handler := NewHandler(mockAuthClient, logger)

			// Создаем HTTP-запрос
			req := httptest.NewRequest(tt.method, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			handler.Handle(rr, req)

			// Проверяем статус код
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			// Проверяем тело ответа
			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}

			// Дополнительные проверки для успешного выхода
			if tt.name == "Successful Logout" {
				// Проверяем, что cookie был очищен
				cookies := rr.Result().Cookies()
				if len(cookies) == 0 {
					t.Errorf("expected cookie to be set, but none were found")
				} else {
					logoutCookie := cookies[0]
					if logoutCookie.Name != consts.SessionCookie || logoutCookie.Value != "" {
						t.Errorf("expected cookie %q to be cleared, but got %q", consts.SessionCookie, logoutCookie.Value)
					}
					if !logoutCookie.Expires.Before(time.Now()) {
						t.Errorf("expected cookie to be expired")
					}
				}
			}
		})
	}
}
