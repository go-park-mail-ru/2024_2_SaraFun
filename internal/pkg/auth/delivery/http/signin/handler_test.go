package signin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	signin_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signin/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_PersonalitiesClient.go -package=signin_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen PersonalitiesClient
//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=signin_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name              string
		method            string
		requestBody       interface{}
		mockUserClient    func(mock *signin_mocks.MockPersonalitiesClient)
		mockSessionClient func(mock *signin_mocks.MockAuthClient)
		expectedStatus    int
		expectedResponse  string
		expectedCookie    *http.Cookie
	}{
		{
			name:   "Successful SignIn",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockUserClient: func(mock *signin_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					CheckPassword(gomock.Any(), &generatedPersonalities.CheckPasswordRequest{
						Username: "testuser",
						Password: "testpass",
					}).
					Return(&generatedPersonalities.CheckPasswordResponse{
						User: &generatedPersonalities.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
							Password: "hashedpassword",
							Profile:  1, // Изменено
						},
					}, nil)
			},
			mockSessionClient: func(mock *signin_mocks.MockAuthClient) {
				mock.EXPECT().
					CreateSession(gomock.Any(), &generatedAuth.CreateSessionRequest{
						User: &generatedAuth.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
							Password: "hashedpassword",
							Profile:  1, // Изменено
						},
					}).
					Return(&generatedAuth.CreateSessionResponse{
						Session: &generatedAuth.Session{
							SessionID: "session-id",
						},
					}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "ok",
			expectedCookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "session-id",
			},
		},
		// Остальные тестовые случаи...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем контроллер gomock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем моки клиентов
			mockUserClient := signin_mocks.NewMockPersonalitiesClient(ctrl)
			mockSessionClient := signin_mocks.NewMockAuthClient(ctrl)

			// Настраиваем поведение моков
			tt.mockUserClient(mockUserClient)
			tt.mockSessionClient(mockSessionClient)

			// Создаем логгер
			logger := zap.NewNop()

			// Создаем обработчик
			handler := NewHandler(mockUserClient, mockSessionClient, logger)

			// Создаем HTTP-запрос
			var req *http.Request
			if tt.requestBody != nil {
				var bodyBytes []byte
				switch v := tt.requestBody.(type) {
				case string:
					bodyBytes = []byte(v)
				default:
					bodyBytes, _ = json.Marshal(v)
				}
				req = httptest.NewRequest(tt.method, "/", bytes.NewReader(bodyBytes))
			} else {
				req = httptest.NewRequest(tt.method, "/", nil)
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
			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("expected response body %q, got %q", tt.expectedResponse, rr.Body.String())
			}

			// Проверяем cookie
			if tt.expectedCookie != nil {
				cookies := rr.Result().Cookies()
				if len(cookies) == 0 {
					t.Errorf("expected cookie to be set, but none were found")
				} else {
					found := false
					for _, cookie := range cookies {
						if cookie.Name == tt.expectedCookie.Name && cookie.Value == tt.expectedCookie.Value {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected cookie %v, but it was not found", tt.expectedCookie)
					}
				}
			} else {
				if len(rr.Result().Cookies()) > 0 {
					t.Errorf("expected no cookies to be set, but found %v", rr.Result().Cookies())
				}
			}
		})
	}
}
