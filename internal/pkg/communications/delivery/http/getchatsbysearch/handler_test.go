package getchatsbysearch

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
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
//	imageservicemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getchatsbysearch/mocks"
//	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
//	messagemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen/mocks"
//	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
//	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//)
//
////nolint:all
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
//	messageClient := messagemocks.NewMockMessageClient(mockCtrl)
//
//	handler := NewHandler(communicationsClient, sessionClient, personalitiesClient, imageService, messageClient, logger)
//
//	validUserID := int32(10)
//	validAuthorID := int32(100)
//	validImages := []models.Image{{Id: 1, Link: "http://example.com/img1.jpg"}}
//
//	validRequest := Request{
//		Search: "john",
//		Page:   1,
//	}
//	validBody, _ := json.Marshal(validRequest)
//
//	tests := []struct {
//		name                     string
//		cookieValue              string
//		userID                   int32
//		userIDError              error
//		searchMatchesAuthors     []int32
//		searchMatchesError       error
//		profileError             error
//		usernameError            error
//		imageError               error
//		lastMessageError         error
//		matchTimeError           error
//		noLastMessage            bool
//		messagesBySearchError    error
//		requestBody              []byte
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//		{
//			name:                     "no cookie",
//			cookieValue:              "",
//			requestBody:              validBody,
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "session user error",
//			cookieValue:              "bad_session",
//			userIDError:              errors.New("session error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "bad json",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			requestBody:              []byte(`{bad json`),
//			expectedStatus:           http.StatusBadRequest,
//			expectedResponseContains: "bad request",
//		},
//		{
//			name:                     "bad get matches",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			searchMatchesError:       errors.New("matches error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get matches",
//		},
//		{
//			name:                     "bad get profile",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			searchMatchesAuthors:     []int32{validAuthorID},
//			profileError:             errors.New("profile error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get profile",
//		},
//		{
//			name:                     "bad get username",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			searchMatchesAuthors:     []int32{validAuthorID},
//			usernameError:            errors.New("username error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get username",
//		},
//		{
//			name:                     "image error",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			searchMatchesAuthors:     []int32{validAuthorID},
//			imageError:               errors.New("image error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "image error",
//		},
//		{
//			name:                     "match time error",
//			cookieValue:              "valid_session",
//			userID:                   validUserID,
//			searchMatchesAuthors:     []int32{validAuthorID},
//			noLastMessage:            true,
//			matchTimeError:           errors.New("match time error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "bad get match time",
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
//			canProceed := tt.userIDError == nil && tt.cookieValue != "" && len(tt.requestBody) > 0 && isValidJSON(tt.requestBody)
//			if canProceed {
//				var req Request
//				_ = json.Unmarshal(tt.requestBody, &req)
//				getMatchesReq := &generatedCommunications.GetMatchesBySearchRequest{UserID: tt.userID, Search: req.Search}
//
//				if tt.searchMatchesError == nil {
//					resp := &generatedCommunications.GetMatchesBySearchResponse{Authors: tt.searchMatchesAuthors}
//					communicationsClient.EXPECT().GetMatchesBySearch(gomock.Any(), getMatchesReq).
//						Return(resp, nil).Times(1)
//				} else {
//					communicationsClient.EXPECT().GetMatchesBySearch(gomock.Any(), getMatchesReq).
//						Return(nil, tt.searchMatchesError).Times(1)
//				}
//			}
//
//			_ = canProceed &&
//				tt.searchMatchesError == nil &&
//				len(tt.searchMatchesAuthors) > 0 &&
//				tt.profileError == nil &&
//				tt.usernameError == nil &&
//				tt.imageError == nil &&
//				tt.lastMessageError == nil &&
//				tt.matchTimeError == nil
//
//			if canProceed && tt.searchMatchesError == nil && len(tt.searchMatchesAuthors) > 0 {
//				author := tt.searchMatchesAuthors[0]
//				getProfileReq := &generatedPersonalities.GetProfileRequest{Id: author}
//				if tt.profileError == nil {
//					profResp := &generatedPersonalities.GetProfileResponse{
//						Profile: &generatedPersonalities.Profile{
//							ID:        author,
//							FirstName: "John",
//							LastName:  "Doe",
//						},
//					}
//					personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//						Return(profResp, nil).Times(1)
//
//					getUsernameReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: author}
//					if tt.usernameError == nil {
//						userResp := &generatedPersonalities.GetUsernameByUserIDResponse{Username: "johndoe"}
//						personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//							Return(userResp, nil).Times(1)
//
//						if tt.imageError == nil {
//							imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(author)).
//								Return(validImages, nil).Times(1)
//
//							getLastReq := &generatedMessage.GetLastMessageRequest{AuthorID: tt.userID, ReceiverID: author}
//							if tt.lastMessageError == nil {
//								if tt.noLastMessage {
//									messageClient.EXPECT().GetLastMessage(gomock.Any(), getLastReq).
//										Return(&generatedMessage.GetLastMessageResponse{Message: ""}, nil).Times(1)
//									getMatchTimeReq := &generatedCommunications.GetMatchTimeRequest{FirstUser: tt.userID, SecondUser: author}
//									if tt.matchTimeError == nil {
//										timeResp := &generatedCommunications.GetMatchTimeResponse{Time: "2024-12-12T10:00:00Z"}
//										communicationsClient.EXPECT().GetMatchTime(gomock.Any(), getMatchTimeReq).
//											Return(timeResp, nil).Times(1)
//									} else {
//										communicationsClient.EXPECT().GetMatchTime(gomock.Any(), getMatchTimeReq).
//											Return(nil, tt.matchTimeError).Times(1)
//									}
//								} else {
//									messageClient.EXPECT().GetLastMessage(gomock.Any(), getLastReq).
//										Return(&generatedMessage.GetLastMessageResponse{
//											Message: "Hello",
//											Self:    true,
//											Time:    "2024-12-12T10:00:00Z",
//										}, nil).Times(1)
//								}
//							} else {
//								messageClient.EXPECT().GetLastMessage(gomock.Any(), getLastReq).
//									Return(nil, tt.lastMessageError).Times(1)
//							}
//						} else {
//							imageService.EXPECT().GetImageLinksByUserId(gomock.Any(), int(author)).
//								Return(nil, tt.imageError).Times(1)
//						}
//					} else {
//						personalitiesClient.EXPECT().GetUsernameByUserID(gomock.Any(), getUsernameReq).
//							Return(nil, tt.usernameError).Times(1)
//					}
//				} else {
//					personalitiesClient.EXPECT().GetProfile(gomock.Any(), getProfileReq).
//						Return(nil, tt.profileError).Times(1)
//				}
//			}
//
//			noErrorsForMessagesSearch := canProceed && tt.searchMatchesError == nil &&
//				tt.profileError == nil && tt.usernameError == nil && tt.imageError == nil && tt.lastMessageError == nil && tt.matchTimeError == nil
//
//			if noErrorsForMessagesSearch {
//				var req Request
//				_ = json.Unmarshal(tt.requestBody, &req)
//				getMsgsReq := &generatedMessage.GetMessagesBySearchRequest{
//					UserID: tt.userID,
//					Page:   int32(req.Page),
//					Search: req.Search,
//				}
//				if tt.messagesBySearchError == nil {
//					msgsResp := &generatedMessage.GetMessagesBySearchResponse{Messages: []*generatedMessage.ChatMessage{}}
//					messageClient.EXPECT().GetMessagesBySearch(gomock.Any(), getMsgsReq).
//						Return(msgsResp, nil).Times(1)
//				} else {
//					messageClient.EXPECT().GetMessagesBySearch(gomock.Any(), getMsgsReq).
//						Return(nil, tt.messagesBySearchError).Times(1)
//				}
//			}
//
//			req := httptest.NewRequest(http.MethodPost, "/search/chats", bytes.NewBuffer(tt.requestBody))
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
//func isValidJSON(b []byte) bool {
//	var js interface{}
//	return json.Unmarshal(b, &js) == nil
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
