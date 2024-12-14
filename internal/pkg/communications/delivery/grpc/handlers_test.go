package communicationsgrpc

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddReaction(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")

	tests := []struct {
		name     string
		reaction models.Reaction
		addError error
		addCount int
		logger   *zap.Logger
	}{
		{
			name:     "good test",
			reaction: models.Reaction{Receiver: 1, Type: true},
			addError: nil,
			addCount: 1,
			logger:   logger,
		},
		{
			name:     "bad test",
			reaction: models.Reaction{Receiver: 100, Type: false},
			addError: errors.New("test error"),
			addCount: 1,
			logger:   logger,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := mocks.NewMockReactionUseCase(mockCtrl)
			usecase.EXPECT().AddReaction(ctx, tt.reaction).Return(tt.addError).Times(tt.addCount)

			s := NewGrpcCommunicationHandler(usecase, logger)

			reaction := &generatedCommunications.Reaction{
				ID:       int32(tt.reaction.Id),
				Author:   int32(tt.reaction.Author),
				Receiver: int32(tt.reaction.Receiver),
				Type:     tt.reaction.Type,
			}
			req := &generatedCommunications.AddReactionRequest{Reaction: reaction}
			_, err := s.AddReaction(ctx, req)
			require.ErrorIs(t, err, tt.addError)
		})
	}

}

func TestGetMatchList(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")

	tests := []struct {
		name       string
		userId     int
		returnRepo []int
		repoError  error
		repoCount  int
		wantList   []int
		logger     *zap.Logger
	}{
		{
			name:       "good test",
			userId:     1,
			returnRepo: []int{1},
			repoError:  nil,
			wantList:   []int{1},
		},
		{
			name:       "bad test",
			userId:     1,
			returnRepo: []int{1},
			repoError:  errors.New("test error"),
			wantList:   nil,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := mocks.NewMockReactionUseCase(mockCtrl)
			usecase.EXPECT().GetMatchList(ctx, tt.userId).Return(tt.returnRepo, tt.repoError).Times(1)

			s := NewGrpcCommunicationHandler(usecase, logger)
			req := &generatedCommunications.GetMatchListRequest{UserID: int32(tt.userId)}
			list, err := s.GetMatchList(ctx, req)
			if list == nil {
				list = &generatedCommunications.GetMatchListResponse{}
			}
			authors := list.Authors
			require.ErrorIs(t, err, tt.repoError)
			for i, v := range tt.wantList {
				if int32(v) != authors[i] {
					t.Errorf("Bad list result: want %d, got %d", v, authors[i])
				}
			}

		})
	}

}

func TestGetReaction(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")

	tests := []struct {
		name            string
		userId          int
		returnReceivers []int
		returnError     error
		returnCount     int
		wantReceivers   []int
		logger          *zap.Logger
	}{
		{
			name:            "good test",
			userId:          1,
			returnReceivers: []int{1, 3},
			returnError:     nil,
			returnCount:     1,
			wantReceivers:   []int{1, 3},
			logger:          logger,
		},
		{
			name:            "bad test",
			userId:          1,
			returnReceivers: nil,
			returnError:     errors.New("test error"),
			returnCount:     1,
			logger:          logger,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := mocks.NewMockReactionUseCase(mockCtrl)
			usecase.EXPECT().GetReactionList(ctx, tt.userId).Return(tt.returnReceivers, tt.returnError).Times(tt.returnCount)

			s := NewGrpcCommunicationHandler(usecase, logger)
			req := &generatedCommunications.GetReactionListRequest{UserId: int32(tt.userId)}
			receivers, err := s.GetReactionList(ctx, req)
			if err != nil {
				receivers = &generatedCommunications.GetReactionListResponse{}
			}
			recs := receivers.Receivers

			require.ErrorIs(t, err, tt.returnError)
			for i, v := range tt.wantReceivers {
				if int32(v) != recs[i] {
					t.Errorf("Bad reaction list result: want %d, got %d", v, recs[i])
				}
			}
		})
	}
}

func TestGetMatchTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockReactionUseCase(mockCtrl)
	tests := []struct {
		name         string
		firstUser    int
		secondUser   int
		repoReturn   string
		repoError    error
		repoCount    int
		expectedTime string
	}{
		{
			name:         "successfull test",
			firstUser:    1,
			secondUser:   2,
			repoReturn:   time.DateTime,
			repoError:    nil,
			repoCount:    1,
			expectedTime: time.DateTime,
		},
		{
			name:         "bad test",
			firstUser:    1,
			secondUser:   2,
			repoReturn:   "",
			repoError:    errors.New("test error"),
			repoCount:    1,
			expectedTime: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetMatchTime(ctx, tt.firstUser, tt.secondUser).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := NewGrpcCommunicationHandler(repo, logger)
			req := &generatedCommunications.GetMatchTimeRequest{
				FirstUser:  int32(tt.firstUser),
				SecondUser: int32(tt.secondUser),
			}
			time, err := s.GetMatchTime(ctx, req)
			if err != nil {
				time = &generatedCommunications.GetMatchTimeResponse{}
			}
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedTime, time.Time)
		})
	}
}

func TestGetMatchesBySearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockReactionUseCase(mockCtrl)
	tests := []struct {
		name         string
		userId       int
		search       string
		repoReturn   []int
		repoError    error
		repoCount    int
		expectedList []int32
	}{
		{
			name:         "good test",
			userId:       1,
			search:       "sparkit",
			repoReturn:   []int{1, 2, 3},
			repoError:    nil,
			repoCount:    1,
			expectedList: []int32{1, 2, 3},
		},
		{
			name:         "bad test",
			userId:       1,
			search:       "",
			repoReturn:   nil,
			repoError:    errors.New("test error"),
			repoCount:    1,
			expectedList: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetMatchesBySearch(ctx, tt.userId, tt.search).Return(tt.repoReturn, tt.repoError).Times(tt.repoCount)
			s := NewGrpcCommunicationHandler(repo, logger)
			req := &generatedCommunications.GetMatchesBySearchRequest{
				UserID: int32(tt.userId),
				Search: tt.search,
			}
			list, err := s.GetMatchesBySearch(ctx, req)
			if err != nil {
				list = &generatedCommunications.GetMatchesBySearchResponse{}
			}
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedList, list.Authors)
		})
	}
}

func TestUpdateOrCreateReaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockReactionUseCase(mockCtrl)

	tests := []struct {
		name      string
		reaction  models.Reaction
		repoError error
		repoCount int
	}{
		{
			name: "good test",
			reaction: models.Reaction{
				Author:   1,
				Receiver: 2,
				Type:     true,
			},
			repoError: nil,
			repoCount: 1,
		},
		{
			name:      "bad test",
			reaction:  models.Reaction{},
			repoError: errors.New("test error"),
			repoCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().UpdateOrCreateReaction(ctx, tt.reaction).Return(tt.repoError).Times(tt.repoCount)
			s := NewGrpcCommunicationHandler(repo, logger)
			reaction := &generatedCommunications.Reaction{
				ID:       int32(tt.reaction.Id),
				Author:   int32(tt.reaction.Author),
				Receiver: int32(tt.reaction.Receiver),
				Type:     tt.reaction.Type,
			}
			req := &generatedCommunications.UpdateOrCreateReactionRequest{Reaction: reaction}
			_, err := s.UpdateOrCreateReaction(ctx, req)
			require.ErrorIs(t, err, tt.repoError)
		})
	}
}
