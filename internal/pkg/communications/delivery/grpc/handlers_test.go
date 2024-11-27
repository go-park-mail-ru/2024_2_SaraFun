package communicationsgrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGrpcCommunicationsHandler_GetReactionList(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedCommunications.GetReactionListRequest
		mockSetup      func(mockUC *mocks.MockReactionUseCase)
		expectedResp   *generatedCommunications.GetReactionListResponse
		expectedErrMsg string
	}{
		{
			name: "Successful GetReactionList",
			request: &generatedCommunications.GetReactionListRequest{
				UserId: 100,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), 100).
					Return([]int{200, 201, 202}, nil)
			},
			expectedResp: &generatedCommunications.GetReactionListResponse{
				Receivers: []int32{200, 201, 202},
			},
			expectedErrMsg: "",
		},
		{
			name: "Error GetReactionList",
			request: &generatedCommunications.GetReactionListRequest{
				UserId: 101,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), 101).
					Return(nil, errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get reaction list error: usecase error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			logger := zap.NewNop()
			handler := NewGrpcCommunicationHandler(mockUC, logger)

			resp, err := handler.GetReactionList(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedResp.Receivers, resp.Receivers)
			}
		})
	}
}

func TestGrpcCommunicationsHandler_GetMatchTime(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedCommunications.GetMatchTimeRequest
		mockSetup      func(mockUC *mocks.MockReactionUseCase)
		expectedResp   *generatedCommunications.GetMatchTimeResponse
		expectedErrMsg string
	}{
		{
			name: "Successful GetMatchTime",
			request: &generatedCommunications.GetMatchTimeRequest{
				FirstUser:  100,
				SecondUser: 200,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchTime(gomock.Any(), 100, 200).
					Return("2024-11-27T12:00:00Z", nil)
			},
			expectedResp: &generatedCommunications.GetMatchTimeResponse{
				Time: "2024-11-27T12:00:00Z",
			},
			expectedErrMsg: "",
		},
		{
			name: "Error GetMatchTime",
			request: &generatedCommunications.GetMatchTimeRequest{
				FirstUser:  101,
				SecondUser: 201,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchTime(gomock.Any(), 101, 201).
					Return("", errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get match time error: usecase error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			logger := zap.NewNop()
			handler := NewGrpcCommunicationHandler(mockUC, logger)

			resp, err := handler.GetMatchTime(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedResp.Time, resp.Time)
			}
		})
	}
}

func TestGrpcCommunicationsHandler_UpdateOrCreateReaction(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedCommunications.UpdateOrCreateReactionRequest
		mockSetup      func(mockUC *mocks.MockReactionUseCase)
		expectedResp   *generatedCommunications.UpdateOrCreateReactionResponse
		expectedErrMsg string
	}{
		{
			name: "Successful UpdateOrCreateReaction",
			request: &generatedCommunications.UpdateOrCreateReactionRequest{
				Reaction: &generatedCommunications.Reaction{
					Author:   100,
					Receiver: 200,
					Type:     true,
				},
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					UpdateOrCreateReaction(gomock.Any(), models.Reaction{
						Author:   100,
						Receiver: 200,
						Type:     true,
					}).
					Return(nil)
			},
			expectedResp:   &generatedCommunications.UpdateOrCreateReactionResponse{},
			expectedErrMsg: "",
		},
		{
			name: "Error UpdateOrCreateReaction",
			request: &generatedCommunications.UpdateOrCreateReactionRequest{
				Reaction: &generatedCommunications.Reaction{
					Author:   101,
					Receiver: 201,
					Type:     false,
				},
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					UpdateOrCreateReaction(gomock.Any(), models.Reaction{
						Author:   101,
						Receiver: 201,
						Type:     false,
					}).
					Return(errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc update reaction error: usecase error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			logger := zap.NewNop()
			handler := NewGrpcCommunicationHandler(mockUC, logger)

			resp, err := handler.UpdateOrCreateReaction(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
