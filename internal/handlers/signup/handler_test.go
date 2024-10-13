package signup

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	sign_up_mocks "sparkit/internal/handlers/signup/mocks"
	"sparkit/internal/models"
	"sparkit/internal/utils/hashing"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name                    string
		method                  string
		path                    string
		body                    []byte
		addUserError            error
		addUserCallsCount       int
		createSessionError      error
		createSessionCallsCount int
		expectedStatus          int
		expectedMessage         string
	}{
		{
			name:                    "success",
			method:                  "POST",
			path:                    "http://localhost:8080/signup",
			body:                    []byte(`{"username":"Username1", "password":"Password1"}`),
			addUserError:            nil,
			addUserCallsCount:       1,
			createSessionError:      nil,
			createSessionCallsCount: 1,
			expectedStatus:          http.StatusOK,
			expectedMessage:         "ok",
		},
		{
			name:                    "wrong method",
			method:                  "GET",
			path:                    "http://localhost:8080/signup",
			body:                    nil,
			addUserError:            nil,
			addUserCallsCount:       0,
			createSessionCallsCount: 0,
			expectedStatus:          http.StatusMethodNotAllowed,
			expectedMessage:         "Method not allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := sign_up_mocks.NewMockUserService(mockCtrl)
			sessionService := sign_up_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(userService, sessionService)

			var user models.User
			if tt.body != nil {
				if err := json.Unmarshal(tt.body, &user); err != nil {
					t.Error(err)
				}
			}
			user.Password, _ = hashing.HashPassword(user.Password)
			userService.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(tt.addUserError).Times(tt.addUserCallsCount)
			sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, user models.User) (*models.Session, error) {
				session := &models.Session{
					SessionID: uuid.New().String(),
					UserID:    user.ID,
					CreatedAt: time.Now(),
				}
				return session, tt.createSessionError
			}).Times(tt.createSessionCallsCount)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			handler.Handle(w, req)
			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}

		})
	}

}
