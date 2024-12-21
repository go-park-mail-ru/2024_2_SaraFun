package setConnection

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/delivery/setConnection/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func TestHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase := mocks.NewMockUseCase(ctrl)
	authClient := mocks.NewMockAuthClient(ctrl)
	logger := zap.NewNop()

	h := NewHandler(useCase, authClient, logger)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")

	t.Run("no cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/setConnection", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Handle(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status code mismatch: got %v, want %v", w.Code, http.StatusUnauthorized)
		}
		if w.Body.String() != "cookie not found\n" {
			t.Errorf("body mismatch: got %v, want %v", w.Body.String(), "cookie not found\n")
		}
	})

	t.Run("auth error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/setConnection", nil).WithContext(ctx)
		req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: "valid_session"})

		authClient.EXPECT().
			GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{SessionID: "valid_session"}).
			Return(nil, errors.New("auth error"))

		w := httptest.NewRecorder()

		h.Handle(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status code mismatch: got %v, want %v", w.Code, http.StatusUnauthorized)
		}
		if w.Body.String() != "get user id error\n" {
			t.Errorf("body mismatch: got %v, want %v", w.Body.String(), "get user id error\n")
		}
	})

	t.Run("success scenario", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx)
			r.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: "valid_session"})
			authClient.EXPECT().
				GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{SessionID: "valid_session"}).
				Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 999}, nil)

			useCase.EXPECT().AddConnection(gomock.Any(), gomock.Any(), 999).Return(nil)
			useCase.EXPECT().DeleteConnection(gomock.Any(), 999).Return(nil).AnyTimes() // вызывается при defer после закрытия соединения

			h.Handle(w, r)
		}))
		defer server.Close()

		u := "ws" + server.URL[len("http"):]
		dialer := websocket.Dialer{}
		conn, resp, err := dialer.Dial(u, nil)
		defer resp.Body.Close()
		if err != nil {
			t.Fatalf("failed to dial websocket: %v", err)
		}
		defer conn.Close()

		// Отправим CloseMessage чтобы закрыть соединение и вызвать deleteConnection
		err = conn.WriteMessage(websocket.CloseMessage, []byte("close"))
		if err != nil {
			t.Errorf("failed to write close message: %v", err)
		}
	})
}
