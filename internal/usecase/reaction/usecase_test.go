package reaction

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/usecase/reaction/mocks"
	"sparkit/internal/utils/consts"
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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().AddReaction(ctx, tt.reaction).Return(tt.addError).Times(tt.addCount)

			s := New(repo, logger)
			err := s.AddReaction(ctx, tt.reaction)
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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetMatchList(ctx, tt.userId).Return(tt.returnRepo, tt.repoError).Times(1)

			s := New(repo, logger)
			list, err := s.GetMatchList(ctx, tt.userId)
			require.ErrorIs(t, err, tt.repoError)
			for i, v := range tt.wantList {
				if v != list[i] {
					t.Errorf("Bad list result: want %d, got %d", v, list[i])
				}
			}

		})
	}

}

func TestGetReaction(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")

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
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetReactionList(ctx, tt.userId).Return(tt.returnReceivers, tt.returnError).Times(tt.returnCount)

			s := New(repo, logger)

			receivers, err := s.GetReactionList(ctx, tt.userId)

			require.ErrorIs(t, err, tt.returnError)
			for i, v := range tt.wantReceivers {
				if v != receivers[i] {
					t.Errorf("Bad reaction list result: want %d, got %d", v, receivers[i])
				}
			}
		})
	}
}
