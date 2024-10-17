package signin

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	sparkiterrors "sparkit/internal/errors"
	signin_mocks "sparkit/internal/handlers/signin/mocks"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"testing"
)

func TestSigninHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                string
		method              string
		path                string
		body                []byte
		checkPasswordError  error
		createSessionError  error
		expectedStatus      int
		expectedMessage     string
		checkPasswordCalled bool
		createSessionCalled bool
	}{
		{
			name:                "successful login",
			method:              "POST",
			path:                "http://localhost:8080/signin",
			body:                []byte(`{"username":"user1", "password":"password1"}`),
			checkPasswordError:  nil,
			createSessionError:  nil,
			expectedStatus:      http.StatusOK,
			expectedMessage:     "ok",
			checkPasswordCalled: true,
			createSessionCalled: true,
		},
		{
			name:                "wrong credentials",
			method:              "POST",
			path:                "http://localhost:8080/signin",
			body:                []byte(`{"username":"user1", "password":"wrongpassword"}`),
			checkPasswordError:  sparkiterrors.ErrWrongCredentials,
			expectedStatus:      http.StatusPreconditionFailed,
			expectedMessage:     "wrong credentials\n",
			checkPasswordCalled: true,
			createSessionCalled: false,
		},
		{
			name:                "failed session creation",
			method:              "POST",
			path:                "http://localhost:8080/signin",
			body:                []byte(`{"username":"user1", "password":"password1"}`),
			checkPasswordError:  nil,
			createSessionError:  errors.New("session creation error"),
			expectedStatus:      http.StatusInternalServerError,
			expectedMessage:     "Не удалось создать сессию\n",
			checkPasswordCalled: true,
			createSessionCalled: true,
		},
		{
			name:                "wrong method",
			method:              "GET",
			path:                "http://localhost:8080/signin",
			body:                nil,
			expectedStatus:      http.StatusMethodNotAllowed,
			expectedMessage:     "Method not allowed\n",
			checkPasswordCalled: false,
			createSessionCalled: false,
		},
		{
			name:                "invalid request format",
			method:              "POST",
			path:                "http://localhost:8080/signin",
			body:                []byte(`invalid_json`),
			expectedStatus:      http.StatusBadRequest,
			expectedMessage:     "Неверный формат данных\n",
			checkPasswordCalled: false,
			createSessionCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := signin_mocks.NewMockUserService(mockCtrl)
			sessionService := signin_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(userService, sessionService)

			// Настройка вызовов `CheckPassword`
			if tt.checkPasswordCalled {
				userService.EXPECT().CheckPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.User{Username: "user1", Password: "hashedpassword"}, tt.checkPasswordError).Times(1)
			} else {
				userService.EXPECT().CheckPassword(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			}

			// Настройка вызовов `CreateSession`
			if tt.createSessionCalled {
				sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(models.Session{SessionID: "session_id"}, tt.createSessionError).Times(1)
			} else {
				sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			// Проверка статуса и тела ответа
			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}

			// Проверка установки куки для успешного логина
			if tt.expectedStatus == http.StatusOK && tt.createSessionError == nil {
				cookie := w.Result().Cookies()
				if len(cookie) == 0 || cookie[0].Name != consts.SessionCookie {
					t.Errorf("expected session cookie to be set")
				}
			}
		})
	}
}
