package updateprofile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/handlers/updateprofile"
	"sparkit/internal/handlers/updateprofile/mocks"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUpdateProfileHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProfileService := mocks.NewMockProfileService(mockCtrl)
	mockSessionService := mocks.NewMockSessionService(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)

	logger := zaptest.NewLogger(t)

	handler := updateprofile.NewHandler(mockProfileService, mockSessionService, mockUserService, logger)

	tests := []struct {
		name               string
		cookieValue        string
		userId             int
		getUserIDError     error
		profileId          int64
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			var req *http.Request
			if tt.sendInvalidJSON {

				req = httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer([]byte("invalid-json")))
				req.Header.Set("Content-Type", "application/json")
			} else if tt.profile.FirstName != "" || tt.profile.LastName != "" {

				bodyBytes, err := json.Marshal(tt.profile)
				if err != nil {
					t.Fatalf("Не удалось сериализовать профиль: %v", err)
				}
				req = httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {

				req = httptest.NewRequest(http.MethodPost, "/updateprofile", nil)
			}

			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
				mockSessionService.EXPECT().
					GetUserIDBySessionID(gomock.Any(), tt.cookieValue).
					Return(tt.userId, tt.getUserIDError).
					Times(1)
			}

			if tt.cookieValue != "" && tt.getUserIDError == nil {
				if tt.getProfileIDError == nil {
					mockUserService.EXPECT().
						GetProfileIdByUserId(gomock.Any(), tt.userId).
						Return(tt.profileId, nil).
						Times(1)
				} else {
					mockUserService.EXPECT().
						GetProfileIdByUserId(gomock.Any(), tt.userId).
						Return(int64(0), tt.getProfileIDError).
						Times(1)
				}
			}

			if tt.cookieValue != "" && tt.getUserIDError == nil && tt.getProfileIDError == nil && !tt.sendInvalidJSON {
				if tt.updateProfileError == nil {
					mockProfileService.EXPECT().
						UpdateProfile(gomock.Any(), tt.profileId, tt.profile).
						Return(nil).
						Times(1)
				} else {
					mockProfileService.EXPECT().
						UpdateProfile(gomock.Any(), tt.profileId, tt.profile).
						Return(tt.updateProfileError).
						Times(1)
				}
			}

			w := httptest.NewRecorder()

			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
			}
		})
	}
}
