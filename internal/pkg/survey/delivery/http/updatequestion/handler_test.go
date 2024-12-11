package updatequestion

import (
	"bytes"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	surveyClient := surveymocks.NewMockSurveyClient(mockCtrl)
	handler := NewHandler(authClient, surveyClient, logger)

	tests := []struct {
		name            string
		method          string
		path            string
		body            []byte
		content         string
		newContent      string
		grade           int
		cookieValue     string
		authReturn      int
		authError       error
		authTimes       int
		surveyReturn    int
		surveyError     error
		surveyTimes     int
		expectedCode    int
		expectedMessage string
	}{
		{
			name:   "good test",
			method: "PUT",
			path:   "/api/survey/question",
			body: []byte(`{
    "old_content": "",
    "new_content": ""}`),
			cookieValue: "sparkit",
			content:     "",
			newContent:  "",
			//grade:           1,
			authReturn:      1,
			authError:       nil,
			authTimes:       1,
			surveyReturn:    1,
			surveyError:     nil,
			surveyTimes:     1,
			expectedCode:    http.StatusOK,
			expectedMessage: "{\"id\":1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserIDRequest).Return(getUserIDResponse, tt.authError).
				Times(tt.authTimes)
			updateQuestion := &generatedSurvey.AdminQuestion{
				Content: tt.newContent,
				Grade:   int32(tt.grade),
			}
			t.Log(updateQuestion)
			updateQuestionRequest := &generatedSurvey.UpdateQuestionRequest{
				Question: updateQuestion,
				Content:  tt.content,
			}
			updateQuestionResponse := &generatedSurvey.UpdateQuestionResponse{Id: int32(tt.surveyReturn)}
			t.Log(updateQuestionRequest)
			surveyClient.EXPECT().UpdateQuestion(gomock.Any(), gomock.Any()).Return(updateQuestionResponse, tt.surveyError).
				Times(tt.surveyTimes)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: tt.cookieValue,
			}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			handler.Handle(w, req)
			if w.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedCode)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}
		})
	}
}
