package repo

import (
	"context"
	"fmt"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

type Storage struct {
	wConns map[int]*ws.Conn
	mu     sync.RWMutex
	logger *zap.Logger
}

//type JsonMessage struct {
//	AuthorID int `json:"author_id"`
//	//ReceiverID int    `json:"receiver_id"`
//	Message string `json:"message"`
//}
//
//type JsonNotification struct {
//	Username  string `json:"username"`
//	Imagelink string `json:"imagelink"`
//	Type string `json:"type"`
//}

type JsonWS struct {
	AuthorID  int    `json:"author_id"`
	Message   string `json:"message"`
	Username  string `json:"username"`
	Imagelink string `json:"imagelink"`
	Type      string `json:"type"`
}

type JsonNotification struct {
	Username  string `json:"username"`
	Imagelink string `json:"imagelink"`
}

func New(conns map[int]*ws.Conn, logger *zap.Logger) *Storage {
	return &Storage{
		wConns: conns,
		mu:     sync.RWMutex{},
		logger: logger,
	}
}

func (s *Storage) AddConnection(ctx context.Context, conn *ws.Conn, userId int) error {
	s.logger.Info("Repo websocket addConnection", zap.Int("userId", userId))
	s.mu.Lock()
	s.wConns[userId] = conn
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteConnection(ctx context.Context, userId int) error {
	s.logger.Info("Repo websocket deleteConnection", zap.Int("userId", userId))
	s.mu.Lock()
	delete(s.wConns, userId)
	s.mu.Unlock()
	return nil
}

func (s *Storage) WriteMessage(ctx context.Context, authorID int, receiverID int,
	message string, username string) error {
	s.logger.Info("Repo websocket writeMessage", zap.Int("receiverID", receiverID))
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.wConns[receiverID]
	if !ok {
		return fmt.Errorf("user ws conn not found: %v", receiverID)
	}
	msg := JsonWS{
		AuthorID: authorID,
		Message:  message,
		Username: username,
		Type:     "message",
	}
	err := conn.WriteJSON(&msg)
	if err != nil {
		return fmt.Errorf("cannot write message: %w", err)
	}
	return nil
}

func (s *Storage) SendNotification(ctx context.Context, receiverID int, authorImageLink string, authorUsername string) error {
	s.logger.Info("Repo websocket sendNotification", zap.Int("receiverID", receiverID))
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.wConns[receiverID]
	if !ok {
		return fmt.Errorf("user ws conn not found: %v", receiverID)
	}
	notification := JsonWS{
		Username:  authorUsername,
		Imagelink: authorImageLink,
		Type:      "notification",
	}
	err := conn.WriteJSON(&notification)
	if err != nil {
		return fmt.Errorf("cannot send notification: %w", err)
	}
	return nil
}
