// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/task.go

// Package mock_usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/skantay/todo-list/internal/entity"
)

// MocktaskRepo is a mock of taskRepo interface.
type MocktaskRepo struct {
	ctrl     *gomock.Controller
	recorder *MocktaskRepoMockRecorder
}

// MocktaskRepoMockRecorder is the mock recorder for MocktaskRepo.
type MocktaskRepoMockRecorder struct {
	mock *MocktaskRepo
}

// NewMocktaskRepo creates a new mock instance.
func NewMocktaskRepo(ctrl *gomock.Controller) *MocktaskRepo {
	mock := &MocktaskRepo{ctrl: ctrl}
	mock.recorder = &MocktaskRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktaskRepo) EXPECT() *MocktaskRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MocktaskRepo) Create(ctx context.Context, task entity.Task) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, task)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MocktaskRepoMockRecorder) Create(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MocktaskRepo)(nil).Create), ctx, task)
}

// Delete mocks base method.
func (m *MocktaskRepo) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MocktaskRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MocktaskRepo)(nil).Delete), ctx, id)
}

// List mocks base method.
func (m *MocktaskRepo) List(ctx context.Context, status string, now time.Time) ([]entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, status, now)
	ret0, _ := ret[0].([]entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MocktaskRepoMockRecorder) List(ctx, status, now interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MocktaskRepo)(nil).List), ctx, status, now)
}

// MarkDone mocks base method.
func (m *MocktaskRepo) MarkDone(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDone", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkDone indicates an expected call of MarkDone.
func (mr *MocktaskRepoMockRecorder) MarkDone(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDone", reflect.TypeOf((*MocktaskRepo)(nil).MarkDone), ctx, id)
}

// Update mocks base method.
func (m *MocktaskRepo) Update(ctx context.Context, task entity.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MocktaskRepoMockRecorder) Update(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MocktaskRepo)(nil).Update), ctx, task)
}
