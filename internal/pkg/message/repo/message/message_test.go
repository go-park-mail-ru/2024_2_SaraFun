package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock New error: %v", err)
	}
	logger := zap.NewNop()
	defer db.Close()

	repo := New(db, logger)

	tests := []struct {
		name          string
		message       *models.Message
		queryID       int
		queryAuthor   int
		queryError    error
		expectedValue int
		expectedError error
	}{
		{
			name: "successfull test",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "success",
			},
			queryID:       1,
			queryAuthor:   1,
			queryError:    nil,
			expectedValue: 1,
			expectedError: nil,
		},
		{
			name: "bad test",
			message: &models.Message{
				Author:   0,
				Receiver: -3,
				Body:     "",
			},
			queryID:       -1,
			queryAuthor:   -1,
			queryError:    errors.New("bad"),
			expectedValue: -1,
			expectedError: fmt.Errorf("AddMessage error: %w", errors.New("bad")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("INSERT INTO message").WithArgs(tt.message.Author, tt.message.Receiver,
					tt.message.Body).WillReturnRows(mock.NewRows([]string{"id", "author"}).AddRow(tt.queryID, tt.queryAuthor))
			} else {
				mock.ExpectQuery("INSERT INTO message").WithArgs(tt.message.Author, tt.message.Receiver, tt.message.Body).WillReturnError(tt.queryError)
			}
			id, err := repo.AddMessage(ctx, tt.message)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedValue, id)
		})
	}
}

func TestGetLastMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock New error: %v", err)
	}
	logger := zap.NewNop()
	defer db.Close()
	repo := New(db, logger)

	tests := []struct {
		name            string
		authorID        int
		receiverID      int
		queryBody       string
		queryAuthor     int
		queryCreated_at string
		queryError      error
		expectedMessage models.Message
	}{
		{
			name:            "successfull test",
			authorID:        1,
			receiverID:      2,
			queryBody:       "success",
			queryAuthor:     1,
			queryCreated_at: time.DateTime,
			queryError:      nil,
			expectedMessage: models.Message{
				Author: 1,
				Body:   "success",
				Time:   time.DateTime,
			},
		},
		{
			name:            "bad test",
			authorID:        1,
			receiverID:      2,
			queryBody:       "",
			queryAuthor:     1,
			queryCreated_at: time.Now().String(),
			queryError:      errors.New("bad"),
			expectedMessage: models.Message{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT body, author, created_at FROM message").
					WithArgs(tt.authorID, tt.receiverID).
					WillReturnRows(mock.NewRows([]string{"body", "author", "created_at"}).
						AddRow(tt.queryBody, tt.queryAuthor, tt.queryCreated_at))
			} else {
				mock.ExpectQuery("SELECT body, author, created_at FROM message").
					WithArgs(tt.authorID, tt.receiverID).
					WillReturnError(tt.queryError)
			}
			msg, err := repo.GetLastMessage(ctx, tt.authorID, tt.receiverID)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedMessage, msg)
		})
	}
}

func TestGetChatMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock New error: %v", err)
	}
	logger := zap.NewNop()
	defer db.Close()
	repo := New(db, logger)

	successQueryRows := mock.NewRows([]string{"body", "author", "receiver", "created_at"}).
		AddRow("success1", 1, 2, time.DateTime).AddRow("success2", 2, 1, time.DateTime)

	tests := []struct {
		name             string
		firstUserID      int
		secondUserID     int
		queryRows        *sqlmock.Rows
		queryError       error
		expectedMessages []models.Message
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			queryRows:    successQueryRows,
			queryError:   nil,
			expectedMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "success1",
					Time:     time.DateTime,
				},
				{
					Author:   2,
					Receiver: 1,
					Body:     "success2",
					Time:     time.DateTime,
				},
			},
		},
		{
			name:             "bad test",
			firstUserID:      1,
			secondUserID:     2,
			queryRows:        nil,
			queryError:       errors.New("bad"),
			expectedMessages: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT body, author, receiver, created_at FROM message").
					WithArgs(tt.firstUserID, tt.secondUserID).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT body, author, receiver, created_at FROM message").
					WithArgs(tt.firstUserID, tt.secondUserID).WillReturnError(tt.queryError)
			}
			msgs, err := repo.GetChatMessages(ctx, tt.firstUserID, tt.secondUserID)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedMessages, msgs)
		})
	}
}

func TestGetMessagesBySearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock New error: %v", err)
	}
	logger := zap.NewNop()
	defer db.Close()
	repo := New(db, logger)

	successRows := mock.NewRows([]string{"body", "author", "receiver", "created_at"}).
		AddRow("success1", 1, 2, time.DateTime).
		AddRow("success2", 2, 1, time.DateTime)

	tests := []struct {
		name             string
		userID           int
		page             int
		search           string
		queryRows        *sqlmock.Rows
		queryError       error
		expectedMessages []models.Message
	}{
		{
			name:       "successfull test",
			userID:     1,
			page:       1,
			search:     "success",
			queryRows:  successRows,
			queryError: nil,
			expectedMessages: []models.Message{
				{
					Author:   1,
					Receiver: 2,
					Body:     "success1",
					Time:     time.DateTime,
				},
				{
					Author:   2,
					Receiver: 1,
					Body:     "success2",
					Time:     time.DateTime,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT body, author, receiver, created_at FROM message").
					WithArgs(tt.userID, tt.search, tt.page*25, (tt.page-1)*25).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT body, author, receiver, created_at FROM message").
					WithArgs(tt.userID, tt.search, tt.page*25, (tt.page-1)*25).
					WillReturnError(tt.queryError)
			}

			msgs, err := repo.GetMessagesBySearch(ctx, tt.userID, tt.page, tt.search)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedMessages, msgs)
		})
	}
}
