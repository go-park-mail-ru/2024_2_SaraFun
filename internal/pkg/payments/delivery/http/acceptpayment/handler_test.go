//nolint:golint
package acceptpayment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"

	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
)

//nolint:all
func TestHandler_Handle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")

	logger := zap.NewNop()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	paymentsClient := paymentsmocks.NewMockPaymentClient(mockCtrl)
	authClient := authmocks.NewMockAuthClient(mockCtrl)

	handler := NewHandler(authClient, paymentsClient, logger)

	tests := []struct {
		name                string
		method              string
		path                string
		body                []byte
		cookieValue         string
		paymentsReturnError error
		expectedStatus      int
		expectedMessage     string
	}{
		{
			name:                "good test",
			method:              http.MethodPost,
			path:                "/acceptpayment",
			body:                []byte(`{"object": {"amount": {"value": "100.00"}, "description": "10"}}`),
			cookieValue:         "sparkit",
			paymentsReturnError: nil,
			expectedStatus:      http.StatusOK,
			expectedMessage:     "",
		},
		{
			name:                "bad json",
			method:              http.MethodPost,
			path:                "/acceptpayment",
			body:                []byte(`{bad json`),
			cookieValue:         "sparkit",
			paymentsReturnError: nil,
			expectedStatus:      http.StatusBadRequest,
			expectedMessage:     "decode json error\n",
		},
		{
			name:                "invalid amount value",
			method:              http.MethodPost,
			path:                "/acceptpayment",
			body:                []byte(`{"object": {"amount": {"value": "invalid_float"}, "description": "10"}}`),
			cookieValue:         "sparkit",
			paymentsReturnError: nil,
			expectedStatus:      http.StatusBadRequest,
			expectedMessage:     "parse json error\n",
		},
		{
			name:                "invalid description",
			method:              http.MethodPost,
			path:                "/acceptpayment",
			body:                []byte(`{"object": {"amount": {"value": "100.00"}, "description": "invalid_int"}}`),
			cookieValue:         "sparkit",
			paymentsReturnError: nil,
			expectedStatus:      http.StatusBadRequest,
			expectedMessage:     "parse json error\n",
		},
		{
			name:                "change balance error",
			method:              http.MethodPost,
			path:                "/acceptpayment",
			body:                []byte(`{"object": {"amount": {"value": "100.00"}, "description": "10"}}`),
			cookieValue:         "sparkit",
			paymentsReturnError: errors.New("change balance failed"),
			expectedStatus:      http.StatusUnauthorized,
			expectedMessage:     "change balance error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.paymentsReturnError != nil {
				var jsonData map[string]interface{}
				err := json.Unmarshal(tt.body, &jsonData)
				if err != nil {
					t.Fatalf("failed to unmarshal test case body: %v", err)
				}

				object, ok := jsonData["object"].(map[string]interface{})
				if !ok {
					t.Fatalf("invalid test case body structure: missing 'object'")
				}
				amount, ok := object["amount"].(map[string]interface{})
				if !ok {
					t.Fatalf("invalid test case body structure: missing 'amount'")
				}
				valueStr, ok := amount["value"].(string)
				if !ok {
					t.Fatalf("invalid test case body structure: 'amount.value' is not a string")
				}
				price, err := strconv.ParseFloat(valueStr, 32)
				if err != nil {
					t.Fatalf("invalid test case body 'amount.value': %v", err)
				}
				descriptionStr, ok := object["description"].(string)
				if !ok {
					t.Fatalf("invalid test case body structure: 'description' is not a string")
				}
				payerID, err := strconv.Atoi(descriptionStr)
				if err != nil {
					t.Fatalf("invalid test case body 'description': %v", err)
				}

				changeBalanceReq := &generatedPayments.ChangeBalanceRequest{
					UserID: int32(payerID),
					Amount: int32(price),
				}
				paymentsClient.EXPECT().ChangeBalance(ctx, changeBalanceReq).
					Return(nil, tt.paymentsReturnError).Times(1)
			} else if tt.expectedStatus == http.StatusOK {
				var jsonData map[string]interface{}
				err := json.Unmarshal(tt.body, &jsonData)
				if err != nil {
					t.Fatalf("failed to unmarshal test case body: %v", err)
				}

				object, ok := jsonData["object"].(map[string]interface{})
				if !ok {
					t.Fatalf("invalid test case body structure: missing 'object'")
				}
				amount, ok := object["amount"].(map[string]interface{})
				if !ok {
					t.Fatalf("invalid test case body structure: missing 'amount'")
				}
				valueStr, ok := amount["value"].(string)
				if !ok {
					t.Fatalf("invalid test case body structure: 'amount.value' is not a string")
				}
				price, err := strconv.ParseFloat(valueStr, 32)
				if err != nil {
					t.Fatalf("invalid test case body 'amount.value': %v", err)
				}
				descriptionStr, ok := object["description"].(string)
				if !ok {
					t.Fatalf("invalid test case body structure: 'description' is not a string")
				}
				payerID, err := strconv.Atoi(descriptionStr)
				if err != nil {
					t.Fatalf("invalid test case body 'description': %v", err)
				}

				changeBalanceReq := &generatedPayments.ChangeBalanceRequest{
					UserID: int32(payerID),
					Amount: int32(price),
				}
				paymentsClient.EXPECT().ChangeBalance(ctx, changeBalanceReq).
					Return(&generatedPayments.ChangeBalanceResponse{}, nil).Times(1)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
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
