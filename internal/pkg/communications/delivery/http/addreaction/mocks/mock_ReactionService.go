// Code generated by MockGen. DO NOT EDIT.
// Source: sparkit/internal/handlers/addreaction (interfaces: ReactionService)

// Package sign_up_mocks is a generated GoMock package.
package sign_up_mocks

import (
	context "context"
	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReactionService is a mock of ReactionService interface.
type MockReactionService struct {
	ctrl     *gomock.Controller
	recorder *MockReactionServiceMockRecorder
}

// MockReactionServiceMockRecorder is the mock recorder for MockReactionService.
type MockReactionServiceMockRecorder struct {
	mock *MockReactionService
}

// NewMockReactionService creates a new mock instance.
func NewMockReactionService(ctrl *gomock.Controller) *MockReactionService {
	mock := &MockReactionService{ctrl: ctrl}
	mock.recorder = &MockReactionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReactionService) EXPECT() *MockReactionServiceMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockReactionService) AddReaction(arg0 context.Context, arg1 models.Reaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockReactionServiceMockRecorder) AddReaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockReactionService)(nil).AddReaction), arg0, arg1)
}
