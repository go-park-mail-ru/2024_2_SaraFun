package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func TestUseCase_AddConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()
	userID := 10

	var conn *websocket.Conn = nil

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().AddConnection(ctx, conn, userID).Return(nil)
		err := uc.AddConnection(ctx, conn, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().AddConnection(ctx, conn, userID).Return(errors.New("db error"))
		err := uc.AddConnection(ctx, conn, userID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "repo Add connection call in Usecase failed: db error" {
			t.Errorf("error message mismatch: got %v, want %v", err.Error(), "repo Add connection call in Usecase failed: db error")
		}
	})
}

func TestUseCase_DeleteConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()
	userID := 20

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().DeleteConnection(ctx, userID).Return(nil)
		err := uc.DeleteConnection(ctx, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().DeleteConnection(ctx, userID).Return(errors.New("db error"))
		err := uc.DeleteConnection(ctx, userID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "repo Delete connection call in Usecase failed: db error" {
			t.Errorf("error message mismatch: got %v, want %v", err.Error(), "repo Delete connection call in Usecase failed: db error")
		}
	})
}

func TestUseCase_WriteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()
	authorID := 30
	receiverID := 40
	message := "Hello"
	username := "UserA"

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().WriteMessage(ctx, authorID, receiverID, message, username).Return(nil)
		err := uc.WriteMessage(ctx, authorID, receiverID, message, username)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().WriteMessage(ctx, authorID, receiverID, message, username).Return(errors.New("db error"))
		err := uc.WriteMessage(ctx, authorID, receiverID, message, username)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "repo WriteMessage call in Usecase failed: db error" {
			t.Errorf("error message mismatch: got %v, want %v", err.Error(), "repo WriteMessage call in Usecase failed: db error")
		}
	})
}

func TestUseCase_SendNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()
	uc := New(repo, logger)

	ctx := context.Background()
	receiverID := 50
	authorUsername := "AuthorUser"
	authorImageLink := "http://example.com/image.jpg"

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().SendNotification(ctx, receiverID, authorImageLink, authorUsername).Return(nil)
		err := uc.SendNotification(ctx, receiverID, authorUsername, authorImageLink)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().SendNotification(ctx, receiverID, authorImageLink, authorUsername).Return(errors.New("db error"))
		err := uc.SendNotification(ctx, receiverID, authorUsername, authorImageLink)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "repo SendNotification call in Usecase failed: db error" {
			t.Errorf("error message mismatch: got %v, want %v", err.Error(), "repo SendNotification call in Usecase failed: db error")
		}
	})
}
