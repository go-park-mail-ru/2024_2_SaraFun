// internal/pkg/communications/delivery/grpc/handlers_test.go
package communicationsgrpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGrpcCommunicationsHandler_AddReaction(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedCommunications.AddReactionRequest
		mockSetup      func(mockUC *mocks.MockReactionUseCase)
		expectedResp   *generatedCommunications.AddReactionResponse
		expectedErrMsg string
	}{
		{
			name: "Successful AddReaction",
			request: &generatedCommunications.AddReactionRequest{
				Reaction: &generatedCommunications.Reaction{
					ID:       1,
					Author:   100,
					Receiver: 200,
					Type:     true,
				},
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					AddReaction(gomock.Any(), models.Reaction{
						Id:       1,
						Author:   100,
						Receiver: 200,
						Type:     true,
					}).
					Return(nil)
			},
			expectedResp:   &generatedCommunications.AddReactionResponse{},
			expectedErrMsg: "",
		},
		{
			name: "UseCase AddReaction Error",
			request: &generatedCommunications.AddReactionRequest{
				Reaction: &generatedCommunications.Reaction{
					ID:       2,
					Author:   101,
					Receiver: 201,
					Type:     false,
				},
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					AddReaction(gomock.Any(), models.Reaction{
						Id:       2,
						Author:   101,
						Receiver: 201,
						Type:     false,
					}).
					Return(errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc add reaction error: usecase error",
		},
		{
			name: "Invalid Reaction Type",
			request: &generatedCommunications.AddReactionRequest{
				Reaction: &generatedCommunications.Reaction{
					ID:       3,
					Author:   102,
					Receiver: 202,
				},
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					AddReaction(gomock.Any(), models.Reaction{
						Id:       3,
						Author:   102,
						Receiver: 202,
					}).
					Return(errors.New("invalid reaction type"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc add reaction error: invalid reaction type",
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной цикла
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			handler := NewGrpcCommunicationHandler(mockUC)

			resp, err := handler.AddReaction(context.Background(), tt.request)

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

func TestGrpcCommunicationsHandler_GetMatchList(t *testing.T) {
	tests := []struct {
		name           string
		request        *generatedCommunications.GetMatchListRequest
		mockSetup      func(mockUC *mocks.MockReactionUseCase)
		expectedResp   *generatedCommunications.GetMatchListResponse
		expectedErrMsg string
	}{
		{
			name: "Successful GetMatchList",
			request: &generatedCommunications.GetMatchListRequest{
				UserID: 100,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchList(gomock.Any(), 100).
					Return([]int{200, 201, 202}, nil)
			},
			expectedResp: &generatedCommunications.GetMatchListResponse{
				Authors: []int32{200, 201, 202},
			},
			expectedErrMsg: "",
		},
		{
			name: "UseCase GetMatchList Error",
			request: &generatedCommunications.GetMatchListRequest{
				UserID: 101,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchList(gomock.Any(), 101).
					Return(nil, errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get match list error: usecase error",
		},
		{
			name: "No Matches Found",
			request: &generatedCommunications.GetMatchListRequest{
				UserID: 102,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchList(gomock.Any(), 102).
					Return([]int{}, nil)
			},
			expectedResp: &generatedCommunications.GetMatchListResponse{
				Authors: []int32{},
			},
			expectedErrMsg: "",
		},
		{
			name: "Invalid User ID",
			request: &generatedCommunications.GetMatchListRequest{
				UserID: -1,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetMatchList(gomock.Any(), -1).
					Return(nil, fmt.Errorf("invalid user ID: %d", -1))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get match list error: invalid user ID: -1",
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной цикла
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			handler := NewGrpcCommunicationHandler(mockUC)

			resp, err := handler.GetMatchList(context.Background(), tt.request)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedResp.Authors, resp.Authors)
			}
		})
	}
}

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
				UserId: 200,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), 200).
					Return([]int{300, 301, 302}, nil)
			},
			expectedResp: &generatedCommunications.GetReactionListResponse{
				Receivers: []int32{300, 301, 302},
			},
			expectedErrMsg: "",
		},
		{
			name: "UseCase GetReactionList Error",
			request: &generatedCommunications.GetReactionListRequest{
				UserId: 201,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), 201).
					Return(nil, errors.New("usecase error"))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get reaction list error: usecase error",
		},
		{
			name: "No Reactions Found",
			request: &generatedCommunications.GetReactionListRequest{
				UserId: 202,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), 202).
					Return([]int{}, nil)
			},
			expectedResp: &generatedCommunications.GetReactionListResponse{
				Receivers: []int32{},
			},
			expectedErrMsg: "",
		},
		{
			name: "Invalid User ID",
			request: &generatedCommunications.GetReactionListRequest{
				UserId: -2,
			},
			mockSetup: func(mockUC *mocks.MockReactionUseCase) {
				mockUC.EXPECT().
					GetReactionList(gomock.Any(), -2).
					Return(nil, fmt.Errorf("invalid user ID: %d", -2))
			},
			expectedResp:   nil,
			expectedErrMsg: "grpc get reaction list error: invalid user ID: -2",
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной цикла
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockReactionUseCase(ctrl)
			tt.mockSetup(mockUC)

			handler := NewGrpcCommunicationHandler(mockUC)

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
