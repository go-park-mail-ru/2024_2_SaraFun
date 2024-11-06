package getmatches

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	getmatches_mocks "sparkit/internal/handlers/getmatches/mocks"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                         string
		method                       string
		path                         string
		body                         []byte
		GetMatchList_List            []int
		GetMatchList_Err             error
		GetUserIDBySessionID_Id      int
		GetUserIDBYSessionID_Err     error
		GetProfile_Profiles          []models.Profile
		GetProfile_Err               error
		GetUsernameByUserId_Username string
		GetUsernameByUserId_Err      error
		GetMatchListCount            int
		GetProfileCount              int
		GetUsernameByUserIdCount     int
		GetUserIdBySessionIDCount    int
		GetImageLinksArr             []models.Image
		GetImageLinksErr             error
		GetImageLinksCount           int
		expectedStatus               int
		expectedMessage              string
		logger                       *zap.Logger
	}{
		{
			name:                         "successfull test",
			method:                       "GET",
			path:                         "http://localhost:8080/matches",
			GetMatchList_List:            []int{1},
			GetMatchList_Err:             nil,
			GetUserIDBySessionID_Id:      1,
			GetUserIDBYSessionID_Err:     nil,
			GetProfile_Profiles:          []models.Profile{{FirstName: "1"}},
			GetProfile_Err:               nil,
			GetUsernameByUserId_Username: "username",
			GetUsernameByUserId_Err:      nil,
			GetMatchListCount:            1,
			GetProfileCount:              1,
			GetUserIdBySessionIDCount:    1,
			GetUsernameByUserIdCount:     1,
			GetImageLinksArr:             []models.Image{{Id: 1, Link: "link1"}, {Id: 2, Link: "link2"}, {Id: 3, Link: "link3"}},
			GetImageLinksErr:             nil,
			GetImageLinksCount:           1,
			expectedMessage:              "[{\"user\":1,\"username\":\"username\",\"profile\":{\"id\":0,\"first_name\":\"1\"},\"images\":[{\"id\":1,\"link\":\"link1\"},{\"id\":2,\"link\":\"link2\"},{\"id\":3,\"link\":\"link3\"}]}]",
			expectedStatus:               http.StatusOK,
			logger:                       logger,
		},
		{
			name:                         "bad test",
			method:                       "GET",
			path:                         "http://localhost:8080/matches",
			GetMatchList_List:            []int{1},
			GetMatchList_Err:             nil,
			GetUserIDBySessionID_Id:      1,
			GetUserIDBYSessionID_Err:     errors.New("ERROR"),
			GetProfile_Profiles:          []models.Profile{{FirstName: "1"}},
			GetProfile_Err:               nil,
			GetUsernameByUserId_Username: "username",
			GetUsernameByUserId_Err:      nil,
			GetMatchListCount:            0,
			GetProfileCount:              0,
			GetUserIdBySessionIDCount:    1,
			GetUsernameByUserIdCount:     0,
			GetImageLinksCount:           0,
			expectedMessage:              "session not found\n",
			expectedStatus:               http.StatusUnauthorized,
			logger:                       logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reactionService := getmatches_mocks.NewMockReactionService(mockCtrl)
			sessionService := getmatches_mocks.NewMockSessionService(mockCtrl)
			profileService := getmatches_mocks.NewMockProfileService(mockCtrl)
			userService := getmatches_mocks.NewMockUserService(mockCtrl)
			imageService := getmatches_mocks.NewMockImageService(mockCtrl)

			handler := NewHandler(reactionService, sessionService, profileService, userService, imageService, logger)

			reactionService.EXPECT().GetMatchList(gomock.Any(), gomock.Any()).
				Return(tt.GetMatchList_List, tt.GetMatchList_Err).
				Times(tt.GetMatchListCount)
			sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), gomock.Any()).
				Return(tt.GetUserIDBySessionID_Id, tt.GetUserIDBYSessionID_Err).
				Times(tt.GetUserIdBySessionIDCount)
			//profileService.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(tt)
			for i, userId := range tt.GetMatchList_List {
				profileService.EXPECT().GetProfile(gomock.Any(), userId).
					Return(tt.GetProfile_Profiles[i], tt.GetProfile_Err).
					Times(tt.GetProfileCount)
				imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), userId).Return(tt.GetImageLinksArr, tt.GetImageLinksErr).
					Times(tt.GetImageLinksCount)
			}
			userService.EXPECT().GetUsernameByUserId(gomock.Any(), gomock.Any()).
				Return(tt.GetUsernameByUserId_Username, tt.GetUsernameByUserId_Err).
				Times(tt.GetUsernameByUserIdCount)

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
