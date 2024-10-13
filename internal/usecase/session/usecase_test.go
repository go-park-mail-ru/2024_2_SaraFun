package session

import (
	"context"
	"github.com/golang/mock/gomock"
	sparkiterrors "sparkit/internal/errors"
	"sparkit/internal/usecase/session/mocks"
	"testing"
)

//func TestCreateSession(t *testing.T) {
//	session1 := models.Session{SessionID: uuid.New().String(),
//		UserID:    1,
//		CreatedAt: time.Now()}
//	session2 := models.Session{}
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	tests := []struct {
//		name        string
//		user        models.User
//		session     models.Session
//		wantErr     error
//		wantSession models.Session
//	}{
//		{
//			name:    "successful test",
//			user:    models.User{ID: 1},
//			session: session1,
//			wantErr: nil,
//			wantSession: models.Session{SessionID: uuid.New().String(),
//				UserID:    1,
//				CreatedAt: time.Now()},
//		},
//		{
//			name:        "bad test",
//			user:        models.User{ID: 2},
//			session:     session2,
//			wantErr:     errors.New("failed to create session"),
//			wantSession: models.Session{},
//		},
//	}
//	repo := mocks.NewMockRepository(ctrl)
//	repo.EXPECT().AddSession(gomock.Any(), session1).Return(nil)
//	repo.EXPECT().AddSession(gomock.Any(), session2).Return(sparkiterrors.ErrInvalidSession)
//	s := New(repo)
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := s.CreateSession(context.Background(), tt.user)
//			if err != tt.wantErr {
//				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestGetUserIDBySessionID(t *testing.T) {
	tests := []struct {
		name             string
		sessionID        string
		getUserError     error
		getUserCallCount int
		want             int
		wantErr          error
	}{
		{
			name:             "successfull test",
			sessionID:        "12345",
			want:             1,
			getUserError:     nil,
			getUserCallCount: 1,
			wantErr:          nil,
		},
		{
			name:             "bad test",
			sessionID:        "12342",
			getUserError:     sparkiterrors.ErrInvalidSession,
			getUserCallCount: 1,
			want:             0,
			wantErr:          sparkiterrors.ErrInvalidSession,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			//repo.EXPECT().GetUserIDBySessionID(gomock.Any(), "12345").Return(1, nil).AnyTimes()
			//repo.EXPECT().GetUserIDBySessionID(gomock.Any(), "12342").Return(0, sparkiterrors.ErrInvalidSession).AnyTimes()
			repo.EXPECT().GetUserIDBySessionID(gomock.Any(), tt.sessionID).Return(1, tt.getUserError).Times(tt.getUserCallCount)

			s := New(repo)
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
	tests := []struct {
		name                  string
		sessionID             string
		checkSessionError     error
		checkSessionCallCount int
		wantErr               error
	}{
		{
			name:                  "successfull test",
			sessionID:             "12345",
			checkSessionError:     nil,
			checkSessionCallCount: 1,
			wantErr:               nil,
		},
		{
			name:                  "bad test",
			sessionID:             "12342",
			checkSessionError:     sparkiterrors.ErrInvalidSession,
			checkSessionCallCount: 1,
			wantErr:               sparkiterrors.ErrInvalidSession,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)
			repo.EXPECT().CheckSession(gomock.Any(), tt.sessionID).Return(tt.checkSessionError).Times(tt.checkSessionCallCount)
			s := New(repo)
			err := s.CheckSession(context.Background(), tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("CheckSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name                   string
		sessionID              string
		deleteSessionError     error
		deleteSessionCallCount int
		wantErr                error
	}{
		{
			name:                   "successfull test",
			sessionID:              "12345",
			deleteSessionError:     nil,
			deleteSessionCallCount: 1,
			wantErr:                nil,
		},
		{
			name:                   "bad test",
			sessionID:              "12342",
			deleteSessionError:     sparkiterrors.ErrInvalidSession,
			deleteSessionCallCount: 1,
			wantErr:                sparkiterrors.ErrInvalidSession,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(mockCtrl)

			repo.EXPECT().DeleteSession(gomock.Any(), tt.sessionID).Return(tt.deleteSessionError).Times(tt.deleteSessionCallCount)

			s := New(repo)
			err := s.DeleteSession(context.Background(), tt.sessionID)
			if err != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
