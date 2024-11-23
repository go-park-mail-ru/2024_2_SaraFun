package likehandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	sparkiterrors "sparkit/internal/errors"
	likehandler_mocks "sparkit/internal/handlers/likehandler/mocks"
	"sparkit/internal/utils/consts"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLikeHandler_Handle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                     string
		method                   string
		path                     string
		body                     []byte
		sessionCookie            *http.Cookie
		getUserIDError           error
		processLikeError         error
		expectedStatus           int
		expectedMessage          string
		getUserIDBySessionCalled bool
		processLikeCalled        bool
	}{
		{
			name:                     "успешный лайк",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 2, "is_like": true}`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "valid_session"},
			getUserIDError:           nil,
			processLikeError:         nil,
			expectedStatus:           http.StatusOK,
			expectedMessage:          "ok",
			getUserIDBySessionCalled: true,
			processLikeCalled:        true,
		},
		{
			name:                     "отсутствует сессионный куки",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 2, "is_like": true}`),
			sessionCookie:            nil,
			expectedStatus:           http.StatusUnauthorized,
			expectedMessage:          "Unauthorized\n",
			getUserIDBySessionCalled: false,
			processLikeCalled:        false,
		},
		{
			name:                     "недействительная сессия",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 2, "is_like": true}`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "invalid_session"},
			getUserIDError:           errors.New("invalid session"),
			expectedStatus:           http.StatusUnauthorized,
			expectedMessage:          "Unauthorized\n",
			getUserIDBySessionCalled: true,
			processLikeCalled:        false,
		},
		{
			name:                     "неправильный метод запроса",
			method:                   "GET",
			path:                     "/like",
			expectedStatus:           http.StatusMethodNotAllowed,
			expectedMessage:          "Method not allowed\n",
			getUserIDBySessionCalled: false,
			processLikeCalled:        false,
		},
		{
			name:                     "недопустимое тело запроса",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`invalid_json`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "valid_session"},
			expectedStatus:           http.StatusBadRequest,
			expectedMessage:          "Invalid request body\n",
			getUserIDBySessionCalled: true,
			processLikeCalled:        false,
		},
		{
			name:                     "пользователь не найден",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 999, "is_like": true}`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "valid_session"},
			getUserIDError:           nil,
			processLikeError:         sparkiterrors.ErrUserNotFound,
			expectedStatus:           http.StatusNotFound,
			expectedMessage:          "Target user not found\n",
			getUserIDBySessionCalled: true,
			processLikeCalled:        true,
		},
		{
			name:                     "нельзя лайкнуть себя",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 1, "is_like": true}`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "valid_session"},
			getUserIDError:           nil,
			processLikeError:         sparkiterrors.ErrCannotLikeSelf,
			expectedStatus:           http.StatusBadRequest,
			expectedMessage:          "Cannot like/dislike self\n",
			getUserIDBySessionCalled: true,
			processLikeCalled:        true,
		},
		{
			name:                     "внутренняя ошибка сервера",
			method:                   "POST",
			path:                     "/like",
			body:                     []byte(`{"target_user_id": 2, "is_like": true}`),
			sessionCookie:            &http.Cookie{Name: consts.SessionCookie, Value: "valid_session"},
			getUserIDError:           nil,
			processLikeError:         errors.New("some internal error"),
			expectedStatus:           http.StatusInternalServerError,
			expectedMessage:          "Failed to process like/dislike\n",
			getUserIDBySessionCalled: true,
			processLikeCalled:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			likeService := likehandler_mocks.NewMockLikeService(mockCtrl)
			sessionService := likehandler_mocks.NewMockSessionService(mockCtrl)
			handler := NewHandler(likeService, sessionService)

			if tt.getUserIDBySessionCalled {
				sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), gomock.Any()).Return(1, tt.getUserIDError)
			} else {
				sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), gomock.Any()).Times(0)
			}

			if tt.processLikeCalled {
				likeService.EXPECT().ProcessLike(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.processLikeError)
			} else {
				likeService.EXPECT().ProcessLike(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %v, получили %v", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedMessage {
				t.Errorf("ожидалось сообщение %q, получили %q", tt.expectedMessage, w.Body.String())
			}
		})
	}
}
