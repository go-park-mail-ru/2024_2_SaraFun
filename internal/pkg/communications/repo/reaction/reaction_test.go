package reaction

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddReaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name       string
		reaction   models.Reaction
		queryError error
	}{
		{
			name: "successfull test",
			reaction: models.Reaction{
				Author:   1,
				Receiver: 2,
				Type:     true,
			},
			queryError: nil,
		},
		{
			name:       "bad test",
			reaction:   models.Reaction{},
			queryError: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO reaction").
				WithArgs(tt.reaction.Author, tt.reaction.Receiver, tt.reaction.Type).
				WillReturnResult(sqlmock.NewResult(0, 0)).
				WillReturnError(tt.queryError)
			err := repo.AddReaction(ctx, tt.reaction)
			require.ErrorIs(t, err, tt.queryError)
		})
	}
}

func TestGetMatchList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name            string
		userID          int
		queryRows       *sqlmock.Rows
		queryErr        error
		expectedAuthors []int
	}{
		{
			name:   "successfull test",
			userID: 1,
			queryRows: sqlmock.NewRows([]string{"Author"}).
				AddRow(1).
				AddRow(2).
				AddRow(3),
			queryErr:        nil,
			expectedAuthors: []int{1, 2, 3},
		},
		{
			name:   "bad test",
			userID: 2,
			queryRows: sqlmock.NewRows([]string{"Author"}).
				AddRow(1),
			queryErr:        errors.New("error"),
			expectedAuthors: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT author FROM reaction").
					WithArgs(tt.userID, tt.userID).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT author FROM reaction").
					WithArgs(tt.userID, tt.userID).
					WillReturnError(tt.queryErr)
			}

			authors, err := repo.GetMatchList(ctx, tt.userID)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedAuthors, authors)
		})
	}
}

func TestGetReactionList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	repo := New(db, logger)
	tests := []struct {
		name              string
		userID            int
		queryRows         *sqlmock.Rows
		queryErr          error
		expectedReceivers []int
	}{
		{
			name:   "successfull test",
			userID: 1,
			queryRows: sqlmock.NewRows([]string{"Receiver"}).
				AddRow(3).
				AddRow(2).
				AddRow(4),
			queryErr:          nil,
			expectedReceivers: []int{3, 2, 4},
		},
		{
			name:              "bad test",
			userID:            2,
			queryRows:         nil,
			queryErr:          errors.New("error"),
			expectedReceivers: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT receiver FROM reaction").
					WithArgs(tt.userID).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT receiver FROM reaction").
					WithArgs(tt.userID).
					WillReturnError(tt.queryErr)
			}
			receivers, err := repo.GetReactionList(ctx, tt.userID)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedReceivers, receivers)
		})
	}
}

func TestGetMatchTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name         string
		firstUserID  int
		secondUserID int
		queryRow     *sqlmock.Rows
		queryErr     error
		expectedTime string
	}{
		{
			name:         "successfull test",
			firstUserID:  1,
			secondUserID: 2,
			queryRow: sqlmock.NewRows([]string{"created_at"}).
				AddRow(time.DateTime),
			queryErr:     nil,
			expectedTime: time.DateTime,
		},
		{
			name:         "bad test",
			firstUserID:  2,
			secondUserID: 1,
			queryRow:     nil,
			queryErr:     errors.New("error"),
			expectedTime: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT created_at FROM reaction").
					WithArgs(tt.firstUserID, tt.secondUserID).
					WillReturnRows(tt.queryRow)
			} else {
				mock.ExpectQuery("SELECT created_at FROM reaction").
					WithArgs(tt.firstUserID, tt.secondUserID).
					WillReturnError(tt.queryErr)
			}
			time, err := repo.GetMatchTime(ctx, tt.firstUserID, tt.secondUserID)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedTime, time)
		})
	}
}

func TestGetMatchesByFirstName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name            string
		userID          int
		firstname       string
		queryRow        *sqlmock.Rows
		queryErr        error
		expectedAuthors []int
	}{
		{
			name:      "successfull test",
			userID:    1,
			firstname: "sparkit",
			queryRow: sqlmock.NewRows([]string{"r.author"}).
				AddRow(2).
				AddRow(3),
			queryErr:        nil,
			expectedAuthors: []int{2, 3},
		},
		{
			name:            "bad test",
			userID:          2,
			firstname:       "sparkit",
			queryRow:        nil,
			queryErr:        errors.New("error"),
			expectedAuthors: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.firstname).
					WillReturnRows(tt.queryRow)
			} else {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.firstname).
					WillReturnError(tt.queryErr)
			}
			authors, err := repo.GetMatchesByFirstName(ctx, tt.userID, tt.firstname)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedAuthors, authors)
		})
	}
}

func TestGetMatchesByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name            string
		userID          int
		username        string
		queryRows       *sqlmock.Rows
		queryErr        error
		expectedAuthors []int
	}{
		{
			name:     "successfull test",
			userID:   1,
			username: "sparkit",
			queryRows: sqlmock.NewRows([]string{"r.author"}).
				AddRow(2).
				AddRow(3),
			queryErr:        nil,
			expectedAuthors: []int{2, 3},
		},
		{
			name:            "bad test",
			userID:          2,
			username:        "sparkit",
			queryRows:       nil,
			queryErr:        errors.New("error"),
			expectedAuthors: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.username).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.username).
					WillReturnError(tt.queryErr)
			}
			authors, err := repo.GetMatchesByUsername(ctx, tt.userID, tt.username)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedAuthors, authors)
		})
	}
}

func TestGetMatchesByString(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name            string
		userID          int
		search          string
		queryRows       *sqlmock.Rows
		queryErr        error
		expectedAuthors []int
	}{
		{
			name:   "successfull test",
			userID: 1,
			search: "sparkit",
			queryRows: sqlmock.NewRows([]string{"r.author"}).
				AddRow(2).
				AddRow(3),
			queryErr:        nil,
			expectedAuthors: []int{2, 3},
		},
		{
			name:            "bad test",
			userID:          2,
			search:          "sparkit",
			queryRows:       nil,
			queryErr:        errors.New("error"),
			expectedAuthors: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryErr == nil {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.search).
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.userID, tt.userID, tt.search).
					WillReturnError(tt.queryErr)
			}
			authors, err := repo.GetMatchesByString(ctx, tt.userID, tt.search)
			require.ErrorIs(t, err, tt.queryErr)
			require.Equal(t, tt.expectedAuthors, authors)
		})
	}
}

func TestUpdateOrCreateReaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name         string
		reaction     models.Reaction
		updateResult driver.Result
		updateErr    error
		createError  error
	}{
		{
			name: "successfull test",
			reaction: models.Reaction{
				Author:   1,
				Receiver: 2,
				Type:     true,
			},
			updateResult: sqlmock.NewResult(0, 0),
			updateErr:    nil,
			createError:  nil,
		},
		{
			name: "bad test",
			reaction: models.Reaction{
				Author:   1,
				Receiver: 2,
				Type:     true,
			},
			updateResult: sqlmock.NewResult(0, 0),
			updateErr:    nil,
			createError:  errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createError == nil {
				mock.ExpectExec("UPDATE").
					WithArgs(tt.reaction.Type, tt.reaction.Author, tt.reaction.Receiver).
					WillReturnResult(tt.updateResult)
				mock.ExpectExec("INSERT INTO").
					WithArgs(tt.reaction.Author, tt.reaction.Receiver, tt.reaction.Type).
					WillReturnResult(sqlmock.NewResult(1, 1))

			} else {
				mock.ExpectExec("UPDATE").
					WithArgs(tt.reaction.Type, tt.reaction.Author, tt.reaction.Receiver).
					WillReturnResult(tt.updateResult)
				mock.ExpectExec("INSERT INTO").
					WithArgs(tt.reaction.Author, tt.reaction.Receiver, tt.reaction.Type).
					WillReturnError(tt.createError)
			}
			err := repo.UpdateOrCreateReaction(ctx, tt.reaction)
			require.ErrorIs(t, err, tt.createError)
		})
	}
}
