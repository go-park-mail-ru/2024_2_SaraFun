package getProducts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"

	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

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
		expectedResponse    *Response
	}{
		{
			name:                "good test",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "valid_session",
			authReturnUserID:    10,
			authReturnError:     nil,
			paymentsReturnError: nil,
			expectedStatus:      http.StatusOK,
			expectedMessage:     "",
			expectedResponse: &Response{
				Responses: []models.Product{
					{
						Title:       "Product1",
						Description: "A great product",
						ImageLink:   "http://example.com/image1.jpg",
						Price:       100,
					},
					{
						Title:       "Product2",
						Description: "Another great product",
						ImageLink:   "http://example.com/image2.jpg",
						Price:       200,
					},
				},
			},
		},
		{
			name:                "no session cookie",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "",
			authReturnUserID:    0,
			authReturnError:     nil,
			paymentsReturnError: nil,
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "bad cookie\n",
			expectedResponse:    nil,
		},
		{
			name:                "invalid session cookie",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "invalid_session",
			authReturnUserID:    0,
			authReturnError:     errors.New("invalid session"),
			paymentsReturnError: nil,
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "get user id by session id\n",
			expectedResponse:    nil,
		},
		{
			name:                "error getting user ID",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "valid_session",
			authReturnUserID:    0,
			authReturnError:     errors.New("some auth error"),
			paymentsReturnError: nil,
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "get user id by session id\n",
			expectedResponse:    nil,
		},
		{
			name:                "error getting products",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "valid_session",
			authReturnUserID:    10,
			authReturnError:     nil,
			paymentsReturnError: errors.New("some payment error"),
			expectedStatus:      http.StatusInternalServerError,
			expectedMessage:     "get products error\n",
			expectedResponse:    nil,
		},
		{
			name:                "error marshalling JSON",
			method:              http.MethodGet,
			path:                "/getproducts",
			body:                nil,
			cookieValue:         "valid_session",
			authReturnUserID:    10,
			authReturnError:     nil,
			paymentsReturnError: nil,
			expectedStatus:      http.StatusOK,
			expectedMessage:     "",
			expectedResponse:    &Response{},
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

					if tt.paymentsReturnError != nil || tt.expectedResponse != nil {
						getProductsReq := &generatedPayments.GetProductsRequest{}
						if tt.paymentsReturnError != nil {
							paymentsClient.EXPECT().GetProducts(gomock.Any(), getProductsReq).
								Return(nil, tt.paymentsReturnError).Times(1)
						} else if tt.expectedResponse != nil {

							productsResponse := &generatedPayments.GetProductsResponse{
								Products: []*generatedPayments.Product{
									{
										Title:       "Product1",
										Description: "A great product",
										ImageLink:   "http://example.com/image1.jpg",
										Price:       100,
									},
									{
										Title:       "Product2",
										Description: "Another great product",
										ImageLink:   "http://example.com/image2.jpg",
										Price:       200,
									},
								},
							}
							paymentsClient.EXPECT().GetProducts(gomock.Any(), getProductsReq).
								Return(productsResponse, nil).Times(1)
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

			if tt.expectedResponse != nil && tt.name != "error marshalling JSON" {

				var resp Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("failed to unmarshal response JSON: %v", err)
				}
				if len(resp.Responses) != len(tt.expectedResponse.Responses) {
					t.Errorf("handler returned unexpected number of products: got %d want %d", len(resp.Responses), len(tt.expectedResponse.Responses))
				} else {
					for i, product := range resp.Responses {
						expectedProduct := tt.expectedResponse.Responses[i]
						if product.Title != expectedProduct.Title ||
							product.Description != expectedProduct.Description ||
							product.ImageLink != expectedProduct.ImageLink ||
							product.Price != expectedProduct.Price {
							t.Errorf("handler returned unexpected product at index %d: got %+v want %+v", i, product, expectedProduct)
						}
					}
				}
			} else if tt.expectedMessage != "" {

				if w.Body.String() != tt.expectedMessage {
					t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
				}
			}
		})
	}
}
