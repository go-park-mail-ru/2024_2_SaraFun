package deletequestion

//import (
//	"context"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//	surveymocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//	"github.com/golang/mock/gomock"
//	"go.uber.org/zap"
//	"testing"
//	"time"
//)
//
//func TestHandler(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	authClient := authmocks.NewMockAuthClient(mockCtrl)
//	surveyClient := surveymocks.NewMockSurveyClient(mockCtrl)
//
//}
