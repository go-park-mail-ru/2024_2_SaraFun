package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-redis/redismock/v9"
	"go.uber.org/zap"
)

func TestInMemoryStorage_AddSession(t *testing.T) {
	logger := zap.NewNop()
	db, mock := redismock.NewClientMock()
	repo := New(db, logger)

	ctx := context.Background()
	session := models.Session{
		SessionID: "session123",
		UserID:    10,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	t.Run("error", func(t *testing.T) {
		mock.ExpectSet("session123", "10", 24*time.Hour).SetErr(errors.New("set error"))

		err := repo.AddSession(ctx, session)
		if err == nil || !contains(err.Error(), "add session failed") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestInMemoryStorage_GetUserIDBySessionID(t *testing.T) {
	logger := zap.NewNop()
	db, mock := redismock.NewClientMock()
	repo := New(db, logger)

	ctx := context.Background()
	sessionID := "sessionXYZ"

	t.Run("success", func(t *testing.T) {
		mock.ExpectGet(sessionID).SetVal("123")

		userID, err := repo.GetUserIDBySessionID(ctx, sessionID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if userID != 123 {
			t.Errorf("got %d, want 123", userID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectGet(sessionID).RedisNil()

		_, err := repo.GetUserIDBySessionID(ctx, sessionID)
		if err == nil || !errors.Is(err, sparkiterrors.ErrInvalidSession) {
			t.Errorf("expected ErrInvalidSession got %v", err)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		mock.ExpectGet(sessionID).SetVal("not-a-number")

		_, err := repo.GetUserIDBySessionID(ctx, sessionID)
		if err == nil || !contains(err.Error(), "convert session id") {
			t.Errorf("expected convert error got %v", err)
		}
	})
}

func TestInMemoryStorage_CheckSession(t *testing.T) {
	logger := zap.NewNop()
	db, mock := redismock.NewClientMock()
	repo := New(db, logger)

	ctx := context.Background()
	sessionID := "checkSession"

	t.Run("valid session", func(t *testing.T) {
		mock.ExpectGet(sessionID).SetVal("456")

		err := repo.CheckSession(ctx, sessionID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectGet(sessionID).RedisNil()

		err := repo.CheckSession(ctx, sessionID)
		if err == nil || !errors.Is(err, sparkiterrors.ErrInvalidSession) {
			t.Errorf("expected ErrInvalidSession got %v", err)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		mock.ExpectGet(sessionID).SetVal("abc")

		err := repo.CheckSession(ctx, sessionID)
		if err == nil || !contains(err.Error(), "convert session id") {
			t.Errorf("expected convert error got %v", err)
		}
	})
}

func TestInMemoryStorage_DeleteSession(t *testing.T) {
	logger := zap.NewNop()
	db, mock := redismock.NewClientMock()
	repo := New(db, logger)

	ctx := context.Background()
	sessionID := "deleteSession"

	t.Run("success", func(t *testing.T) {
		mock.ExpectDel(sessionID).SetVal(1)

		err := repo.DeleteSession(ctx, sessionID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectDel(sessionID).SetErr(errors.New("del error"))

		err := repo.DeleteSession(ctx, sessionID)
		if err == nil || !contains(err.Error(), "delete session failed") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
