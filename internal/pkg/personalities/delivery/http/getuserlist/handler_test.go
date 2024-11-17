package getuserlist

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getuserlist/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestGetUserListHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gfctx")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name               string
		method             string
		path               string
		UserBySessionId    int
		UserBySessionErr   error
		UserBySessionCount int
		GetProfile         models.Profile
		GetProfileError    error
		GetProfileCount    int
		UsernameByID       string
		UsernameByIDError  error
		UsernameByIDCount  int
		GetUserList        []models.User
		GetUserListError   error
		GetUserListCount   int
		GetImages          []models.Image
		GetImagesError     error
		GetImagesCount     int
		GetReactions       []int
		GetReactionsError  error
		GetReactionsCount  int
		expectedStatus     int
		expectedMessage    string
		cookieValue        string
		logger             *zap.Logger
	}{
		{
			name:               "successfull test",
			method:             http.MethodGet,
			path:               "/users",
			UserBySessionId:    1,
			UserBySessionErr:   nil,
			UserBySessionCount: 1,
			GetProfile:         models.Profile{FirstName: "Kirill"},
			GetProfileError:    nil,
			GetProfileCount:    1,
			UsernameByID:       "username",
			UsernameByIDError:  nil,
			UsernameByIDCount:  1,
			GetUserList:        []models.User{{ID: 2, Username: "Andrey"}},
			GetUserListError:   nil,
			GetUserListCount:   1,
			GetImages:          []models.Image{{Id: 1, Link: "link"}},
			GetImagesError:     nil,
			GetImagesCount:     1,
			GetReactions:       []int{2},
			GetReactionsError:  nil,
			GetReactionsCount:  1,
			expectedStatus:     http.StatusOK,
			expectedMessage:    "[{\"user\":2,\"username\":\"username\",\"profile\":{\"id\":0,\"first_name\":\"Kirill\"},\"images\":[{\"id\":1,\"link\":\"link\"}]}]",
			logger:             logger,
		},
		{
			name:               "bad test",
			method:             http.MethodGet,
			path:               "/users",
			UserBySessionId:    1,
			UserBySessionErr:   nil,
			UserBySessionCount: 1,
			GetProfile:         models.Profile{FirstName: "Kirill"},
			GetProfileError:    errors.New("error"),
			GetProfileCount:    1,
			UsernameByID:       "username",
			UsernameByIDError:  nil,
			UsernameByIDCount:  0,
			GetUserList:        []models.User{{ID: 2, Username: "Andrey"}},
			GetUserListError:   nil,
			GetUserListCount:   1,
			GetImages:          []models.Image{{Id: 1, Link: "link"}},
			GetImagesError:     nil,
			GetImagesCount:     0,
			GetReactions:       []int{2},
			GetReactionsError:  nil,
			GetReactionsCount:  1,
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "bad get profile\n",
			logger:             logger,
		},
		{
			name:               "bad test",
			method:             http.MethodGet,
			path:               "/users",
			UserBySessionId:    1,
			UserBySessionErr:   nil,
			UserBySessionCount: 1,
			GetProfile:         models.Profile{FirstName: "Kirill"},
			GetProfileError:    nil,
			GetProfileCount:    1,
			UsernameByID:       "username",
			UsernameByIDError:  nil,
			UsernameByIDCount:  0,
			GetUserList:        []models.User{{ID: 2, Username: "Andrey"}},
			GetUserListError:   nil,
			GetUserListCount:   1,
			GetImages:          []models.Image{{Id: 1, Link: "link"}},
			GetImagesError:     errors.New("error"),
			GetImagesCount:     1,
			GetReactions:       []int{2},
			GetReactionsError:  nil,
			GetReactionsCount:  1,
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "error\n",
			logger:             logger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionService := getuserlist_mocks.NewMockSessionService(mockCtrl)
			profileService := getuserlist_mocks.NewMockProfileService(mockCtrl)
			userService := getuserlist_mocks.NewMockUserService(mockCtrl)
			imageService := getuserlist_mocks.NewMockImageService(mockCtrl)
			reactionService := getuserlist_mocks.NewMockReactionService(mockCtrl)

			sessionService.EXPECT().GetUserIDBySessionID(ctx, gomock.Any()).
				Return(tt.UserBySessionId, tt.UserBySessionErr).Times(tt.UserBySessionCount)
			imageService.EXPECT().GetImageLinksByUserId(ctx, gomock.Any()).
				Return(tt.GetImages, tt.GetImagesError).Times(tt.GetImagesCount)
			reactionService.EXPECT().GetReactionList(ctx, tt.UserBySessionId).
				Return(tt.GetReactions, tt.GetReactionsError).Times(tt.GetReactionsCount)
			userService.EXPECT().GetFeedList(ctx, gomock.Any(), gomock.Any()).
				Return(tt.GetUserList, tt.GetUserListError).Times(tt.GetUserListCount)
			for _, user := range tt.GetUserList {
				profileService.EXPECT().GetProfile(ctx, user.ID).
					Return(tt.GetProfile, tt.GetProfileError).Times(tt.GetProfileCount)
				userService.EXPECT().GetUsernameByUserId(ctx, user.ID).
					Return(tt.UsernameByID, tt.UsernameByIDError).Times(tt.UsernameByIDCount)
			}

			handler := NewHandler(sessionService, profileService, userService, imageService, reactionService, tt.logger)

			req := httptest.NewRequest(tt.method, tt.path, nil).WithContext(ctx)

			req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
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
