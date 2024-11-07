package getcurrentprofile

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/models"
	"sparkit/internal/pkg/profile/delivery/getcurrentprofile/mocks"
	"sparkit/internal/utils/consts"
	"testing"
	"time"
)

type TestResponse struct {
	Profile models.Profile `json:"profile"`
	Images  []models.Image `json:"images"`
}

func TestHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
		//id int
		//GetImageLinks
		expectedGetImageLinksByUserId_Images []models.Image
		expectedGetImageLinksByUserId_Error  error
		expectedGetImageLinksByUserId_Count  int
		//GetProfile
		expectedGetProfile_Profile models.Profile
		expectedGetProfile_Error   error
		expectedGetProfile_Count   int
		//GetProfileByUser
		expectedGetProfileIdByUserId_ProfileId int
		expectedGetProfileIdByUserId_Error     error
		expectedGetProfileIdByUserId_Count     int
		//session
		expectedGetUserIdBySessionId_UserId int
		expectedGetUserIdBySessionId_Error  error
		expectedGetUserIdBySessionId_Count  int
		expectedStatus                      int
		expectedMessage                     string
		logger                              *zap.Logger
	}{
		{
			name:   "succesfull test",
			method: "GET",
			path:   "http://localhost:8080/profile/{1}",
			expectedGetImageLinksByUserId_Images: []models.Image{{Id: 1, Link: "link1"},
				{Id: 2, Link: "link2"},
			},
			expectedGetImageLinksByUserId_Error:    nil,
			expectedGetImageLinksByUserId_Count:    1,
			expectedGetProfile_Profile:             models.Profile{FirstName: "Kirill"},
			expectedGetProfile_Error:               nil,
			expectedGetProfile_Count:               1,
			expectedGetProfileIdByUserId_ProfileId: 1,
			expectedGetProfileIdByUserId_Error:     nil,
			expectedGetProfileIdByUserId_Count:     1,
			expectedGetUserIdBySessionId_UserId:    1,
			expectedGetUserIdBySessionId_Error:     nil,
			expectedGetUserIdBySessionId_Count:     1,
			expectedStatus:                         http.StatusOK,
			expectedMessage:                        "{\"profile\":{\"id\":0,\"first_name\":\"Kirill\"},\"images\":[{\"id\":1,\"link\":\"link1\"},{\"id\":2,\"link\":\"link2\"}]}",
			logger:                                 logger,
		},
		{
			name:                                 "bad test",
			method:                               "GET",
			path:                                 "http://localhost:8080/profile/{2}",
			expectedGetImageLinksByUserId_Images: []models.Image{},
			expectedGetImageLinksByUserId_Error:  errors.New("error"),
			expectedGetImageLinksByUserId_Count:  1,
			expectedGetProfileIdByUserId_Count:   1,
			expectedGetProfile_Count:             0,
			expectedGetUserIdBySessionId_Count:   1,
			expectedStatus:                       http.StatusInternalServerError,
			expectedMessage:                      "error\n",
			logger:                               logger,
		},
		{
			name:                                 "bad get profile test",
			method:                               "GET",
			path:                                 "http://localhost:8080/profile/{2}",
			expectedGetImageLinksByUserId_Images: []models.Image{},
			expectedGetImageLinksByUserId_Error:  nil,
			expectedGetImageLinksByUserId_Count:  1,
			expectedGetProfileIdByUserId_Count:   1,
			expectedGetProfile_Count:             1,
			expectedGetProfile_Error:             errors.New("error"),
			expectedGetUserIdBySessionId_Count:   1,
			expectedStatus:                       http.StatusInternalServerError,
			expectedMessage:                      "error\n",
			logger:                               logger,
		},
		{
			name:                                 "bad get profile test",
			method:                               "GET",
			path:                                 "http://localhost:8080/profile/{2}",
			expectedGetImageLinksByUserId_Images: []models.Image{},
			expectedGetImageLinksByUserId_Error:  nil,
			expectedGetImageLinksByUserId_Count:  0,
			expectedGetProfileIdByUserId_Count:   0,
			expectedGetProfile_Count:             0,
			expectedGetUserIdBySessionId_Error:   errors.New("error"),
			expectedGetUserIdBySessionId_Count:   1,
			expectedStatus:                       http.StatusUnauthorized,
			expectedMessage:                      "user not found\n",
			logger:                               logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageService := sign_up_mocks.NewMockImageService(mockCtrl)
			profileService := sign_up_mocks.NewMockProfileService(mockCtrl)
			userService := sign_up_mocks.NewMockUserService(mockCtrl)
			sessionService := sign_up_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(imageService, profileService, userService, sessionService, tt.logger)

			imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), gomock.Any()).
				Return(tt.expectedGetImageLinksByUserId_Images, tt.expectedGetImageLinksByUserId_Error).
				Times(tt.expectedGetImageLinksByUserId_Count)
			profileService.EXPECT().GetProfile(gomock.Any(), gomock.Any()).
				Return(tt.expectedGetProfile_Profile, tt.expectedGetProfile_Error).
				Times(tt.expectedGetProfile_Count)
			userService.EXPECT().GetProfileIdByUserId(gomock.Any(), gomock.Any()).
				Return(tt.expectedGetProfileIdByUserId_ProfileId, tt.expectedGetProfileIdByUserId_Error).
				Times(tt.expectedGetProfileIdByUserId_Count)
			sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), gomock.Any()).
				Return(tt.expectedGetUserIdBySessionId_UserId, tt.expectedGetUserIdBySessionId_Error).
				Times(tt.expectedGetUserIdBySessionId_Count)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "4gg-4gfd6-445gfdf",
			}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel() // Отменяем контекст после завершения работы
			ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
			req = req.WithContext(ctx)
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
