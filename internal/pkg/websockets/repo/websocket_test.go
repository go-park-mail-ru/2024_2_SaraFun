package repo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func TestStorage(t *testing.T) {
	logger := zap.NewNop()

	s := New(make(map[int]*websocket.Conn), logger)

	ctx := context.Background()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var serverConn *websocket.Conn
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		serverConn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade error", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	clientConn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	defer clientConn.Close()

	userID := 123

	err = s.AddConnection(ctx, serverConn, userID)
	if err != nil {
		t.Errorf("AddConnection error: %v", err)
	}

	s.mu.RLock()
	if s.wConns[userID] == nil {
		t.Errorf("connection not added to map")
	}
	s.mu.RUnlock()

	message := "Hello"
	authorID := 10
	username := "testuser"

	err = s.WriteMessage(ctx, authorID, userID, message, username)
	if err != nil {
		t.Errorf("WriteMessage error: %v", err)
	}

	setErr := clientConn.SetReadDeadline(time.Now().Add(time.Second))
	if setErr != nil {
		t.Errorf("SetReadDeadline error")
	}
	mt, msgData, err := clientConn.ReadMessage()
	if err != nil {
		t.Errorf("failed to read message: %v", err)
	}
	if mt != websocket.TextMessage {
		t.Errorf("expected text message, got %v", mt)
	}

	msgStr := string(msgData)
	if !contains(msgStr, `"type":"message"`) ||
		!contains(msgStr, fmt.Sprintf(`"message":"%s"`, message)) ||
		!contains(msgStr, fmt.Sprintf(`"author_id":%d`, authorID)) ||
		!contains(msgStr, fmt.Sprintf(`"username":"%s"`, username)) {
		t.Errorf("message format mismatch: %v", msgStr)
	}

	authorImageLink := "http://example.com/image.png"
	err = s.SendNotification(ctx, userID, authorImageLink, username)
	if err != nil {
		t.Errorf("SendNotification error: %v", err)
	}

	setErr = clientConn.SetReadDeadline(time.Now().Add(time.Second))
	if setErr != nil {
		t.Errorf("SetReadDeadline error")
	}
	mt, notifData, err := clientConn.ReadMessage()
	if err != nil {
		t.Errorf("failed to read notification: %v", err)
	}
	if mt != websocket.TextMessage {
		t.Errorf("expected text message for notification, got %v", mt)
	}

	notifStr := string(notifData)
	if !contains(notifStr, `"type":"notification"`) ||
		!contains(notifStr, fmt.Sprintf(`"imagelink":"%s"`, authorImageLink)) ||
		!contains(notifStr, fmt.Sprintf(`"username":"%s"`, username)) {
		t.Errorf("notification format mismatch: %v", notifStr)
	}

	err = s.WriteMessage(ctx, authorID, 999, "no user", "no user")
	if err == nil {
		t.Errorf("expected error for non-existent user, got nil")
	}
	if !contains(err.Error(), "user ws conn not found") {
		t.Errorf("error message mismatch: got %v", err.Error())
	}

	err = s.SendNotification(ctx, 999, "img", "user")
	if err == nil {
		t.Errorf("expected error for non-existent user in notification, got nil")
	}
	if !contains(err.Error(), "user ws conn not found") {
		t.Errorf("error message mismatch: got %v", err.Error())
	}

	err = s.DeleteConnection(ctx, userID)
	if err != nil {
		t.Errorf("DeleteConnection error: %v", err)
	}

	s.mu.RLock()
	if _, ok := s.wConns[userID]; ok {
		t.Errorf("connection not deleted from map")
	}
	s.mu.RUnlock()
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
