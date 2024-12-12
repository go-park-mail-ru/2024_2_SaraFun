package message

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/usecase/message/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)

	tests := []struct {
		name          string
		message       *models.Message
		repoMessageId int
		repoErr       error
		repoCount     int
		expectedId    int
	}{
		{
			name: "successfull test",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "привет",
				Time:     time.DateTime,
			},
			repoMessageId: 1,
			repoErr:       nil,
			repoCount:     1,
			expectedId:    1,
		},
		{
			name: "bad test",
			message: &models.Message{
				Author:   2,
				Receiver: 1,
				Body:     "",
				Time:     "",
			},
			repoMessageId: -1,
			repoErr:       errors.New("error"),
			repoCount:     1,
			expectedId:    -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo.EXPECT().AddMessage(ctx, test.message).Return(test.repoMessageId, test.repoErr).Times(test.repoCount)
			usecase := New(repo, logger)
			id, err := usecase.AddMessage(ctx, test.message)
			require.ErrorIs(t, err, test.repoErr)
			require.Equal(t, test.expectedId, id)
		})
	}
}

func TestGetLastMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name             string
		authorID         int
		receiverID       int
		repoMessage      models.Message
		repoErr          error
		repoCount        int
		expectedMessages models.Message
		expectedSelf     bool
	}{
		{
			name:       "successfull test",
			authorID:   1,
			receiverID: 2,
			repoMessage: models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "привет",
				Time:     time.DateTime,
			},
			repoErr:      nil,
			repoCount:    1,
			expectedSelf: true,
			expectedMessages: models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "привет",
				Time:     time.DateTime,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo.EXPECT().GetLastMessage(ctx, test.authorID, test.receiverID).Return(test.repoMessage, test.repoErr).Times(test.repoCount)
			usecase := New(repo, logger)
			message, self, err := usecase.GetLastMessage(ctx, test.authorID, test.receiverID)
			require.ErrorIs(t, err, test.repoErr)
			require.Equal(t, test.expectedSelf, self)
			require.Equal(t, test.expectedMessages, message)
		})
	}
}

func TestGetChatMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name             string
		firstUserID      int
		secondUserID     int
		repoMessages     []models.Message
		repoErr          error
		repoCount        int
		expectedMessages []models.Message
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			repoMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "privet",
					Time:     time.DateTime,
				},
			},
			repoErr:   nil,
			repoCount: 1,
			expectedMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "privet",
					Time:     time.DateTime,
				},
			},
		},
		{
			name:         "bad test",
			firstUserID:  1,
			secondUserID: 2,
			repoMessages: []models.Message{
				{
					Author: 1,
				},
			},
			repoErr:          errors.New("error"),
			repoCount:        1,
			expectedMessages: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo.EXPECT().GetChatMessages(ctx, test.firstUserID, test.secondUserID).Return(test.repoMessages, test.repoErr).Times(test.repoCount)
			usecase := New(repo, logger)
			messages, err := usecase.GetChatMessages(ctx, test.firstUserID, test.secondUserID)
			require.ErrorIs(t, err, test.repoErr)
			require.Equal(t, test.expectedMessages, messages)
		})
	}
}

func TestGetMessagesBySearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name             string
		userID           int
		page             int
		search           string
		repoMessages     []models.Message
		repoErr          error
		repoCount        int
		expectedMessages []models.Message
	}{
		{
			name:   "successfull test",
			userID: 1,
			page:   1,
			search: "privet",
			repoMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "privet",
					Time:     time.DateTime,
				},
			},
			repoErr:   nil,
			repoCount: 1,
			expectedMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "privet",
					Time:     time.DateTime,
				},
			},
		},
		{
			name:             "bad test",
			userID:           1,
			page:             1,
			search:           "privet",
			repoMessages:     nil,
			repoErr:          errors.New("error"),
			repoCount:        1,
			expectedMessages: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo.EXPECT().GetMessagesBySearch(ctx, test.userID, test.page, test.search).Return(test.repoMessages, test.repoErr).Times(test.repoCount)
			usecase := New(repo, logger)
			messages, err := usecase.GetMessagesBySearch(ctx, test.userID, test.page, test.search)
			require.ErrorIs(t, err, test.repoErr)
			require.Equal(t, test.expectedMessages, messages)
		})
	}
}
