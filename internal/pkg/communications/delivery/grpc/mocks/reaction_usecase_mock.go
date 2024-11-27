// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockReactionUseCase is a mock of ReactionUseCase interface.
type MockReactionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockReactionUseCaseMockRecorder
}

// MockReactionUseCaseMockRecorder is the mock recorder for MockReactionUseCase.
type MockReactionUseCaseMockRecorder struct {
	mock *MockReactionUseCase
}

// NewMockReactionUseCase creates a new mock instance.
func NewMockReactionUseCase(ctrl *gomock.Controller) *MockReactionUseCase {
	mock := &MockReactionUseCase{ctrl: ctrl}
	mock.recorder = &MockReactionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReactionUseCase) EXPECT() *MockReactionUseCaseMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockReactionUseCase) AddReaction(ctx context.Context, reaction models.Reaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", ctx, reaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockReactionUseCaseMockRecorder) AddReaction(ctx, reaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockReactionUseCase)(nil).AddReaction), ctx, reaction)
}

// GetMatchList mocks base method.
func (m *MockReactionUseCase) GetMatchList(ctx context.Context, userId int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchList", ctx, userId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchList indicates an expected call of GetMatchList.
func (mr *MockReactionUseCaseMockRecorder) GetMatchList(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchList", reflect.TypeOf((*MockReactionUseCase)(nil).GetMatchList), ctx, userId)
}

// GetReactionList mocks base method.
func (m *MockReactionUseCase) GetReactionList(ctx context.Context, userId int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReactionList", ctx, userId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReactionList indicates an expected call of GetReactionList.
func (mr *MockReactionUseCaseMockRecorder) GetReactionList(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReactionList", reflect.TypeOf((*MockReactionUseCase)(nil).GetReactionList), ctx, userId)
}