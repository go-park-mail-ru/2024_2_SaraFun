package repo

import (
	"context"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("INSERT INTO balance").WithArgs(tt.userID, tt.amount).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("INSERT INTO balance").WithArgs(tt.userID, tt.amount).WillReturnResult(tt.execResult)
			}
			err := repo.AddBalance(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestAddDailyLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("INSERT INTO daily_likes").WithArgs(tt.userID, tt.amount).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("INSERT INTO daily_likes").WithArgs(tt.userID, tt.amount).WillReturnResult(tt.execResult)
			}
			err := repo.AddDailyLikeCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestAddPurchasedLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("INSERT INTO purchased_likes").WithArgs(tt.userID, tt.amount).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("INSERT INTO purchased_likes").WithArgs(tt.userID, tt.amount).WillReturnResult(tt.execResult)
			}
			err := repo.AddPurchasedLikeCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestChangeBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE balance").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE balance").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.ChangeBalance(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestChangeDailyLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.ChangeDailyLikeCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestChangePurchasedLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE purchased_likes").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE purchased_likes").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.ChangePurchasedLikeCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestSetBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE balance").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE balance").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.SetBalance(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestSetDailyLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.SetDailyLikesCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestSetDailyLikesCountToAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE daily_likes").WithArgs(tt.amount).WillReturnResult(tt.execResult)
			}
			err := repo.SetDailyLikesCountToAll(ctx, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestSetPurchasedLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		userID     int
		amount     int
		execResult driver.Result
		execError  error
	}{
		{
			name:       "good test",
			userID:     1,
			amount:     10,
			execError:  nil,
			execResult: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectExec("UPDATE purchased_likes").WithArgs(tt.amount, tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectExec("UPDATE purchased_likes").WithArgs(tt.amount, tt.userID).WillReturnResult(tt.execResult)
			}
			err := repo.SetPurchasedLikesCount(ctx, tt.userID, tt.amount)
			require.Equal(t, tt.execError, err)

		})
	}
}

func TestGetPurchasedLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		name           string
		userID         int
		expectedAmount int
		execResult     *sqlmock.Rows
		execError      error
	}{
		{
			name:           "good test",
			userID:         1,
			expectedAmount: 10,
			execError:      nil,
			execResult:     sqlmock.NewRows([]string{"amount"}).AddRow(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectQuery("SELECT likes_count FROM purchased_likes").WithArgs(tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectQuery("SELECT likes_count FROM purchased_likes").WithArgs(tt.userID).WillReturnRows(tt.execResult)
			}
			amount, err := repo.GetPurchasedLikesCount(ctx, tt.userID)
			require.Equal(t, tt.execError, err)
			require.Equal(t, amount, tt.expectedAmount)

		})
	}
}

func TestGetDailyLikesCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		name           string
		userID         int
		expectedAmount int
		execResult     *sqlmock.Rows
		execError      error
	}{
		{
			name:           "good test",
			userID:         1,
			expectedAmount: 10,
			execError:      nil,
			execResult:     sqlmock.NewRows([]string{"amount"}).AddRow(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectQuery("SELECT likes_count FROM daily_likes").WithArgs(tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectQuery("SELECT likes_count FROM daily_likes").WithArgs(tt.userID).WillReturnRows(tt.execResult)
			}
			amount, err := repo.GetDailyLikesCount(ctx, tt.userID)
			require.Equal(t, tt.execError, err)
			require.Equal(t, amount, tt.expectedAmount)

		})
	}
}

func TestGetBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		name           string
		userID         int
		expectedAmount int
		execResult     *sqlmock.Rows
		execError      error
	}{
		{
			name:           "good test",
			userID:         1,
			expectedAmount: 10,
			execError:      nil,
			execResult:     sqlmock.NewRows([]string{"amount"}).AddRow(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.execError != nil {
				mock.ExpectQuery("SELECT balance FROM balance").WithArgs(tt.userID).WillReturnError(tt.execError)
			} else {
				mock.ExpectQuery("SELECT balance FROM balance").WithArgs(tt.userID).WillReturnRows(tt.execResult)
			}
			amount, err := repo.GetBalance(ctx, tt.userID)
			require.Equal(t, tt.execError, err)
			require.Equal(t, amount, tt.expectedAmount)

		})
	}
}
