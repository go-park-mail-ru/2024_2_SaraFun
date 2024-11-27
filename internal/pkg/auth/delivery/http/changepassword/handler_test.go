package changepassword

import (
	"bytes"
	_ "context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	changepassword_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/changepassword/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    Request
		mockAuth       func(mock *changepassword_mocks.MockAuthClient)
		mockPersonal   func(mock *changepassword_mocks.MockPersonalitiesClient)
		cookie         *http.Cookie
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful password change",
			requestBody: Request{
				CurrentPassword: "currentpass",
				NewPassword:     "newpass",
			},
			mockAuth: func(mock *changepassword_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "valid-session",
					}).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 1}, nil)
			},
			mockPersonal: func(mock *changepassword_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), &generatedPersonalities.GetUsernameByUserIDRequest{
						UserID: 1,
					}).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{Username: "testuser"}, nil)

				mock.EXPECT().
					CheckPassword(gomock.Any(), &generatedPersonalities.CheckPasswordRequest{
						Username: "testuser",
						Password: "currentpass",
					}).
					Return(&generatedPersonalities.CheckPasswordResponse{}, nil)

				mock.EXPECT().
					ChangePassword(gomock.Any(), &generatedPersonalities.ChangePasswordRequest{
						UserID:   1,
						Password: "newpass",
					}).
					Return(&generatedPersonalities.ChangePasswordResponse{}, nil)
			},
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Вы успешно поменяли пароль!",
		},

		{
			name: "Missing session cookie",
			requestBody: Request{
				CurrentPassword: "currentpass",
				NewPassword:     "newpass",
			},
			mockAuth:       func(mock *changepassword_mocks.MockAuthClient) {},
			mockPersonal:   func(mock *changepassword_mocks.MockPersonalitiesClient) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Вы не авторизованы!",
		},
		{
			name: "Incorrect current password",
			requestBody: Request{
				CurrentPassword: "wrongpass",
				NewPassword:     "newpass",
			},
			mockAuth: func(mock *changepassword_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "valid-session",
					}).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 1}, nil)
			},
			mockPersonal: func(mock *changepassword_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), &generatedPersonalities.GetUsernameByUserIDRequest{
						UserID: 1,
					}).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{Username: "testuser"}, nil)

				mock.EXPECT().
					CheckPassword(gomock.Any(), &generatedPersonalities.CheckPasswordRequest{
						Username: "testuser",
						Password: "wrongpass",
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid password"))
			},
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session",
			},
			expectedStatus: http.StatusPreconditionFailed,
			expectedBody:   "Неправильный текущий пароль",
		},
		{
			name: "Change password failed",
			requestBody: Request{
				CurrentPassword: "currentpass",
				NewPassword:     "newpass",
			},
			mockAuth: func(mock *changepassword_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "valid-session",
					}).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{UserId: 1}, nil)
			},
			mockPersonal: func(mock *changepassword_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), &generatedPersonalities.GetUsernameByUserIDRequest{
						UserID: 1,
					}).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{Username: "testuser"}, nil)

				mock.EXPECT().
					CheckPassword(gomock.Any(), &generatedPersonalities.CheckPasswordRequest{
						Username: "testuser",
						Password: "currentpass",
					}).
					Return(&generatedPersonalities.CheckPasswordResponse{}, nil)

				mock.EXPECT().
					ChangePassword(gomock.Any(), &generatedPersonalities.ChangePasswordRequest{
						UserID:   1,
						Password: "newpass",
					}).
					Return(nil, errors.New("failed to change password"))
			},
			cookie: &http.Cookie{
				Name:  consts.SessionCookie,
				Value: "valid-session",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Не получилось поменять пароль :(",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := changepassword_mocks.NewMockAuthClient(ctrl)
			mockPersonal := changepassword_mocks.NewMockPersonalitiesClient(ctrl)

			tt.mockAuth(mockAuth)
			tt.mockPersonal(mockPersonal)

			handler := NewHandler(mockAuth, mockPersonal, zap.NewNop())

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()

			handler.Handle(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
			require.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
