package sendmessage

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	communicationsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen/mocks"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	messagemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen/mocks"
	websocketmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/http/sendmessage/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	messageClient := messagemocks.NewMockMessageClient(mockCtrl)
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	communicationsClient := communicationsmocks.NewMockCommunicationsClient(mockCtrl)
	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
	wsClient := websocketmocks.NewMockWebSocketService(mockCtrl)
	handler := NewHandler(messageClient, wsClient, authClient, communicationsClient, personalitiesClient, logger)

	tests := []struct {
		name            string
		method          string
		path            string
		body            []byte
		message         models.Message
		cookieValue     string
		username        string
		status          string
		secondUserID    int
		authReturn      int
		authError       error
		authTimes       int
		usernameReturn  string
		usernameError   error
		usernameTimes   int
		messageIDReturn int
		messageError    error
		messageTimes    int
		wsError         error
		wsTimes         int
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:   "good test",
			method: "POST",
			path:   "/api/message/message",
			body: []byte(`{
    "receiver": 2,
    "body": "тестовое сообщение"
}`),
			cookieValue:  "sparkit",
			authReturn:   1,
			authError:    nil,
			authTimes:    1,
			secondUserID: 2,
			message: models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "тестовое сообщение",
				Time:     time.DateTime,
			},
			status:          "",
			usernameReturn:  "sparkit",
			usernameError:   nil,
			usernameTimes:   1,
			messageIDReturn: 1,
			messageError:    nil,
			messageTimes:    1,
			wsError:         nil,
			wsTimes:         1,
			expectedStatus:  http.StatusOK,
		},
		{

			name:   "check block exists",
			method: "POST",
			path:   "/api/message/message",
			body: []byte(`{
    "receiver": 2,
    "body": "тестовое сообщение"
}`),
			cookieValue:  "sparkit",
			authReturn:   1,
			authError:    nil,
			authTimes:    1,
			secondUserID: 2,
			message: models.Message{
				Author:   1,
				Receiver: 2,
				Body:     "тестовое сообщение",
				Time:     time.DateTime,
			},
			status:          "block exists",
			usernameReturn:  "sparkit",
			usernameError:   nil,
			usernameTimes:   1,
			messageIDReturn: 1,
			messageError:    nil,
			messageTimes:    1,
			wsError:         nil,
			wsTimes:         1,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "{\"info\":\"block exists\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getuserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserIDReq).Return(getuserIDResponse, tt.authError).Times(tt.authTimes)

			checkBlockReq := &generatedMessage.CheckUsersBlockNotExistsRequest{
				FirstUserID:  int32(tt.authReturn),
				SecondUserID: int32(tt.secondUserID),
			}
			checkBlockResponse := &generatedMessage.CheckUsersBlockNotExistsResponse{Status: tt.status}
			messageClient.EXPECT().CheckUsersBlockNotExists(ctx, checkBlockReq).Return(checkBlockResponse, tt.authError).Times(tt.authTimes)

			chatMessage := &generatedMessage.ChatMessage{
				ID:       0,
				Author:   int32(tt.message.Author),
				Receiver: int32(tt.message.Receiver),
				Body:     tt.message.Body,
			}
			addMessageReq := &generatedMessage.AddMessageRequest{Message: chatMessage}
			addMessageResponse := &generatedMessage.AddMessageResponse{MessageID: int32(tt.messageIDReturn)}
			messageClient.EXPECT().AddMessage(ctx, addMessageReq).Return(addMessageResponse, tt.messageError).Times(tt.messageTimes)

			getUsernameReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: int32(tt.authReturn)}
			getUsernameResponse := &generatedPersonalities.GetUsernameByUserIDResponse{Username: tt.usernameReturn}
			personalitiesClient.EXPECT().GetUsernameByUserID(ctx, getUsernameReq).Return(getUsernameResponse, tt.usernameError).Times(tt.usernameTimes)

			wsClient.EXPECT().WriteMessage(ctx, tt.message.Author, tt.message.Receiver, tt.message.Body, tt.usernameReturn).Return(tt.wsError).Times(tt.wsTimes)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: tt.cookieValue,
			}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}
		})
	}
}
