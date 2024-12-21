package addProduct

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"

	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

//nolint:all
func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		path                string
		body                []byte
		cookieValue         string
		authReturnUserID    int32
		authReturnError     error
		paymentsReturnError error
		expectedStatus      int
		expectedMessage     string
	}{

		{
			name:                "bad json",
			method:              http.MethodPost,
			path:                "/addProduct",
			body:                []byte(`{bad json`),
			cookieValue:         "valid_session",
			authReturnUserID:    10,
			authReturnError:     nil,
			paymentsReturnError: nil,
			expectedStatus:      http.StatusBadRequest,
			expectedMessage:     "unmarshal data\n",
		},

		{
			name:                "bad cookie",
			method:              http.MethodPost,
			path:                "/addProduct",
			body:                []byte(`{"Title": "Product1", "Description": "A great product", "ImageLink": "http://example.com/image.jpg", "Price": 100}`),
			cookieValue:         "invalid_session",
			authReturnUserID:    -1,
			authReturnError:     errors.New("invalid session"),
			paymentsReturnError: nil,
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "get user id by session id\n",
		},

		{
			name:                "no session cookie",
			method:              http.MethodPost,
			path:                "/addProduct",
			body:                []byte(`{"Title": "Product1", "Description": "A great product", "ImageLink": "http://example.com/image.jpg", "Price": 100}`),
			cookieValue:         "",
			authReturnUserID:    0,
			authReturnError:     nil,
			paymentsReturnError: nil,
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "bad cookie\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authClient := authmocks.NewMockAuthClient(mockCtrl)
			paymentsClient := paymentsmocks.NewMockPaymentClient(mockCtrl)

			logger := zap.NewNop()
			handler := NewHandler(authClient, paymentsClient, logger)

			if tt.cookieValue != "" {
				getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
				if tt.authReturnError != nil {
					authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
						Return(nil, tt.authReturnError).Times(1)
				} else {
					authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDReq).
						Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: tt.authReturnUserID}, nil).Times(1)
				}

				var jsonData map[string]interface{}
				err := json.Unmarshal(tt.body, &jsonData)
				if err == nil {
					object, ok := jsonData["object"].(map[string]interface{})
					if ok {
						amount, ok := object["amount"].(map[string]interface{})
						if ok {
							value, ok := amount["value"]
							if ok {
								var price int32
								switch v := value.(type) {
								case string:
									priceFloat, err := strconv.ParseFloat(v, 32)
									if err == nil {
										price = int32(priceFloat)
									}
								case float64:
									price = int32(v)
								default:

								}
								descriptionStr, ok := object["description"].(string)
								if ok {
									_, err := strconv.Atoi(descriptionStr)
									if err == nil {
										createProductReq := &generatedPayments.CreateProductRequest{
											Product: &generatedPayments.Product{
												Title:       jsonData["Title"].(string),
												Description: jsonData["Description"].(string),
												ImageLink:   jsonData["ImageLink"].(string),
												Price:       price,
											},
										}
										if tt.paymentsReturnError != nil {
											paymentsClient.EXPECT().CreateProduct(gomock.Any(), createProductReq).
												Return(nil, tt.paymentsReturnError).Times(1)
										} else {
											paymentsClient.EXPECT().CreateProduct(gomock.Any(), createProductReq).
												Return(&generatedPayments.CreateProductResponse{ID: 1}, nil).Times(1)
										}
									}
								}
							}
						}
					}
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req = req.WithContext(context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id"))
			if tt.cookieValue != "" {
				cookie := &http.Cookie{
					Name:  consts.SessionCookie,
					Value: tt.cookieValue,
				}
				req.AddCookie(cookie)
			}

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
