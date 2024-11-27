package updateprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
)

////go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=updateprofile_mocks . ProfileService
//type ProfileService interface {
//	UpdateProfile(ctx context.Context, id int, profile models.Profile) error
//}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=updateprofile_mocks . SessionService
type SessionService interface {
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

////go:generate mockgen -destination=./mocks/mock_UserService.go -package=updateprofile_mocks . UserService
//type UserService interface {
//	GetProfileIdByUserId(ctx context.Context, userId int) (int, error)
//}

type PersonalitiesClient interface {
	UpdateProfile(ctx context.Context,
		in *generatedPersonalities.UpdateProfileRequest) (*generatedPersonalities.UpdateProfileResponse, error)
	GetProfileIdByUserId(ctx context.Context,
		in *generatedPersonalities.GetProfileIDByUserIDRequest) (*generatedPersonalities.GetProfileIDByUserIDResponse, error)
}

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

type Handler struct {
	personalitiesClient generatedPersonalities.PersonalitiesClient
	sessionClient       generatedAuth.AuthClient
	logger              *zap.Logger
}

func NewHandler(personalitiesClient generatedPersonalities.PersonalitiesClient, sessionClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{personalitiesClient: personalitiesClient, sessionClient: sessionClient, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("error getting session cookie", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("error getting user id", zap.Error(err))
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	getProfileIdRequest := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: userId.UserId}
	profileId, err := h.personalitiesClient.GetProfileIDByUserID(ctx, getProfileIdRequest)
	if err != nil {
		h.logger.Error("error getting profile id", zap.Error(err))
		http.Error(w, "profile not found", http.StatusUnauthorized)
		return
	}
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		h.logger.Error("error decoding profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	profile.ID = int(profileId.ProfileID)
	genProfile := &generatedPersonalities.Profile{
		ID:        int32(profile.ID),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       int32(profile.Age),
		Gender:    profile.Gender,
		Target:    profile.Target,
		About:     profile.About,
	}
	updateProfileRequest := &generatedPersonalities.UpdateProfileRequest{
		Id:      profileId.ProfileID,
		Profile: genProfile,
	}
	_, err = h.personalitiesClient.UpdateProfile(ctx, updateProfileRequest)
	if err != nil {
		h.logger.Error("error updating profile", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("profile updated sucessfully")
	fmt.Fprintf(w, "ok")
}
