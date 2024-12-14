package updateprofile

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
	imageservicemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/updateprofile/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

// 1. Успешный сценарий: проверяется корректное обновление профиля при валидной сессии и корректном JSON.
// 2. Отсутствие cookie: проверка, что без cookie возвращается 401 Unauthorized и соответствующее сообщение.
// 3. Ошибка при получении userID: проверка, что при ошибке в AuthClient возвращается 401 и сообщение о том, что пользователь не найден.
// 4. Ошибка при получении profileID: проверка, что при ошибке в PersonalitiesClient при получении profileID возвращается 401 и сообщение о ненайденном профиле.
// 5. Некорректный JSON: проверка, что при ошибке парсинга тела запроса возвращается 400 Bad Request.

func TestHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")

	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
	imageService := imageservicemocks.NewMockImageService(mockCtrl)

	handler := NewHandler(personalitiesClient, sessionClient, imageService, logger)

	validBody := []byte(`{
		"first_name": "John",
		"last_name": "Doe",
		"gender": "male",
		"age": 30,
		"target": "friendship",
		"about": "Hello!",
		"imgNumbers": [{"id":1,"number":2},{"id":2,"number":3}]
	}`)

	tests := []struct {
		name                     string
		method                   string
		cookieValue              string
		userID                   int32
		userIDError              error
		profileID                int32
		profileIDError           error
		requestBody              []byte
		updateProfileError       error
		updateImagesError        error
		expectedStatus           int
		expectedResponseContains string
	}{
		{
			name:                     "good test",
			method:                   http.MethodPost,
			cookieValue:              "valid_session",
			userID:                   10,
			requestBody:              validBody,
			profileID:                100,
			expectedStatus:           http.StatusOK,
			expectedResponseContains: "ok",
		},
		{
			name:                     "no cookie",
			method:                   http.MethodPost,
			cookieValue:              "",
			requestBody:              validBody,
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "session not found",
		},
		{
			name:                     "session user error",
			method:                   http.MethodPost,
			cookieValue:              "bad_session",
			userIDError:              errors.New("session error"),
			requestBody:              validBody,
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "user not found",
		},
		{
			name:                     "profileID error",
			method:                   http.MethodPost,
			cookieValue:              "valid_session",
			userID:                   10,
			profileIDError:           errors.New("profile id error"),
			requestBody:              validBody,
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "profile not found",
		},
		{
			name:           "bad json",
			method:         http.MethodPost,
			cookieValue:    "valid_session",
			userID:         10,
			profileID:      100,
			requestBody:    []byte(`{bad json`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cookieValue != "" {
				getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
				if tt.userIDError == nil {
					userResp := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: tt.userID}
					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
						Return(userResp, nil).Times(1)
				} else {
					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
						Return(nil, tt.userIDError).Times(1)
				}
			}

			if tt.userIDError == nil && tt.cookieValue != "" {
				getProfileIDReq := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: tt.userID}
				if tt.profileIDError == nil && tt.profileID != 0 {
					profileIDResp := &generatedPersonalities.GetProfileIDByUserIDResponse{ProfileID: tt.profileID}
					personalitiesClient.EXPECT().GetProfileIDByUserID(gomock.Any(), getProfileIDReq).
						Return(profileIDResp, nil).Times(1)
				} else if tt.profileIDError != nil {
					personalitiesClient.EXPECT().GetProfileIDByUserID(gomock.Any(), getProfileIDReq).
						Return(nil, tt.profileIDError).Times(1)
				}
			}

			if tt.userIDError == nil && tt.profileIDError == nil && tt.profileID != 0 &&
				tt.cookieValue != "" && bytes.Equal(tt.requestBody, validBody) {
				if tt.updateProfileError == nil {
					personalitiesClient.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).
						Return(&generatedPersonalities.UpdateProfileResponse{}, nil).Times(1)
				} else {
					personalitiesClient.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).
						Return(nil, tt.updateProfileError).Times(1)
				}

				if tt.updateProfileError == nil {
					if tt.updateImagesError == nil {
						imageService.EXPECT().UpdateOrdNumbers(gomock.Any(), gomock.Any()).
							Return(nil).Times(1)
					} else {
						imageService.EXPECT().UpdateOrdNumbers(gomock.Any(), gomock.Any()).
							Return(tt.updateImagesError).Times(1)
					}
				}
			}

			req := httptest.NewRequest(tt.method, "/update/profile", bytes.NewBuffer(tt.requestBody))
			req = req.WithContext(ctx)
			if tt.cookieValue != "" {
				cookie := &http.Cookie{
					Name:  consts.SessionCookie,
					Value: tt.cookieValue,
				}
				req.AddCookie(cookie)
			}
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v", tt.name, w.Code, tt.expectedStatus)
			}
			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
				t.Errorf("%s: handler returned unexpected body: got %v want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		(len(substr) > 0 && len(s) > len(substr) && findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
