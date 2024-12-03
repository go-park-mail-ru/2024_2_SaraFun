package addquestion

import (
	"bytes"
	"context"
	"errors"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	surveymocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	surveyClient := surveymocks.NewMockSurveyClient(mockCtrl)
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	handler := NewHandler(authClient, surveyClient, logger)
	tests := []struct {
		name                   string
		method                 string
		path                   string
		body                   []byte
		content                string
		grade                  int
		authReturnUserID       int
		authReturnError        error
		authReturnCount        int
		surveyReturnQuestionID int
		surveyReturnError      error
		surveyReturnCount      int
		cookieValue            string
		expectedStatus         int
		expectedMessage        string
	}{
		{
			name:   "good test",
			method: http.MethodPost,
			path:   "/survey/question",
			body: []byte(
				`{
"content": "тест",
"grade": 5
}`),
			content:                "тест",
			grade:                  5,
			authReturnUserID:       1,
			authReturnError:        nil,
			authReturnCount:        1,
			surveyReturnQuestionID: 1,
			surveyReturnError:      nil,
			surveyReturnCount:      1,
			cookieValue:            "sparkit",
			expectedStatus:         http.StatusOK,
			expectedMessage:        `{"id":1}`,
		},
		{
			name:   "bad request",
			method: http.MethodPost,
			path:   "/survey/question",
			body: []byte(
				``),
			authReturnUserID: 1,
			authReturnError:  nil,
			authReturnCount:  1,
			cookieValue:      "sparkit",
			expectedStatus:   http.StatusBadRequest,
			expectedMessage:  "json decode question error\n",
		},
		{
			name:   "bad cookie",
			method: http.MethodPost,
			path:   "/survey/question",
			body: []byte(
				`{
"content": "тест",
"grade": 5
}`),
			authReturnUserID: -1,
			authReturnError:  errors.New("bad cookie"),
			authReturnCount:  1,
			cookieValue:      "badcookie",
			expectedStatus:   http.StatusUnauthorized,
			expectedMessage:  "get user id by session id\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturnUserID)}
			reqQuestion := &generatedSurvey.AdminQuestion{
				Content: tt.content,
				Grade:   int32(tt.grade),
			}
			addQuestionReq := &generatedSurvey.AddQuestionRequest{
				Question: reqQuestion,
			}
			addQuestionResponse := &generatedSurvey.AddQuestionResponse{QuestionID: int32(tt.surveyReturnQuestionID)}

			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserIDReq).Return(getUserIDResponse, tt.authReturnError).
				Times(tt.authReturnCount)
			surveyClient.EXPECT().AddQuestion(ctx, addQuestionReq).Return(addQuestionResponse, tt.surveyReturnError).
				Times(tt.surveyReturnCount)

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
