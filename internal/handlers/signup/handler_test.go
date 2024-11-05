package signup

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	sign_up_mocks "sparkit/internal/handlers/signup/mocks"
	"sparkit/internal/models"
	"sparkit/internal/utils/hashing"
	"testing"
	"time"
)

type TestRequest struct {
	User    models.User    `json:"user"`
	Profile models.Profile `json:"profile"`
}

func TestHandler(t *testing.T) {
	logger := zap.NewNop()
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
		createProfileId         int64
		createProfileError      error
		createProfileCallsCount int
		expectedStatus          int
		expectedMessage         string
		logger                  *zap.Logger
	}{
		{
			name:   "success",
			method: "POST",
			path:   "http://localhost:8080/signup",
			body: []byte(`{
				"user": {
        			"username": "User1",
        			"password": "user234"
   			 	},
    			"profile": {
        			"gender": "user",
        			"age": 40
    			}
			}`),
			addUserError:            nil,
			addUserCallsCount:       1,
			createSessionError:      nil,
			createSessionCallsCount: 1,
			createProfileId:         1,
			createProfileError:      nil,
			createProfileCallsCount: 1,
			expectedStatus:          http.StatusOK,
			expectedMessage:         "ok",
			logger:                  logger,
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
			logger:                  logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := sign_up_mocks.NewMockUserService(mockCtrl)
			sessionService := sign_up_mocks.NewMockSessionService(mockCtrl)
			profileService := sign_up_mocks.NewMockProfileService(mockCtrl)
			handler := NewHandler(userService, sessionService, profileService, tt.logger)

			var reqB TestRequest
			if tt.body != nil {
				if err := json.Unmarshal(tt.body, &reqB); err != nil {
					t.Error(err)
				}
			}
			reqB.User.Password, _ = hashing.HashPassword(reqB.User.Password)
			profileService.EXPECT().CreateProfile(gomock.Any(), reqB.Profile).Return(tt.createProfileId, tt.createProfileError).Times(tt.createProfileCallsCount)
			userService.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(tt.addUserError).Times(tt.addUserCallsCount)
			sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, user models.User) (*models.Session, error) {
				session := &models.Session{
					SessionID: uuid.New().String(),
					UserID:    reqB.User.ID,
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
