package getmatches

//
//import (
//	"bytes"
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
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
//	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
//	communicationsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen/mocks"
//	imageservicemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getmatches/mocks"
//	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
//	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//)
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
//	communicationsClient := communicationsmocks.NewMockCommunicationsClient(mockCtrl)
//	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
//	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
//	imageService := imageservicemocks.NewMockImageService(mockCtrl)
//
//	handler := NewHandler(communicationsClient, sessionClient, personalitiesClient, imageService, logger)
//
//	validUserID := int32(10)
//	validAuthorID := int32(100)
//	validImages := []models.Image{{Id: 1, Link: "http://example.com/img1.jpg"}}
//
//	tests := []struct {
//		name                     string
//		cookieValue              string
//		userID                   int32
//		userIDError              error
//		matchListAuthors         []int32
//		matchListError           error
//		profileError             error
//		imageError               error
//		usernameError            error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//
//		{
//			name:                     "no cookie",
//			cookieValue:              "",
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "session user error",
//			cookieValue:              "bad_session",
//			userIDError:              errors.New("session error"),
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "match list error",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			matchListError:           errors.New("match error"),
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//
//		{
//			name:                     "image error",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			matchListAuthors:         []int32{validAuthorID},
//			imageError:               errors.New("image error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "image error",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.cookieValue != "" {
//				getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
//				if tt.userIDError == nil {
//					userResp := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: tt.userID}
//					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
//						Return(userResp, nil).Times(1)
//				} else {
//					sessionClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
//						Return(nil, tt.userIDError).Times(1)
//				}
//			}
//
//			canProceed := (tt.userIDError == nil && tt.cookieValue != "")
//			if canProceed {
//				getMatchListReq := &generatedCommunications.GetMatchListRequest{UserID: tt.userID}
//				if tt.matchListError == nil {
//					resp := &generatedCommunications.GetMatchListResponse{Authors: tt.matchListAuthors}
//					communicationsClient.EXPECT().GetMatchList(gomock.Any(), getMatchListReq).
//						Return(resp, nil).Times(1)
//				} else {
//					communicationsClient.EXPECT().GetMatchList(gomock.Any(), getMatchListReq).
//						Return(nil, tt.matchListError).Times(1)
//				}
//			}
//
//			if canProceed && tt.matchListError == nil && len(tt.matchListAuthors) > 0 {
//				author := tt.matchListAuthors[0]
//				getProfileReq := &generatedPersonalities.GetProfileRequest{Id: author}
//				if tt.profileError == nil {
//					profileResp := &generatedPersonalities.GetProfileResponse{
//						Profile: &generatedPersonalities.Profile{
//							ID:        author,
//							FirstName: "John",
//							LastName:  "Doe",
//							Age:       25,
//							Gender:    "male",
//							Target:    "friendship",
//							About:     "Hello",
//						},
//					}
//					personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//						Return(profileResp, nil).Times(1)
//
//					if tt.imageError == nil {
//						imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(author)).
//							Return(validImages, nil).Times(1)
//					} else {
//						imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(author)).
//							Return(nil, tt.imageError).Times(1)
//					}
//
//					if tt.imageError == nil {
//						getUsernameReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: author}
//						if tt.usernameError == nil {
//							userResp := &generatedPersonalities.GetUsernameByUserIDResponse{Username: "johndoe"}
//							personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//								Return(userResp, nil).Times(1)
//						} else {
//							personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//								Return(nil, tt.usernameError).Times(1)
//						}
//					}
//				} else {
//					personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//						Return(nil, tt.profileError).Times(1)
//				}
//			}
//
//			req := httptest.NewRequest(http.MethodGet, "/matches", bytes.NewBuffer(nil))
//			if tt.cookieValue != "" {
//				req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
//			}
//			w := httptest.NewRecorder()
//
//			handler.Handle(w, req.WithContext(ctx))
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
