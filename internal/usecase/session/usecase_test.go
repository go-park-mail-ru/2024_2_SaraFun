package session

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/usecase/session/mocks"
	"testing"
)

func TestGetUserIDBySessionID(t *testing.T) {
	logger := zap.NewNop()
	defer logger.Sync()
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
			repo.EXPECT().GetUserIDBySessionID(gomock.Any(), tt.sessionID).Return(1, tt.getUserError).Times(tt.getUserCallCount)

			s := New(repo, logger)
			res, err := s.GetUserIDBySessionID(context.Background(), tt.sessionID)
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
	logger := zap.NewNop()
	defer logger.Sync()
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
			repo.EXPECT().CheckSession(gomock.Any(), tt.sessionID).Return(tt.checkSessionError).Times(tt.checkSessionCallCount)
			s := New(repo, logger)
			err := s.CheckSession(context.Background(), tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("CheckSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDeleteSession(t *testing.T) {
	logger := zap.NewNop()
	defer logger.Sync()
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

			repo.EXPECT().DeleteSession(gomock.Any(), tt.sessionID).Return(tt.deleteSessionError).Times(tt.deleteSessionCallCount)

			s := New(repo, logger)
			err := s.DeleteSession(context.Background(), tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
