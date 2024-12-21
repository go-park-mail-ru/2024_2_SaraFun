package authgrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGrpcAuthHandler_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := authgrpc.NewGRPCAuthHandler(uc, logger)

	ctx := context.Background()
	req := &generatedAuth.CreateSessionRequest{
		User: &generatedAuth.User{
			ID: 123,
		},
	}

	session := models.Session{
		SessionID: "session123",
		UserID:    123,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().CreateSession(ctx, models.User{ID: 123}).Return(session, nil)
		resp, err := h.CreateSession(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Session.SessionID != "session123" || resp.Session.UserID != 123 {
			t.Errorf("response mismatch: got %+v", resp.Session)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().CreateSession(ctx, models.User{ID: 123}).Return(models.Session{}, errors.New("create error"))
		_, err := h.CreateSession(ctx, req)
		if err == nil || !contains(err.Error(), "create error") {
			t.Errorf("expected error, got %v", err)
		}
	})
}

func TestGrpcAuthHandler_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := authgrpc.NewGRPCAuthHandler(uc, logger)

	ctx := context.Background()
	req := &generatedAuth.DeleteSessionRequest{
		SessionID: "sessionXYZ",
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().DeleteSession(ctx, "sessionXYZ").Return(nil)
		_, err := h.DeleteSession(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().DeleteSession(ctx, "sessionXYZ").Return(errors.New("delete error"))
		_, err := h.DeleteSession(ctx, req)
		if err == nil || !contains(err.Error(), "Grpc delete session error") {
			t.Errorf("expected error, got %v", err)
		}
	})
}

func TestGrpcAuthHandler_CheckSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := authgrpc.NewGRPCAuthHandler(uc, logger)

	ctx := context.Background()
	req := &generatedAuth.CheckSessionRequest{
		SessionID: "checkSession",
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().CheckSession(ctx, "checkSession").Return(nil)
		_, err := h.CheckSession(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().CheckSession(ctx, "checkSession").Return(errors.New("check error"))
		_, err := h.CheckSession(ctx, req)
		if err == nil || !contains(err.Error(), "Grpc check session error") {
			t.Errorf("expected error, got %v", err)
		}
	})
}

func TestGrpcAuthHandler_GetUserIDBySessionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := authgrpc.NewGRPCAuthHandler(uc, logger)

	ctx := context.Background()
	req := &generatedAuth.GetUserIDBySessionIDRequest{
		SessionID: "getUser",
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().GetUserIDBySessionID(ctx, "getUser").Return(999, nil)
		resp, err := h.GetUserIDBySessionID(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.UserId != 999 {
			t.Errorf("got %v, want 999", resp.UserId)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetUserIDBySessionID(ctx, "getUser").Return(0, errors.New("user error"))
		_, err := h.GetUserIDBySessionID(ctx, req)
		if err == nil || !contains(err.Error(), "Grpc get user id by session id error") {
			t.Errorf("expected error, got %v", err)
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
