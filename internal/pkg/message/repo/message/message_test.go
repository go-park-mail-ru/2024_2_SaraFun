package message

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStorage_AddMessage(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		message     *models.Message
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedID  int
		expectedErr error
	}{
		{
			name: "Successful AddMessage",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "Hello!",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO message \(author, receiver, body\) VALUES \(\$1, \$2, \$3\) RETURNING id, author`).
					WithArgs(1, 2, "Hello!").
					WillReturnRows(sqlmock.NewRows([]string{"id", "author"}).AddRow(1, 1))
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Database Error",
			message: &models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "Hello!",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO message \(author, receiver, body\) VALUES \(\$1, \$2, \$3\) RETURNING id, author`).
					WithArgs(1, 2, "Hello!").
					WillReturnError(errors.New("database error"))
			},
			expectedID:  -1,
			expectedErr: errors.New("AddMessage error: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			id, err := repo.AddMessage(context.Background(), tt.message)

			require.Equal(t, tt.expectedID, id)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetLastMessage(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		authorID    int
		receiverID  int
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedMsg models.Message
		expectedErr error
	}{

		{
			name:       "No Rows Found",
			authorID:   1,
			receiverID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT body, author, created_at FROM message WHERE`).
					WithArgs(1, 2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedMsg: models.Message{},
			expectedErr: sparkiterrors.ErrNoResult,
		},
		{
			name:       "Database Error",
			authorID:   1,
			receiverID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT body, author, created_at FROM message WHERE`).
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expectedMsg: models.Message{},
			expectedErr: errors.New("GetLastMessage error: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			msg, err := repo.GetLastMessage(context.Background(), tt.authorID, tt.receiverID)

			require.Equal(t, tt.expectedMsg, msg)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetChatMessages(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name         string
		firstUserID  int
		secondUserID int
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedMsgs []models.Message
		expectedErr  error
	}{

		{
			name:         "Database Error",
			firstUserID:  1,
			secondUserID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT body, author, receiver, created_at FROM message WHERE`).
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expectedMsgs: nil,
			expectedErr:  errors.New("GetChatMessages error: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			msgs, err := repo.GetChatMessages(context.Background(), tt.firstUserID, tt.secondUserID)

			require.Equal(t, tt.expectedMsgs, msgs)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetMessagesBySearch(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name         string
		userID       int
		page         int
		search       string
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedMsgs []models.Message
		expectedErr  error
	}{

		{
			name:   "Database Error",
			userID: 1,
			page:   1,
			search: "hello",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT body, author, receiver, created_at FROM message WHERE`).
					WithArgs(1, "hello", 25, 0).
					WillReturnError(errors.New("database error"))
			},
			expectedMsgs: nil,
			expectedErr:  errors.New("GetMessagesBySearch error: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := New(db, logger)
			msgs, err := repo.GetMessagesBySearch(context.Background(), tt.userID, tt.page, tt.search)

			require.Equal(t, tt.expectedMsgs, msgs)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
