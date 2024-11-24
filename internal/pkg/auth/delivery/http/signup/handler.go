package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/hashing"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService
//type UserService interface {
//	RegisterUser(ctx context.Context, user models.User) (int, error)
//	CheckUsernameExists(ctx context.Context, username string) (bool, error)
//}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService
//type SessionService interface {
//	CreateSession(ctx context.Context, user models.User) (models.Session, error)
//}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService
//type ProfileService interface {
//	CreateProfile(ctx context.Context, profile models.Profile) (int, error)
//}

//type UserClient interface {
//	RegisterUser(ctx context.Context,
//		in *generatedPersonalities.RegisterUserRequest) (*generatedPersonalities.RegisterUserResponse, error)
//	CheckUsernameExists(ctx context.Context,
//		in *generatedPersonalities.CheckUsernameExistsRequest) (*generatedPersonalities.CheckUsernameExistsResponse, error)
//}

type SessionClient interface {
	CreateSession(ctx context.Context, in *generatedAuth.CreateSessionRequest) (*generatedAuth.CreateSessionResponse, error)
}

//type ProfileClient interface {
//	CreateProfile(ctx context.Context,
//		in *generatedPersonalities.CreateProfileRequest) (*generatedPersonalities.CreateProfileResponse, error)
//}

type PersonalitiesClient interface {
	RegisterUser(ctx context.Context,
		in *generatedPersonalities.RegisterUserRequest) (*generatedPersonalities.RegisterUserResponse, error)
	CheckUsernameExists(ctx context.Context,
		in *generatedPersonalities.CheckUsernameExistsRequest) (*generatedPersonalities.CheckUsernameExistsResponse, error)
	CreateProfile(ctx context.Context,
		in *generatedPersonalities.CreateProfileRequest) (*generatedPersonalities.CreateProfileResponse, error)
}

//type Handler struct {
//	userService    UserService
//	sessionService SessionService
//	profileService ProfileService
//	logger         *zap.Logger
//}

type Handler struct {
	personalitiesClient generatedPersonalities.PersonalitiesClient
	sessionClient       generatedAuth.AuthClient
	logger              *zap.Logger
}

type Request struct {
	User    models.User
	Profile models.Profile
}

func NewHandler(personalitiesClient generatedPersonalities.PersonalitiesClient, sessionsClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{
		personalitiesClient: personalitiesClient,
		sessionClient:       sessionsClient,
		logger:              logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	if r.Method != http.MethodPost {
		h.logger.Error("bad method", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request := Request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.User.Sanitize()
	request.Profile.Sanitize()
	//personalitiesGRPC
	checkUsernameRequest := &generatedPersonalities.CheckUsernameExistsRequest{Username: request.User.Username}
	exists, err := h.personalitiesClient.CheckUsernameExists(ctx, checkUsernameRequest)
	if err != nil {
		h.logger.Error("failed to check username exists", zap.Error(err))
		http.Error(w, "failed to check username exists", http.StatusInternalServerError)
		return
	}
	if exists.Exists {
		h.logger.Error("user already exists", zap.String("username", request.User.Username))
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	//personalitiesGRPC
	genProfile := &generatedPersonalities.Profile{ID: int32(request.Profile.ID),
		FirstName: request.Profile.FirstName,
		LastName:  request.Profile.LastName,
		Age:       int32(request.Profile.Age),
		Gender:    request.Profile.Gender,
		Target:    request.Profile.Target,
		About:     request.Profile.About,
	}
	createProfileRequest := &generatedPersonalities.CreateProfileRequest{Profile: genProfile}
	profileId, err := h.personalitiesClient.CreateProfile(ctx, createProfileRequest)
	if err != nil {
		h.logger.Error("failed to create Profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.User.Profile = int(profileId.ProfileId)
	hashedPass, err := hashing.HashPassword(request.User.Password)
	if err != nil {
		h.logger.Error("failed to hash password", zap.Error(err))
		http.Error(w, "bad password", http.StatusBadRequest)
		return
	}
	request.User.Password = hashedPass
	// personalities grpc
	genUser := &generatedPersonalities.User{
		ID:       int32(request.User.ID),
		Username: request.User.Username,
		Password: request.User.Password,
		Email:    request.User.Email,
		Profile:  int32(request.User.Profile),
	}
	registerUserRequest := &generatedPersonalities.RegisterUserRequest{User: genUser}
	id, err := h.personalitiesClient.RegisterUser(ctx, registerUserRequest)
	if err != nil {
		h.logger.Error("failed to register User", zap.Error(err))
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	request.User.ID = int(id.UserId)

	//auth grpc
	sessUser := &generatedAuth.User{
		ID:       int32(request.User.ID),
		Username: request.User.Username,
		Password: request.User.Password,
		Email:    request.User.Email,
		Profile:  int32(request.User.Profile),
	}
	createSessionRequest := &generatedAuth.CreateSessionRequest{
		User: sessUser,
	}
	if session, err := h.sessionClient.CreateSession(ctx, createSessionRequest); err != nil {
		h.logger.Error("failed to create session", zap.Error(err))
		http.Error(w, "Не удалось создать сессию", http.StatusInternalServerError)
		return
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     consts.SessionCookie,
			Value:    session.Session.SessionID,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
	}
	h.logger.Info("good signup", zap.String("username", request.User.Username))
	fmt.Fprintf(w, "Вы успешно зарегистрировались")
}
