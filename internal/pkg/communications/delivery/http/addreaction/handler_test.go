package addreaction

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	addreaction_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/addreaction/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_CommunicationsClient.go -package=addreaction_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen CommunicationsClient
//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=addreaction_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name                     string
		cookie                   *http.Cookie
		requestBody              interface{}
		mockAuthClient           func(mock *addreaction_mocks.MockAuthClient)
		mockCommunicationsClient func(mock *addreaction_mocks.MockCommunicationsClient)
		expectedStatus           int
		expectedResponse         string
	}{
		{
			name: "Successful Add Reaction",
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session-id",
			},
			requestBody: models.Reaction{
				Receiver: 2,
				Type:     true, // Type is now a string
			},
			mockAuthClient: func(mock *addreaction_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "valid-session-id",
					}).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *addreaction_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					AddReaction(gomock.Any(), &generatedCommunications.AddReactionRequest{
						Reaction: &generatedCommunications.Reaction{
							ID:       0,
							Author:   1,
							Receiver: 2,
							Type:     true, // Type is a string
						},
					}).
					Return(&generatedCommunications.AddReactionResponse{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "ok",
		},
		{
			name:                     "Missing Session Cookie",
			cookie:                   nil,
			requestBody:              nil,
			mockAuthClient:           func(mock *addreaction_mocks.MockAuthClient) {},
			mockCommunicationsClient: func(mock *addreaction_mocks.MockCommunicationsClient) {},
			expectedStatus:           http.StatusUnauthorized,
			expectedResponse:         "session not found\n",
		},
		{
			name: "Error Getting User ID",
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "invalid-session-id",
			},
			requestBody: nil,
			mockAuthClient: func(mock *addreaction_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "invalid-session-id",
					}).
					Return(nil, errors.New("session not found"))
			},
			mockCommunicationsClient: func(mock *addreaction_mocks.MockCommunicationsClient) {},
			expectedStatus:           http.StatusUnauthorized,
			expectedResponse:         "session not found\n",
		},
		{
			name: "Invalid JSON",
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session-id",
			},
			requestBody: "invalid json",
			mockAuthClient: func(mock *addreaction_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *addreaction_mocks.MockCommunicationsClient) {},
			expectedStatus:           http.StatusBadRequest,
			expectedResponse:         "bad request\n",
		},
		{
			name: "Error Adding Reaction",
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session-id",
			},
			requestBody: models.Reaction{
				Receiver: 2,
				Type:     true,
			},
			mockAuthClient: func(mock *addreaction_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *addreaction_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					AddReaction(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mock clients
			mockAuthClient := addreaction_mocks.NewMockAuthClient(ctrl)
			mockCommunicationsClient := addreaction_mocks.NewMockCommunicationsClient(ctrl)

			// Setup mocks
			tt.mockAuthClient(mockAuthClient)
			tt.mockCommunicationsClient(mockCommunicationsClient)

			// Create logger
			logger := zap.NewNop()

			// Create handler
			handler := NewHandler(mockCommunicationsClient, mockAuthClient, logger)

			// Create HTTP request
			var req *http.Request
			if tt.requestBody != nil {
				var bodyBytes []byte
				switch v := tt.requestBody.(type) {
				case string:
					bodyBytes = []byte(v)
				default:
					bodyBytes, _ = json.Marshal(v)
				}
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/", nil)
			}

			// Add cookie if present
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			// Add context with RequestID
			ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)

			// Create ResponseRecorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Handle(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			// Check response body
			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("expected response body %q, got %q", tt.expectedResponse, rr.Body.String())
			}
		})
	}
}
