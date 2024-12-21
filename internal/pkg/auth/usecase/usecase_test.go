package usecase

import (
	"context"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/usecase/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGetUserIDBySessionID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	tests := []struct {
		name             string
		sessionID        string
		getUserError     error
		getUserCallCount int
		logger           *zap.Logger
		want             int
		wantErr          error
	}{
		{
			name:             "successfull test",
			sessionID:        "12345",
			want:             1,
			getUserError:     nil,
			getUserCallCount: 1,
			logger:           logger,
			wantErr:          nil,
		},
		{
			name:             "bad test",
			sessionID:        "12342",
			getUserError:     sparkiterrors.ErrInvalidSession,
			getUserCallCount: 1,
			logger:           logger,
			want:             0,
			wantErr:          sparkiterrors.ErrInvalidSession,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().GetUserIDBySessionID(ctx, tt.sessionID).Return(1, tt.getUserError).Times(tt.getUserCallCount)

			s := New(repo, logger)
			res, err := s.GetUserIDBySessionID(ctx, tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("GetUserIDBySessionID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res != tt.want {
				t.Errorf("GetUserIDBySessionID() got = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestCheckSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	tests := []struct {
		name                  string
		sessionID             string
		checkSessionError     error
		checkSessionCallCount int
		logger                *zap.Logger
		wantErr               error
	}{
		{
			name:                  "successfull test",
			sessionID:             "12345",
			checkSessionError:     nil,
			checkSessionCallCount: 1,
			logger:                logger,
			wantErr:               nil,
		},
		{
			name:                  "bad test",
			sessionID:             "12342",
			checkSessionError:     sparkiterrors.ErrInvalidSession,
			checkSessionCallCount: 1,
			logger:                logger,
			wantErr:               sparkiterrors.ErrInvalidSession,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().CheckSession(ctx, tt.sessionID).Return(tt.checkSessionError).Times(tt.checkSessionCallCount)
			s := New(repo, logger)
			err := s.CheckSession(ctx, tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("CheckSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDeleteSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст после завершения работы
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	tests := []struct {
		name                   string
		sessionID              string
		deleteSessionError     error
		deleteSessionCallCount int
		logger                 *zap.Logger
		wantErr                error
	}{
		{
			name:                   "successfull test",
			sessionID:              "12345",
			deleteSessionError:     nil,
			deleteSessionCallCount: 1,
			logger:                 logger,
			wantErr:                nil,
		},
		{
			name:                   "bad test",
			sessionID:              "12342",
			deleteSessionError:     sparkiterrors.ErrInvalidSession,
			deleteSessionCallCount: 1,
			logger:                 logger,
			wantErr:                sparkiterrors.ErrInvalidSession,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)

			repo.EXPECT().DeleteSession(ctx, tt.sessionID).Return(tt.deleteSessionError).Times(tt.deleteSessionCallCount)

			s := New(repo, logger)
			err := s.DeleteSession(ctx, tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
