// internal/pkg/auth/delivery/grpc/handlers_test.go
package authgrpc

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateSession(t *testing.T) {
	// Определяем тестовые случаи
	tests := []struct {
		name           string
		request        *generatedAuth.CreateSessionRequest
		mockSetup      func(mockUseCase *mocks.MockUseCase)
		expectedResp   *generatedAuth.CreateSessionResponse
		expectedErrMsg string
	}{
		{
			name: "Successful CreateSession",
			request: &generatedAuth.CreateSessionRequest{
				User: &generatedAuth.User{
					ID: 1,
				},
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				// Настройка ожидания вызова CreateSession с определёнными аргументами
				mockUseCase.EXPECT().
					CreateSession(gomock.Any(), models.User{ID: 1}).
					Return(models.Session{
						SessionID: "session123",
						UserID:    1,
						CreatedAt: time.Now(),
						ExpiresAt: time.Now().Add(24 * time.Hour),
					}, nil)
			},
			expectedResp:   nil,
			expectedErrMsg: "",
		},
		{
			name: "UseCase CreateSession Error",
			request: &generatedAuth.CreateSessionRequest{
				User: &generatedAuth.User{
					ID: 2,
				},
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					CreateSession(gomock.Any(), models.User{ID: 2}).
					Return(models.Session{}, errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "usecase error",
		},
		{
			name: "Invalid User ID",
			request: &generatedAuth.CreateSessionRequest{
				User: &generatedAuth.User{
					ID: -1,
				},
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					CreateSession(gomock.Any(), models.User{ID: -1}).
					Return(models.Session{}, errors.New("invalid user ID"))
			},
			expectedResp:   nil,
			expectedErrMsg: "invalid user ID",
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной цикла
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockUseCase(ctrl)

			tt.mockSetup(mockUseCase)

			logger := zap.NewNop()

			handler := NewGRPCAuthHandler(mockUseCase, logger)

			resp, err := handler.CreateSession(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {

				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.Session)
				require.Equal(t, "session123", resp.Session.SessionID)
				require.Equal(t, int32(1), resp.Session.UserID)
				require.NotEmpty(t, resp.Session.CreatedAt)
				require.NotEmpty(t, resp.Session.ExpiresAt)

				_, err := time.Parse(time.RFC3339, resp.Session.CreatedAt)
				require.NoError(t, err)
				_, err = time.Parse(time.RFC3339, resp.Session.ExpiresAt)
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedAuth.DeleteSessionRequest
		mockSetup      func(mockUseCase *mocks.MockUseCase)
		expectedResp   *generatedAuth.DeleteSessionResponse
		expectedErrMsg string
	}{
		{
			name: "Successful DeleteSession",
			request: &generatedAuth.DeleteSessionRequest{
				SessionID: "session123",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					DeleteSession(gomock.Any(), "session123").
					Return(nil)
			},
			expectedResp:   &generatedAuth.DeleteSessionResponse{},
			expectedErrMsg: "",
		},
		{
			name: "UseCase DeleteSession Error",
			request: &generatedAuth.DeleteSessionRequest{
				SessionID: "session456",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					DeleteSession(gomock.Any(), "session456").
					Return(errors.New("delete error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "delete error",
		},
		{
			name: "Empty SessionID",
			request: &generatedAuth.DeleteSessionRequest{
				SessionID: "",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					DeleteSession(gomock.Any(), "").
					Return(errors.New("session ID cannot be empty"))
			},
			expectedResp:   nil,
			expectedErrMsg: "session ID cannot be empty",
		},
		{
			name: "Successful deletion of empty SessionID",
			request: &generatedAuth.DeleteSessionRequest{
				SessionID: "",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					DeleteSession(gomock.Any(), "").
					Return(nil)
			},
			expectedResp:   &generatedAuth.DeleteSessionResponse{},
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockUseCase(ctrl)

			tt.mockSetup(mockUseCase)

			logger := zap.NewNop()

			handler := NewGRPCAuthHandler(mockUseCase, logger)

			resp, err := handler.DeleteSession(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestCheckSession(t *testing.T) {

	tests := []struct {
		name           string
		request        *generatedAuth.CheckSessionRequest
		mockSetup      func(mockUseCase *mocks.MockUseCase)
		expectedResp   *generatedAuth.CheckSessionResponse
		expectedErrMsg string
	}{
		{
			name: "Successful CheckSession",
			request: &generatedAuth.CheckSessionRequest{
				SessionID: "session123",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					CheckSession(gomock.Any(), "session123").
					Return(nil)
			},
			expectedResp:   &generatedAuth.CheckSessionResponse{},
			expectedErrMsg: "",
		},
		{
			name: "UseCase CheckSession Error",
			request: &generatedAuth.CheckSessionRequest{
				SessionID: "session456",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					CheckSession(gomock.Any(), "session456").
					Return(errors.New("invalid session"))
			},
			expectedResp:   nil,
			expectedErrMsg: "invalid session",
		},
		{
			name: "Empty SessionID",
			request: &generatedAuth.CheckSessionRequest{
				SessionID: "",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					CheckSession(gomock.Any(), "").
					Return(errors.New("session ID cannot be empty"))
			},
			expectedResp:   nil,
			expectedErrMsg: "session ID cannot be empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockUseCase(ctrl)

			tt.mockSetup(mockUseCase)

			logger := zap.NewNop()

			handler := NewGRPCAuthHandler(mockUseCase, logger)

			resp, err := handler.CheckSession(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestGetUserIDBySessionID(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedAuth.GetUserIDBySessionIDRequest
		mockSetup      func(mockUseCase *mocks.MockUseCase)
		expectedResp   *generatedAuth.GetUserIDBYSessionIDResponse
		expectedErrMsg string
	}{
		{
			name: "Successful GetUserIDBySessionID",
			request: &generatedAuth.GetUserIDBySessionIDRequest{
				SessionID: "session123",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					GetUserIDBySessionID(gomock.Any(), "session123").
					Return(42, nil)
			},
			expectedResp: &generatedAuth.GetUserIDBYSessionIDResponse{
				UserId: 42,
			},
			expectedErrMsg: "",
		},
		{
			name: "UseCase GetUserIDBySessionID Error",
			request: &generatedAuth.GetUserIDBySessionIDRequest{
				SessionID: "session456",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					GetUserIDBySessionID(gomock.Any(), "session456").
					Return(0, errors.New("invalid session"))
			},
			expectedResp:   nil,
			expectedErrMsg: "invalid session",
		},
		{
			name: "Invalid SessionID",
			request: &generatedAuth.GetUserIDBySessionIDRequest{
				SessionID: "invalid_session",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					GetUserIDBySessionID(gomock.Any(), "invalid_session").
					Return(-1, fmt.Errorf("convert session id invalid_session to int: strconv.Atoi: parsing \"invalid_session\": invalid syntax"))
			},
			expectedResp:   nil,
			expectedErrMsg: "convert session id invalid_session to int: strconv.Atoi: parsing \"invalid_session\": invalid syntax",
		},
		{
			name: "Empty SessionID",
			request: &generatedAuth.GetUserIDBySessionIDRequest{
				SessionID: "",
			},
			mockSetup: func(mockUseCase *mocks.MockUseCase) {
				mockUseCase.EXPECT().
					GetUserIDBySessionID(gomock.Any(), "").
					Return(0, errors.New("session ID cannot be empty"))
			},
			expectedResp:   nil,
			expectedErrMsg: "session ID cannot be empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockUseCase(ctrl)

			tt.mockSetup(mockUseCase)

			logger := zap.NewNop()

			handler := NewGRPCAuthHandler(mockUseCase, logger)

			resp, err := handler.GetUserIDBySessionID(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedResp.UserId, resp.UserId)
			}
		})
	}
}
