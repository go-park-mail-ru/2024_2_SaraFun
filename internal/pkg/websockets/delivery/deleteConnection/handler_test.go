package deleteConnection

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/delivery/deleteConnection/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase := mocks.NewMockUseCase(ctrl)
	authClient := mocks.NewMockAuthClient(ctrl)
	logger := zap.NewNop()

	h := NewHandler(useCase, authClient, logger)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")

	tests := []struct {
		name            string
		cookieValue     string
		mockAuthSetup   func()
		mockUCSetup     func()
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "no cookie",
			cookieValue:     "",
			mockAuthSetup:   func() {},
			mockUCSetup:     func() {},
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "cookie not found\n",
		},

		{
			name:        "usecase error",
			cookieValue: "valid_session",
			mockAuthSetup: func() {
				authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
					SessionID: "valid_session",
				}).Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 123}, nil)
			},
			mockUCSetup: func() {
				useCase.EXPECT().DeleteConnection(gomock.Any(), 123).Return(errors.New("uc error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "add connection error\n",
		},
		{
			name:        "success",
			cookieValue: "valid_session",
			mockAuthSetup: func() {
				authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
					SessionID: "valid_session",
				}).Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 123}, nil)
			},
			mockUCSetup: func() {
				useCase.EXPECT().DeleteConnection(gomock.Any(), 123).Return(nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuthSetup()
			tt.mockUCSetup()

			req := httptest.NewRequest(http.MethodGet, "/deleteConnection", nil)
			req = req.WithContext(ctx)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: tt.cookieValue,
				})
			}
			w := httptest.NewRecorder()

			h.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code mismatch: got %v, want %v", w.Code, tt.expectedStatus)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("body mismatch: got %v, want %v", w.Body.String(), tt.expectedMessage)
			}
		})
	}
}
