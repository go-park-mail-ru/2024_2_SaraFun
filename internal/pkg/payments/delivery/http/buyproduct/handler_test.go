package buyproduct

//
//import (
//	"bytes"
//	"context"
//
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"go.uber.org/zap"
//
//	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//
//	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
//	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"
//
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//)
//
//// TestHandler_Handle тестирует функцию Handle хэндлера buyproduct.Handler
//func TestHandler_Handle(t *testing.T) {
//	tests := []struct {
//		name                    string
//		method                  string
//		path                    string
//		body                    []byte
//		cookieValue             string
//		authReturnUserID        int32
//		authReturnError         error
//		paymentsReturnError     error
//		expectedStatus          int
//		expectedMessage         string
//		expectedBuyLikesRequest *generatedPayments.BuyLikesRequest
//	}{
//		{
//			name:                "good test",
//			method:              http.MethodPost,
//			path:                "/buyproduct",
//			body:                []byte(`{"title": "Product1", "price": 100}`),
//			cookieValue:         "sparkit",
//			authReturnUserID:    10,
//			authReturnError:     nil,
//			paymentsReturnError: nil,
//			expectedStatus:      http.StatusOK,
//			expectedMessage:     "ok",
//			expectedBuyLikesRequest: &generatedPayments.BuyLikesRequest{
//				Title:  "Product1",
//				Amount: 100,
//				UserID: 10,
//			},
//		},
//		{
//			name:                    "bad json",
//			method:                  http.MethodPost,
//			path:                    "/buyproduct",
//			body:                    []byte(`{bad json`),
//			cookieValue:             "sparkit",
//			authReturnUserID:        10,
//			authReturnError:         nil,
//			paymentsReturnError:     nil,
//			expectedStatus:          http.StatusBadRequest,
//			expectedMessage:         "json unmarshal error\n",
//			expectedBuyLikesRequest: nil,
//		},
//		{
//			name:                "missing fields",
//			method:              http.MethodPost,
//			path:                "/buyproduct",
//			body:                []byte(`{"title": "Product1"}`),
//			cookieValue:         "sparkit",
//			authReturnUserID:    10,
//			authReturnError:     nil,
//			paymentsReturnError: nil,
//			expectedStatus:      http.StatusOK,
//			expectedMessage:     "ok",
//			expectedBuyLikesRequest: &generatedPayments.BuyLikesRequest{
//				Title:  "Product1",
//				Amount: 0,
//				UserID: 10,
//			},
//		},
//		{
//			name:                    "invalid price type",
//			method:                  http.MethodPost,
//			path:                    "/buyproduct",
//			body:                    []byte(`{"title": "Product1", "price": "100"}`),
//			cookieValue:             "sparkit",
//			authReturnUserID:        10,
//			authReturnError:         nil,
//			paymentsReturnError:     nil,
//			expectedStatus:          http.StatusBadRequest,
//			expectedMessage:         "json unmarshal error\n",
//			expectedBuyLikesRequest: nil,
//		},
//
//		{
//			name:                    "no session cookie",
//			method:                  http.MethodPost,
//			path:                    "/buyproduct",
//			body:                    []byte(`{"title": "Product4", "price": 100}`),
//			cookieValue:             "",
//			authReturnUserID:        0,
//			authReturnError:         nil,
//			paymentsReturnError:     nil,
//			expectedStatus:          http.StatusUnauthorized,
//			expectedMessage:         "bad cookie\n",
//			expectedBuyLikesRequest: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			mockCtrl := gomock.NewController(t)
//			defer mockCtrl.Finish()
//
//			authClient := authmocks.NewMockAuthClient(mockCtrl)
//			paymentsClient := paymentsmocks.NewMockPaymentClient(mockCtrl)
//
//			logger := zap.NewNop()
//			handler := NewHandler(authClient, paymentsClient, logger)
//
//			if tt.cookieValue != "" {
//				getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
//				if tt.authReturnError != nil {
//					authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
//						Return(nil, tt.authReturnError).Times(1)
//				} else {
//					authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
//						Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: tt.authReturnUserID}, nil).Times(1)
//				}
//
//				if tt.expectedBuyLikesRequest != nil {
//					paymentsClient.EXPECT().BuyLikes(gomock.Any(), gomock.Eq(tt.expectedBuyLikesRequest)).
//						Return(&generatedPayments.BuyLikesResponse{}, tt.paymentsReturnError).Times(1)
//				}
//			}
//
//			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
//			req = req.WithContext(context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id"))
//			if tt.cookieValue != "" {
//				cookie := &http.Cookie{
//					Name:  consts.SessionCookie,
//					Value: tt.cookieValue,
//				}
//				req.AddCookie(cookie)
//			}
//
//			w := httptest.NewRecorder()
//
//			handler.Handle(w, req)
//
//			if w.Code != tt.expectedStatus {
//				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
//			}
//
//			if w.Body.String() != tt.expectedMessage {
//				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
//			}
//		})
//	}
//}
