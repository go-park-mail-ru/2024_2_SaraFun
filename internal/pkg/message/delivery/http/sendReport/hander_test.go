package sendReport

import (
	"bytes"
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	communicationsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen/mocks"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	messagemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen/mocks"
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
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	messageClient := messagemocks.NewMockMessageClient(mockCtrl)
	communicationsClient := communicationsmocks.NewMockCommunicationsClient(mockCtrl)
	handler := NewHandler(authClient, messageClient, communicationsClient, logger)

	tests := []struct {
		name                string
		method              string
		path                string
		body                []byte
		cookieValue         string
		authReturn          int
		authError           error
		authTimes           int
		receiver            int
		reportBody          string
		reportReason        string
		addReportReturn     int
		updateOrCreateError error
		updateOrCreateTimes int
		expectedStatus      int
		expectedMessage     string
	}{
		{
			name:   "good test",
			method: "POST",
			path:   "/api/message/report",
			body: []byte(`{
    "receiver": 2,
    "reason": "abuse",
    "body": "он нереальный мудак"
}`),
			reportBody:          "он нереальный мудак",
			reportReason:        "abuse",
			receiver:            2,
			cookieValue:         "sparkit",
			authReturn:          1,
			authTimes:           1,
			authError:           nil,
			addReportReturn:     1,
			updateOrCreateError: nil,
			updateOrCreateTimes: 1,
			expectedStatus:      http.StatusOK,
			expectedMessage:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserIDReq).Return(getUserIDResponse, tt.authError).Times(tt.authTimes)

			report := &generatedMessage.Report{
				Author:   int32(tt.authReturn),
				Receiver: int32(tt.receiver),
				Body:     tt.reportBody,
				Reason:   tt.reportReason,
			}
			addReportReq := &generatedMessage.AddReportRequest{Report: report}
			addReportResponse := &generatedMessage.AddReportResponse{ReportID: int32(tt.addReportReturn)}
			messageClient.EXPECT().AddReport(ctx, addReportReq).Return(addReportResponse, tt.authError).Times(tt.authTimes)

			reaction := &generatedCommunications.Reaction{
				Author:   int32(tt.authReturn),
				Receiver: int32(tt.receiver),
				Type:     false,
			}
			updateOrCreateReq := &generatedCommunications.UpdateOrCreateReactionRequest{Reaction: reaction}
			updateOrCreateResponse := &generatedCommunications.UpdateOrCreateReactionResponse{}
			communicationsClient.EXPECT().UpdateOrCreateReaction(ctx, updateOrCreateReq).Return(updateOrCreateResponse, tt.updateOrCreateError).
				Times(tt.updateOrCreateTimes)

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
