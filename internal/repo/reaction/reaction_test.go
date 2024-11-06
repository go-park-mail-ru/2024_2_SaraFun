package reaction

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
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
