package reaction

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

func TestAddReaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		reaction models.Reaction
		queryErr error
		wantErr  error
	}{
		{
			name:     "successful insert",
			reaction: models.Reaction{Author: 1, Receiver: 2, Type: true},
			queryErr: nil,
			wantErr:  nil,
		},

		{
			name:     "invalid author",
			reaction: models.Reaction{Author: 0, Receiver: 2, Type: true},
			queryErr: nil,
			wantErr:  errors.New("failed to insert reaction: author must be greater than 0"),
		},
		{
			name:     "invalid receiver",
			reaction: models.Reaction{Author: 1, Receiver: 0, Type: true},
			queryErr: nil,
			wantErr:  errors.New("failed to insert reaction: receiver must be greater than 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.reaction.Author <= 0 {

				require.Error(t, storage.AddReaction(ctx, tt.reaction))
				return
			}
			if tt.reaction.Receiver <= 0 {
				require.Error(t, storage.AddReaction(ctx, tt.reaction))
				return
			}

			if tt.queryErr != nil {
				mock.ExpectExec("INSERT INTO reaction").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectExec("INSERT INTO reaction").
					WithArgs(tt.reaction.Author, tt.reaction.Receiver, tt.reaction.Type).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := storage.AddReaction(ctx, tt.reaction)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr), "expected error to be %v but got %v", tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetReactionList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		userId   int
		mockRows *sqlmock.Rows
		queryErr error
		wantErr  error
		wantList []int
	}{
		{
			name:     "successful reaction list retrieval",
			userId:   1,
			mockRows: sqlmock.NewRows([]string{"receiver"}).AddRow(2).AddRow(3),
			queryErr: nil,
			wantErr:  nil,
			wantList: []int{2, 3},
		},

		{
			name:     "database error on reaction list retrieval",
			userId:   1,
			mockRows: nil,
			queryErr: errors.New("some select error"),
			wantErr:  fmt.Errorf("failed to select: %v", errors.New("some select error")),
			wantList: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT receiver FROM reaction WHERE author =").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT receiver FROM reaction WHERE author =").
					WithArgs(tt.userId).
					WillReturnRows(tt.mockRows)
			}

			gotList, err := storage.GetReactionList(ctx, tt.userId)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Nil(t, gotList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantList, gotList)
			}
		})
	}
}

func TestGetMatchList(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name     string
		userId   int
		mockRows *sqlmock.Rows
		queryErr error
		wantErr  error
		wantList []int
	}{
		{
			name:     "Successful match list retrieval",
			userId:   1,
			mockRows: sqlmock.NewRows([]string{"author"}).AddRow(2).AddRow(3),
			queryErr: nil,
			wantErr:  nil,
			wantList: []int{2, 3},
		},
		{
			name:     "Database error on match list retrieval",
			userId:   1,
			mockRows: nil,
			queryErr: errors.New("query error"),
			wantErr:  fmt.Errorf("failed to select: %v", errors.New("query error")),
			wantList: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT author FROM reaction").WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT author FROM reaction").
					WithArgs(tt.userId, tt.userId).
					WillReturnRows(tt.mockRows)
			}

			gotList, err := storage.GetMatchList(ctx, tt.userId)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Nil(t, gotList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantList, gotList)
			}
		})
	}
}

func TestGetMatchTime(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name          string
		firstUser     int
		secondUser    int
		mockQueryTime string
		queryErr      error
		wantTime      string
		wantErr       error
	}{
		{
			name:          "Successful match time retrieval",
			firstUser:     1,
			secondUser:    2,
			mockQueryTime: "2024-11-27T10:00:00Z",
			queryErr:      nil,
			wantTime:      "2024-11-27T10:00:00Z",
			wantErr:       nil,
		},
		{
			name:          "Database error on match time retrieval",
			firstUser:     1,
			secondUser:    2,
			mockQueryTime: "",
			queryErr:      errors.New("query error"),
			wantTime:      "",
			wantErr:       fmt.Errorf("failed to get match time: %w", errors.New("query error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.queryErr != nil {
				mock.ExpectQuery("SELECT created_at FROM reaction").
					WithArgs(tt.firstUser, tt.secondUser).
					WillReturnError(tt.queryErr)
			} else {
				mock.ExpectQuery("SELECT created_at FROM reaction").
					WithArgs(tt.firstUser, tt.secondUser).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(tt.mockQueryTime))
			}

			gotTime, err := storage.GetMatchTime(ctx, tt.firstUser, tt.secondUser)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Empty(t, gotTime)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantTime, gotTime)
			}
		})
	}
}

func TestUpdateOrCreateReaction(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name           string
		reaction       models.Reaction
		updateAffected int64
		updateErr      error
		insertErr      error
		wantErr        error
	}{
		{
			name:           "Successful update",
			reaction:       models.Reaction{Author: 1, Receiver: 2, Type: true},
			updateAffected: 1,
			updateErr:      nil,
			insertErr:      nil,
			wantErr:        nil,
		},
		{
			name:           "Update no rows, successful insert",
			reaction:       models.Reaction{Author: 1, Receiver: 2, Type: true},
			updateAffected: 0,
			updateErr:      nil,
			insertErr:      nil,
			wantErr:        nil,
		},
		{
			name:           "Update error",
			reaction:       models.Reaction{Author: 1, Receiver: 2, Type: true},
			updateAffected: 0,
			updateErr:      errors.New("update error"),
			insertErr:      nil,
			wantErr:        fmt.Errorf("failed to update reaction: %w", errors.New("update error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{DB: db, logger: logger}

			if tt.updateErr != nil {
				mock.ExpectExec("UPDATE reaction SET type =").WillReturnError(tt.updateErr)
			} else {
				mock.ExpectExec("UPDATE reaction SET type =").
					WithArgs(tt.reaction.Type, tt.reaction.Author, tt.reaction.Receiver).
					WillReturnResult(sqlmock.NewResult(0, tt.updateAffected))
			}

			if tt.updateAffected == 0 {
				if tt.insertErr != nil {
					mock.ExpectExec("INSERT INTO reaction").WillReturnError(tt.insertErr)
				} else {
					mock.ExpectExec("INSERT INTO reaction").
						WithArgs(tt.reaction.Author, tt.reaction.Receiver, tt.reaction.Type).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
			}

			err := storage.UpdateOrCreateReaction(ctx, tt.reaction)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
