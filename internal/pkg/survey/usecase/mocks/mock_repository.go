// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddQuestion mocks base method.
func (m *MockRepository) AddQuestion(ctx context.Context, question models.AdminQuestion) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddQuestion", ctx, question)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddQuestion indicates an expected call of AddQuestion.
func (mr *MockRepositoryMockRecorder) AddQuestion(ctx, question interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddQuestion", reflect.TypeOf((*MockRepository)(nil).AddQuestion), ctx, question)
}

// AddSurvey mocks base method.
func (m *MockRepository) AddSurvey(ctx context.Context, survey models.Survey) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSurvey", ctx, survey)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSurvey indicates an expected call of AddSurvey.
func (mr *MockRepositoryMockRecorder) AddSurvey(ctx, survey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSurvey", reflect.TypeOf((*MockRepository)(nil).AddSurvey), ctx, survey)
}

// DeleteQuestion mocks base method.
func (m *MockRepository) DeleteQuestion(ctx context.Context, content string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuestion", ctx, content)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuestion indicates an expected call of DeleteQuestion.
func (mr *MockRepositoryMockRecorder) DeleteQuestion(ctx, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuestion", reflect.TypeOf((*MockRepository)(nil).DeleteQuestion), ctx, content)
}

// GetQuestions mocks base method.
func (m *MockRepository) GetQuestions(ctx context.Context) ([]models.AdminQuestion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestions", ctx)
	ret0, _ := ret[0].([]models.AdminQuestion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestions indicates an expected call of GetQuestions.
func (mr *MockRepositoryMockRecorder) GetQuestions(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestions", reflect.TypeOf((*MockRepository)(nil).GetQuestions), ctx)
}

// GetSurveyInfo mocks base method.
func (m *MockRepository) GetSurveyInfo(ctx context.Context) ([]models.Survey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSurveyInfo", ctx)
	ret0, _ := ret[0].([]models.Survey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSurveyInfo indicates an expected call of GetSurveyInfo.
func (mr *MockRepositoryMockRecorder) GetSurveyInfo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSurveyInfo", reflect.TypeOf((*MockRepository)(nil).GetSurveyInfo), ctx)
}

// UpdateQuestion mocks base method.
func (m *MockRepository) UpdateQuestion(ct context.Context, question models.AdminQuestion, content string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuestion", ct, question, content)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateQuestion indicates an expected call of UpdateQuestion.
func (mr *MockRepositoryMockRecorder) UpdateQuestion(ct, question, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuestion", reflect.TypeOf((*MockRepository)(nil).UpdateQuestion), ct, question, content)
}
