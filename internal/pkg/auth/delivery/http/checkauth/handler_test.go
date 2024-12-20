package checkauth

//
//import (
//	"context"
//	"errors"
//	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"
//
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
//	paymentClient := paymentsmocks.NewMockPaymentClient(mockCtrl)
//	handler := NewHandler(sessionClient, paymentClient, logger)
//
//	tests := []struct {
//		name                     string
//		method                   string
//		cookieValue              string
//		checkSessionError        error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//		{
//			name:                     "not GET method",
//			method:                   http.MethodPost,
//			expectedStatus:           http.StatusMethodNotAllowed,
//			expectedResponseContains: "method is not allowed",
//		},
//		{
//			name:                     "no cookie",
//			method:                   http.MethodGet,
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "bad session",
//			method:                   http.MethodGet,
//			cookieValue:              "bad_session",
//			checkSessionError:        errors.New("session error"),
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "bad session",
//		},
//		{
//			name:                     "good test",
//			method:                   http.MethodGet,
//			cookieValue:              "valid_session",
//			expectedStatus:           http.StatusOK,
//			expectedResponseContains: "ok",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.cookieValue != "" && tt.checkSessionError != nil {
//				sessionReq := &generatedAuth.CheckSessionRequest{SessionID: tt.cookieValue}
//				sessionClient.EXPECT().CheckSession(gomock.Any(), sessionReq).
//					Return(nil, tt.checkSessionError).Times(1)
//			} else if tt.cookieValue != "" && tt.checkSessionError == nil {
//				sessionReq := &generatedAuth.CheckSessionRequest{SessionID: tt.cookieValue}
//				sessionResp := &generatedAuth.CheckSessionResponse{}
//				sessionClient.EXPECT().CheckSession(gomock.Any(), sessionReq).
//					Return(sessionResp, nil).Times(1)
//			}
//
//			req := httptest.NewRequest(tt.method, "/checkauth", nil)
//			req = req.WithContext(ctx)
//			if tt.cookieValue != "" {
//				req.AddCookie(&http.Cookie{
//					Name:  consts.SessionCookie,
//					Value: tt.cookieValue,
//				})
//			}
//			w := httptest.NewRecorder()
//
//			handler.Handle(w, req)
//
//			if w.Code != tt.expectedStatus {
//				t.Errorf("%s: handler returned wrong status code: got %v, want %v", tt.name, w.Code, tt.expectedStatus)
//			}
//			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
//				t.Errorf("%s: handler returned unexpected body: got %v, want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
//			}
//		})
//	}
//}
//
//func contains(s, substr string) bool {
//	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
//		(len(s) > 0 && len(substr) > 0 && string(s[0:len(substr)]) == substr) ||
//		(len(s) > len(substr) && string(s[len(s)-len(substr):]) == substr) ||
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
