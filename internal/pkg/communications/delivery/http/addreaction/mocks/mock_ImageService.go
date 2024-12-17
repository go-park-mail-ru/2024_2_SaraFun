// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/addreaction (interfaces: ImageService)

// Package sign_up_mocks is a generated GoMock package.
package sign_up_mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockImageService is a mock of ImageService interface.
type MockImageService struct {
	ctrl     *gomock.Controller
	recorder *MockImageServiceMockRecorder
}

// MockImageServiceMockRecorder is the mock recorder for MockImageService.
type MockImageServiceMockRecorder struct {
	mock *MockImageService
}

// NewMockImageService creates a new mock instance.
func NewMockImageService(ctrl *gomock.Controller) *MockImageService {
	mock := &MockImageService{ctrl: ctrl}
	mock.recorder = &MockImageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageService) EXPECT() *MockImageServiceMockRecorder {
	return m.recorder
}

// GetFirstImage mocks base method.
func (m *MockImageService) GetFirstImage(arg0 context.Context, arg1 int) (models.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFirstImage", arg0, arg1)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFirstImage indicates an expected call of GetFirstImage.
func (mr *MockImageServiceMockRecorder) GetFirstImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstImage", reflect.TypeOf((*MockImageService)(nil).GetFirstImage), arg0, arg1)
}