package signin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
			name:              "Invalid Method",
			method:            http.MethodGet,
			requestBody:       nil,
			mockUserClient:    func(mock *signin_mocks.MockPersonalitiesClient) {},
			mockSessionClient: func(mock *signin_mocks.MockAuthClient) {},
			expectedStatus:    http.StatusMethodNotAllowed,
			expectedResponse:  "Method not allowed\n",
			expectedCookie:    nil,
		},
		{
			name:              "Malformed JSON Body",
			method:            http.MethodPost,
			requestBody:       `{"username": "testuser", "password"}`, // Invalid JSON
			mockUserClient:    func(mock *signin_mocks.MockPersonalitiesClient) {},
			mockSessionClient: func(mock *signin_mocks.MockAuthClient) {},
			expectedStatus:    http.StatusBadRequest,
			expectedResponse:  "Неверный формат данных\n",
			expectedCookie:    nil,
		},
		{
			name:   "Incorrect Credentials",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "wronguser",
				"password": "wrongpass",
			},
			mockUserClient: func(mock *signin_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					CheckPassword(gomock.Any(), &generatedPersonalities.CheckPasswordRequest{
						Username: "wronguser",
						Password: "wrongpass",
					}).
					Return(nil, fmt.Errorf("invalid credentials"))
			},
			mockSessionClient: func(mock *signin_mocks.MockAuthClient) {},
			expectedStatus:    http.StatusPreconditionFailed,
			expectedResponse:  "wrong credentials\n",
			expectedCookie:    nil,
		},
		{
			name:   "Session Creation Failure",
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
							Profile:  1,
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
							Profile:  1,
						},
					}).
					Return(nil, fmt.Errorf("session creation error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Не удалось создать сессию\n",
			expectedCookie:   nil,
		},
		{
			name:              "Missing Request Body",
			method:            http.MethodPost,
			requestBody:       nil,
			mockUserClient:    func(mock *signin_mocks.MockPersonalitiesClient) {},
			mockSessionClient: func(mock *signin_mocks.MockAuthClient) {},
			expectedStatus:    http.StatusBadRequest,
			expectedResponse:  "Неверный формат данных\n",
			expectedCookie:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserClient := signin_mocks.NewMockPersonalitiesClient(ctrl)
			mockSessionClient := signin_mocks.NewMockAuthClient(ctrl)

			tt.mockUserClient(mockUserClient)
			tt.mockSessionClient(mockSessionClient)

			logger := zap.NewNop()
			handler := NewHandler(mockUserClient, mockSessionClient, logger)

			var req *http.Request
			if tt.requestBody != nil {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(tt.method, "/", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tt.method, "/", nil)
			}
			req = req.WithContext(context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id"))

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}

			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("expected response %q, got %q", tt.expectedResponse, rr.Body.String())
			}

			if tt.expectedCookie != nil {
				cookies := rr.Result().Cookies()
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
		})
	}
}
