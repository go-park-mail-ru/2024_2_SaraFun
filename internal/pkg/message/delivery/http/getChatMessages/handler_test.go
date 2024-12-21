package getChatMessages

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	messagemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen/mocks"
	imagemocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/http/getChatMessages/mocks"
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
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	messageClient := messagemocks.NewMockMessageClient(mockCtrl)
	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
	imageService := imagemocks.NewMockImageService(mockCtrl)
	handler := NewHandler(authClient, messageClient, personalitiesClient, imageService, logger)

	successProfile := &generatedPersonalities.Profile{
		ID:        2,
		FirstName: "Имя",
		LastName:  "Фамилия",
		Age:       100,
		Gender:    "пол",
		Target:    "цель",
		About:     "о себе",
		BirthDate: "2000-01-01",
	}
	successMessages := []*generatedMessage.ChatMessage{
		{
			Body: "test",
		},
		{
			Body: "test",
		},
	}
	successImages := []models.Image{
		{
			Link: "link1",
		},
		{
			Link: "link2",
		},
	}

	tests := []struct {
		name                  string
		method                string
		path                  string
		cookieValue           string
		secondUser            int
		authReturn            int
		authError             error
		authTimes             int
		usernameReturn        string
		usernameError         error
		usernameTimes         int
		profileIDReturn       int
		profileIDError        error
		profileIDTimes        int
		getProfileReturn      *generatedPersonalities.Profile
		getProfileError       error
		getProfileTimes       int
		getChatMessagesReturn []*generatedMessage.ChatMessage
		getChatMessagesError  error
		getChatMessageTimes   int
		getImagesLink         []models.Image
		getImagesError        error
		getImagesTimes        int
	}{
		{
			name:                  "good test",
			path:                  "/api/message/getchat?userID=2",
			method:                http.MethodGet,
			cookieValue:           "sparkit",
			secondUser:            2,
			authReturn:            1,
			authError:             nil,
			authTimes:             1,
			profileIDReturn:       1,
			profileIDError:        nil,
			profileIDTimes:        1,
			getProfileReturn:      successProfile,
			getProfileError:       nil,
			getProfileTimes:       1,
			getChatMessagesReturn: successMessages,
			getChatMessagesError:  nil,
			getChatMessageTimes:   1,
			usernameReturn:        "username",
			usernameError:         nil,
			usernameTimes:         1,
			getImagesLink:         successImages,
			getImagesError:        nil,
			getImagesTimes:        1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserIDReq).Return(getUserIDResponse, tt.authError).Times(tt.authTimes)

			getUsernameByUserIDReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: int32(tt.secondUser)}
			getUsernameByUserIDResponse := &generatedPersonalities.GetUsernameByUserIDResponse{Username: tt.usernameReturn}
			personalitiesClient.EXPECT().GetUsernameByUserID(ctx, getUsernameByUserIDReq).Return(getUsernameByUserIDResponse, tt.usernameError).
				Times(tt.usernameTimes)

			getProfileIDReq := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: int32(tt.secondUser)}
			getProfileIDResponse := &generatedPersonalities.GetProfileIDByUserIDResponse{ProfileID: int32(tt.profileIDReturn)}
			personalitiesClient.EXPECT().GetProfileIDByUserID(ctx, getProfileIDReq).Return(getProfileIDResponse, tt.profileIDError).
				Times(tt.profileIDTimes)

			getProfileReq := &generatedPersonalities.GetProfileRequest{Id: int32(tt.profileIDReturn)}
			getProfileResponse := &generatedPersonalities.GetProfileResponse{Profile: tt.getProfileReturn}
			personalitiesClient.EXPECT().GetProfile(ctx, getProfileReq).Return(getProfileResponse, tt.getProfileError).Times(tt.getProfileTimes)

			getChatMessagesReq := &generatedMessage.GetChatMessagesRequest{
				FirstUserID:  int32(tt.authReturn),
				SecondUserID: int32(tt.secondUser),
			}
			getChatMessagesResponse := &generatedMessage.GetChatMessagesResponse{}
			messageClient.EXPECT().GetChatMessages(ctx, getChatMessagesReq).Return(getChatMessagesResponse, tt.getChatMessagesError).
				Times(tt.getChatMessageTimes)
			imageService.EXPECT().GetImageLinksByUserId(ctx, tt.secondUser).Return(tt.getImagesLink, tt.getImagesError).Times(tt.getImagesTimes)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(ctx)
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: tt.cookieValue,
			}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			handler.Handle(w, req)
		})
	}
}
