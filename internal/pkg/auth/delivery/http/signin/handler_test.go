package signin

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
//	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
//	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//	"github.com/mailru/easyjson"
//)
//
////nolint:all
//func TestHandler(t *testing.T) {
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	userClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
//	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
//
//	handler := NewHandler(userClient, sessionClient, logger)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")
//
//	validUser := models.User{
//		Username: "john_doe",
//		Password: "secret",
//	}
//	validBody, _ := easyjson.Marshal(validUser)
//
//	tests := []struct {
//		name                     string
//		method                   string
//		body                     []byte
//		checkPasswordError       error
//		createSessionError       error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//		{
//			name:                     "not POST method",
//			method:                   http.MethodGet,
//			body:                     validBody,
//			expectedStatus:           http.StatusMethodNotAllowed,
//			expectedResponseContains: "Method not allowed",
//		},
//		{
//			name:                     "bad json",
//			method:                   http.MethodPost,
//			body:                     []byte(`{bad json`),
//			expectedStatus:           http.StatusBadRequest,
//			expectedResponseContains: "Неверный формат данных",
//		},
//		{
//			name:                     "wrong credentials",
//			method:                   http.MethodPost,
//			body:                     validBody,
//			checkPasswordError:       errors.New("wrong creds"),
//			expectedStatus:           http.StatusPreconditionFailed,
//			expectedResponseContains: "wrong credentials",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			req := httptest.NewRequest(tt.method, "/signin", bytes.NewBuffer(tt.body))
//			req = req.WithContext(ctx)
//
//			w := httptest.NewRecorder()
//
//			if tt.method == http.MethodPost && isEasyJSONValidUser(tt.body) {
//
//				var u models.User
//				_ = easyjson.Unmarshal(tt.body, &u)
//
//				cpReq := &generatedPersonalities.CheckPasswordRequest{
//					Username: u.Username,
//					Password: u.Password,
//				}
//
//				if tt.checkPasswordError == nil {
//
//					resp := &generatedPersonalities.CheckPasswordResponse{
//						User: &generatedPersonalities.User{
//							ID:       10,
//							Username: u.Username,
//							Email:    "john@example.com",
//							Profile:  100,
//						},
//					}
//					userClient.EXPECT().CheckPassword(gomock.Any(), cpReq).
//						Return(resp, nil).Times(1)
//
//					if tt.createSessionError == nil {
//
//						csReq := &generatedAuth.CreateSessionRequest{
//							User: &generatedAuth.User{
//								ID:       10,
//								Username: u.Username,
//								Email:    "john@example.com",
//								Password: u.Password,
//								Profile:  100,
//							},
//						}
//						sessionResp := &generatedAuth.CreateSessionResponse{
//							Session: &generatedAuth.Session{SessionID: "valid_session_id"},
//						}
//						sessionClient.EXPECT().CreateSession(gomock.Any(), csReq).
//							Return(sessionResp, nil).Times(1)
//					} else {
//
//						csReq := &generatedAuth.CreateSessionRequest{
//							User: &generatedAuth.User{
//								ID:       10,
//								Username: u.Username,
//								Email:    "john@example.com",
//								Password: u.Password,
//								Profile:  100,
//							},
//						}
//						sessionClient.EXPECT().CreateSession(gomock.Any(), csReq).
//							Return(nil, tt.createSessionError).Times(1)
//					}
//
//				} else {
//
//					userClient.EXPECT().CheckPassword(gomock.Any(), cpReq).
//						Return(nil, tt.checkPasswordError).Times(1)
//				}
//			}
//
//			handler.Handle(w, req)
//
//			if w.Code != tt.expectedStatus {
//				t.Errorf("%s: wrong status code: got %v, want %v", tt.name, w.Code, tt.expectedStatus)
//			}
//			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
//				t.Errorf("%s: wrong body: got %v, want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
//			}
//
//			if tt.expectedStatus == http.StatusOK && tt.createSessionError == nil && tt.checkPasswordError == nil && tt.method == http.MethodPost && isEasyJSONValidUser(tt.body) {
//				resp := w.Result()
//				defer resp.Body.Close()
//				cookies := resp.Cookies()
//				found := false
//				for _, c := range cookies {
//					if c.Name == consts.SessionCookie && c.Value == "valid_session_id" {
//						found = true
//						break
//					}
//				}
//				if !found {
//					t.Errorf("%s: expected session cookie to be set", tt.name)
//				}
//			}
//		})
//	}
//}
//
//func isEasyJSONValidUser(b []byte) bool {
//	var u models.User
//	return easyjson.Unmarshal(b, &u) == nil
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
