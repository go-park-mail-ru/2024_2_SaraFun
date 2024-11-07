package updateprofile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/models"
	"sparkit/internal/pkg/profile/delivery/updateprofile/mocks"
	"sparkit/internal/utils/consts"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestUpdateProfileHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := zap.NewNop()
	mockProfileService := updateprofile_mocks.updateprofile_mocks.NewMockProfileService(mockCtrl)
	mockSessionService := updateprofile_mocks.NewMockSessionService(mockCtrl)
	mockUserService := updateprofile_mocks.updateprofile_mocks.NewMockUserService(mockCtrl)

	handler := NewHandler(mockProfileService, mockSessionService, mockUserService, logger)
	tests := []struct {
		name               string
		cookieValue        string
		userId             int
		getUserIDError     error
		profileId          int
		getProfileIDError  error
		profile            models.Profile
		sendInvalidJSON    bool
		updateProfileError error
		expectedStatus     int
		expectedResponse   string
	}{
		{
			name:             "successful profile update",
			cookieValue:      "valid-session-id",
			userId:           1,
			profileId:        1001,
			profile:          models.Profile{FirstName: "John", LastName: "Doe"},
			expectedStatus:   http.StatusOK,
			expectedResponse: "ok",
		},
		{
			name:             "missing session cookie",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "session not found\n",
		},
		{
			name:             "error getting user ID from session",
			cookieValue:      "invalid-session-id",
			getUserIDError:   errors.New("session service error"),
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "user not found\n",
		},
		{
			name:              "error getting profile ID by user ID",
			cookieValue:       "valid-session-id",
			userId:            1,
			getProfileIDError: errors.New("user service error"),
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  "profile not found\n",
		},
		{
			name:             "error decoding JSON body",
			cookieValue:      "valid-session-id",
			userId:           1,
			profileId:        1001,
			sendInvalidJSON:  true, // Будем посылать некорректный JSON
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "invalid character 'i' looking for beginning of value\n", // Ожидаемое сообщение от декодера
		},
		{
			name:               "error updating profile",
			cookieValue:        "valid-session-id",
			userId:             1,
			profileId:          1001,
			profile:            models.Profile{FirstName: "John", LastName: "Doe"},
			updateProfileError: errors.New("database error"),
			expectedStatus:     http.StatusInternalServerError,
			expectedResponse:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.sendInvalidJSON {
				t.Log("sendInvalidJson if")
				req = httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer([]byte("invalid-json")))
				req.Header.Set("Content-Type", "application/json")
			} else if tt.profile.FirstName != "" || tt.profile.LastName != "" {
				t.Log("sendInvalidJson else if")
				bodyBytes, err := json.Marshal(tt.profile)
				if err != nil {
					t.Fatalf("Не удалось сериализовать профиль: %v", err)
				}
				t.Log("make NewRequest")
				req = httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				t.Log("after")
			} else {
				t.Log("sendInvalidJson else")
				req = httptest.NewRequest(http.MethodPost, "/updateprofile", nil)
			}

			if tt.cookieValue != "" {
				t.Log("if cookie value")
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				mockSessionService.EXPECT().
					GetUserIDBySessionID(gomock.Any(), tt.cookieValue).
					Return(tt.userId, tt.getUserIDError).
					Times(1)
				t.Log("good if cookie value")
			}

			if tt.cookieValue != "" && tt.getUserIDError == nil {
				t.Log("cookie getUser")
				if tt.getProfileIDError == nil {
					t.Log("cookie getUser if")
					mockUserService.EXPECT().
						GetProfileIdByUserId(gomock.Any(), tt.userId).
						Return(tt.profileId, nil).
						Times(1)
					t.Log("cookie getUser if good")
				} else {
					t.Log("cookie getUser else")
					mockUserService.EXPECT().
						GetProfileIdByUserId(gomock.Any(), tt.userId).
						Return(0, tt.getProfileIDError).
						Times(1)
				}
			}

			if tt.cookieValue != "" && tt.getUserIDError == nil && tt.getProfileIDError == nil && !tt.sendInvalidJSON {
				t.Log("all ")
				if tt.updateProfileError == nil {
					t.Log("all if")
					mockProfileService.EXPECT().
						UpdateProfile(gomock.Any(), tt.profileId, tt.profile).
						Return(nil).
						Times(1)
					t.Log("good all if")
				} else {
					t.Log("all else")
					mockProfileService.EXPECT().
						UpdateProfile(gomock.Any(), tt.profileId, tt.profile).
						Return(tt.updateProfileError).
						Times(1)
					t.Log("good all else")
				}
			}

			w := httptest.NewRecorder()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel() // Отменяем контекст после завершения работы
			ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
			req = req.WithContext(ctx)
			t.Log("new recorder")
			handler.Handle(w, req)
			t.Log("good handle")
			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
			}
			t.Log("good")
		})
	}
}
