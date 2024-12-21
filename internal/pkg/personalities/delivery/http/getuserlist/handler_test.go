package getuserlist

//
//import (
//	"context"
//	"errors"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"go.uber.org/zap"
//
//	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
//	communicationsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen/mocks"
//	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
//	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
//
//	imageservicemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getuserlist/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
//)
//
//// 1. Успешный сценарий получения списка пользователей.
//// 2. Отсутствие cookie, приводящее к ошибке авторизации.
//// 3. Ошибка при получении userID по сессии.
//// 4. Ошибка при получении списка реакций.
//// 5. Ошибка при получении списка пользователей (feed).
//// 6. Ошибка при получении профиля пользователя.
//// 7. Ошибка при получении изображений пользователя.
//// 8. Ошибка при получении имени пользователя (username).
//
//func TestHandler(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")
//
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
//	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
//	imageService := imageservicemocks.NewMockImageService(mockCtrl)
//	communicationsClient := communicationsmocks.NewMockCommunicationsClient(mockCtrl)
//
//	handler := NewHandler(sessionClient, personalitiesClient, imageService, communicationsClient, logger)
//
//	tests := []struct {
//		name                     string
//		method                   string
//		cookieValue              string
//		authUserID               int32
//		authError                error
//		reactionListError        error
//		feedListError            error
//		feedUsers                []*generatedPersonalities.User
//		profileError             error
//		imageError               error
//		usernameError            error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//		{
//			name:                     "good test",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			feedUsers:                []*generatedPersonalities.User{{ID: 100}},
//			expectedStatus:           http.StatusOK,
//			expectedResponseContains: `"first_name":`,
//		},
//		{
//			name:                     "no cookie",
//			method:                   http.MethodGet,
//			cookieValue:              "",
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "session not found",
//			method:                   http.MethodGet,
//			cookieValue:              "bad_session",
//			authError:                errors.New("session error"),
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "reaction list failed",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			reactionListError:        errors.New("reaction error"),
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "reaction list failed",
//		},
//		{
//			name:                     "feed list error",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			feedListError:            errors.New("feed error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "ошибка в получении списка пользователей",
//		},
//		{
//			name:                     "profile error",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			feedUsers:                []*generatedPersonalities.User{{ID: 100}},
//			profileError:             errors.New("profile error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get profile",
//		},
//		{
//			name:                     "image error",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			feedUsers:                []*generatedPersonalities.User{{ID: 100}},
//			imageError:               errors.New("image error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "image error",
//		},
//		{
//			name:                     "username error",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			authUserID:               10,
//			feedUsers:                []*generatedPersonalities.User{{ID: 100}},
//			usernameError:            errors.New("username error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get username",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.cookieValue != "" {
//				getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
//				if tt.authError == nil {
//					resp := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: tt.authUserID}
//					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).Return(resp, nil).Times(1)
//				} else {
//					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).Return(nil, tt.authError).Times(1)
//				}
//			}
//
//			if tt.authError == nil && tt.cookieValue != "" {
//				getReactionListReq := &generatedCommunications.GetReactionListRequest{UserId: tt.authUserID}
//				if tt.reactionListError == nil {
//					commResp := &generatedCommunications.GetReactionListResponse{Receivers: []int32{}}
//					communicationsClient.EXPECT().GetReactionList(gomock.Any(), getReactionListReq).
//						Return(commResp, nil).Times(1)
//				} else {
//					communicationsClient.EXPECT().GetReactionList(gomock.Any(), getReactionListReq).
//						Return(nil, tt.reactionListError).Times(1)
//				}
//			}
//
//			if tt.authError == nil && tt.reactionListError == nil && tt.cookieValue != "" {
//				// feed list
//				getFeedReq := &generatedPersonalities.GetFeedListRequest{UserID: tt.authUserID, Receivers: []int32{}}
//				if tt.feedListError == nil {
//					usersResp := &generatedPersonalities.GetFeedListResponse{Users: tt.feedUsers}
//					personalitiesClient.EXPECT().GetFeedList(gomock.Any(), getFeedReq).
//						Return(usersResp, nil).Times(1)
//				} else {
//					personalitiesClient.EXPECT().GetFeedList(gomock.Any(), getFeedReq).
//						Return(nil, tt.feedListError).Times(1)
//				}
//			}
//
//			if tt.authError == nil && tt.reactionListError == nil && tt.feedListError == nil && tt.cookieValue != "" {
//				for _, usr := range tt.feedUsers {
//					getProfileReq := &generatedPersonalities.GetProfileRequest{Id: usr.ID}
//					if tt.profileError == nil {
//						profileResp := &generatedPersonalities.GetProfileResponse{
//							Profile: &generatedPersonalities.Profile{
//								ID:        usr.ID,
//								FirstName: "John",
//								LastName:  "Doe",
//								Age:       25,
//								Gender:    "male",
//								Target:    "friendship",
//								About:     "Hi",
//							},
//						}
//						personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//							Return(profileResp, nil).Times(1)
//					} else {
//						personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//							Return(nil, tt.profileError).Times(1)
//						break
//					}
//
//					if tt.profileError == nil {
//						if tt.imageError == nil {
//							imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(usr.ID)).
//								Return([]models.Image{{Link: "http://example.com/img.jpg"}}, nil).Times(1)
//						} else {
//							imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(usr.ID)).
//								Return(nil, tt.imageError).Times(1)
//							break
//						}
//					}
//
//					if tt.profileError == nil && tt.imageError == nil {
//						getUsernameReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: usr.ID}
//						if tt.usernameError == nil {
//							personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//								Return(&generatedPersonalities.GetUsernameByUserIDResponse{Username: "johnny"}, nil).Times(1)
//						} else {
//							personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//								Return(nil, tt.usernameError).Times(1)
//							break
//						}
//					}
//				}
//			}
//
//			req := httptest.NewRequest(tt.method, "/userlist", nil)
//			req = req.WithContext(ctx)
//			if tt.cookieValue != "" {
//				cookie := &http.Cookie{
//					Name:  consts.SessionCookie,
//					Value: tt.cookieValue,
//				}
//				req.AddCookie(cookie)
//			}
//			w := httptest.NewRecorder()
//
//			handler.Handle(w, req)
//
//			if w.Code != tt.expectedStatus {
//				t.Errorf("%s: handler returned wrong status code: got %v want %v", tt.name, w.Code, tt.expectedStatus)
//			}
//			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
//				t.Errorf("%s: handler returned unexpected body: got %v want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
//			}
//		})
//	}
//}
//
//func contains(s, substr string) bool {
//	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
//		(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
//		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
//		(len(substr) > 0 && len(s) > len(substr) && findInString(s, substr)))
//}
//
//func findInString(s, substr string) bool {
//	for i := 0; i+len(substr) <= len(s); i++ {
//		if s[i:i+len(substr)] == substr {
//			return true
//		}
//	}
//	return false
//}
