package signup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	signup_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signup/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_PersonalitiesClient.go -package=signup_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen PersonalitiesClient
//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=signup_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name                    string
		method                  string
		requestBody             interface{}
		mockPersonalitiesClient func(mock *signup_mocks.MockPersonalitiesClient)
		mockAuthClient          func(mock *signup_mocks.MockAuthClient)
		expectedStatus          int
		expectedResponse        string
		expectedCookie          *http.Cookie
	}{
		{
			name:   "Successful Signup",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: false,
					}, nil)
				// Mock CreateProfile
				mock.EXPECT().
					CreateProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.CreateProfileResponse{
						ProfileId: 1,
					}, nil)
				// Mock RegisterUser
				mock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.RegisterUserResponse{
						UserId: 1,
					}, nil)
			},
			mockAuthClient: func(mock *signup_mocks.MockAuthClient) {
				mock.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
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
		{
			name:                    "Method Not Allowed",
			method:                  http.MethodGet,
			requestBody:             nil,
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {},
			mockAuthClient:          func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:          http.StatusMethodNotAllowed,
			expectedResponse:        "Method not allowed\n",
			expectedCookie:          nil,
		},
		{
			name:                    "Invalid JSON",
			method:                  http.MethodPost,
			requestBody:             "invalid json",
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {},
			mockAuthClient:          func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:          http.StatusBadRequest,
			expectedResponse:        "invalid character 'i' looking for beginning of value\n",
			expectedCookie:          nil,
		},
		{
			name:   "Username Already Exists",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "existinguser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "existinguser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: true,
					}, nil)
			},
			mockAuthClient:   func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "user already exists\n",
			expectedCookie:   nil,
		},
		{
			name:   "Error Checking Username Exists",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists with error
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(nil, errors.New("database error"))
			},
			mockAuthClient:   func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "failed to check username exists\n",
			expectedCookie:   nil,
		},
		{
			name:   "Error Creating Profile",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: false,
					}, nil)
				// Mock CreateProfile with error
				mock.EXPECT().
					CreateProfile(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("profile creation error"))
			},
			mockAuthClient:   func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "profile creation error\n",
			expectedCookie:   nil,
		},
		{
			name:   "Error Registering User",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: false,
					}, nil)
				// Mock CreateProfile
				mock.EXPECT().
					CreateProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.CreateProfileResponse{
						ProfileId: 1,
					}, nil)
				// Mock RegisterUser with error
				mock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("user registration error"))
			},
			mockAuthClient:   func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Внутренняя ошибка сервера\n",
			expectedCookie:   nil,
		},
		{
			name:   "Error Creating Session",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: "testpass",
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: false,
					}, nil)
				// Mock CreateProfile
				mock.EXPECT().
					CreateProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.CreateProfileResponse{
						ProfileId: 1,
					}, nil)
				// Mock RegisterUser
				mock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.RegisterUserResponse{
						UserId: 1,
					}, nil)
			},
			mockAuthClient: func(mock *signup_mocks.MockAuthClient) {
				mock.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("session creation error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Не удалось создать сессию\n",
			expectedCookie:   nil,
		},
		{
			name:   "Error Hashing Password",
			method: http.MethodPost,
			requestBody: Request{
				User: models.User{
					Username: "testuser",
					Password: strings.Repeat("a", 100), // 100 символов
					Email:    "test@example.com",
				},
				Profile: models.Profile{
					FirstName: "Test",
					LastName:  "User",
					Age:       30,
					Gender:    "Male",
					Target:    "Friendship",
					About:     "Just testing",
				},
			},
			mockPersonalitiesClient: func(mock *signup_mocks.MockPersonalitiesClient) {
				// Mock CheckUsernameExists
				mock.EXPECT().
					CheckUsernameExists(gomock.Any(), &generatedPersonalities.CheckUsernameExistsRequest{
						Username: "testuser",
					}).
					Return(&generatedPersonalities.CheckUsernameExistsResponse{
						Exists: false,
					}, nil)
				// Mock CreateProfile
				mock.EXPECT().
					CreateProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.CreateProfileResponse{
						ProfileId: 1,
					}, nil)
			},
			mockAuthClient:   func(mock *signup_mocks.MockAuthClient) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "bad password\n",
			expectedCookie:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем контроллер gomock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем моки клиентов
			mockPersonalitiesClient := signup_mocks.NewMockPersonalitiesClient(ctrl)
			mockAuthClient := signup_mocks.NewMockAuthClient(ctrl)

			// Настраиваем моки
			tt.mockPersonalitiesClient(mockPersonalitiesClient)
			tt.mockAuthClient(mockAuthClient)

			// Создаем логгер
			logger := zap.NewNop()

			// Создаем обработчик
			handler := NewHandler(mockPersonalitiesClient, mockAuthClient, logger)

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

			// Добавляем контекст с RequestID
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
			if !strings.Contains(rr.Body.String(), tt.expectedResponse) {
				t.Errorf("expected response body to contain %q, got %q", tt.expectedResponse, rr.Body.String())
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
