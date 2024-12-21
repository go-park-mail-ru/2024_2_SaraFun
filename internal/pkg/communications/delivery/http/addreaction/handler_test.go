package addreaction

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
//	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
//	communicationsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen/mocks"
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
//	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
//	reactionClient := communicationsmocks.NewMockCommunicationsClient(mockCtrl)
//	personalitiesClient :=
//
//	handler := NewHandler(reactionClient, sessionClient, logger)
//
//	validBody := []byte(`{
//		"id":1,
//		"receiver":2,
//		"type":"like"
//	}`)
//
//	tests := []struct {
//		name                     string
//		method                   string
//		cookieValue              string
//		userID                   int32
//		userIDError              error
//		requestBody              []byte
//		addReactionError         error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//
//		{
//			name:                     "no cookie",
//			method:                   http.MethodPost,
//			cookieValue:              "",
//			requestBody:              validBody,
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "session user error",
//			method:                   http.MethodPost,
//			cookieValue:              "bad_session",
//			userIDError:              errors.New("session error"),
//			requestBody:              validBody,
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:           "bad json",
//			method:         http.MethodPost,
//			cookieValue:    "valid_session",
//			userID:         10,
//			requestBody:    []byte(`{bad json`),
//			expectedStatus: http.StatusBadRequest,
//			// "bad request" возвращается при ошибке парсинга
//			expectedResponseContains: "bad request",
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
//			if tt.userIDError == nil && tt.cookieValue != "" && bytes.HasPrefix(tt.requestBody, []byte(`{`)) && !bytes.HasPrefix(tt.requestBody, []byte(`{bad`)) {
//				addReactionReq := gomock.Any()
//				if tt.addReactionError == nil {
//					reactionClient.EXPECT().AddReaction(gomock.Any(), addReactionReq).
//						Return(&generatedCommunications.AddReactionResponse{}, nil).Times(1)
//				} else {
//					reactionClient.EXPECT().AddReaction(gomock.Any(), addReactionReq).
//						Return(nil, tt.addReactionError).Times(1)
//				}
//			}
//
//			req := httptest.NewRequest(tt.method, "/reaction", bytes.NewBuffer(tt.requestBody))
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
