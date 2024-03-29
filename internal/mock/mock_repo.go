// Code generated by MockGen. DO NOT EDIT.
// Source: repo_interface.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	model "github.com/ffajarpratama/boiler-api/internal/model"
	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockIFaceRepository is a mock of IFaceRepository interface.
type MockIFaceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIFaceRepositoryMockRecorder
}

// MockIFaceRepositoryMockRecorder is the mock recorder for MockIFaceRepository.
type MockIFaceRepositoryMockRecorder struct {
	mock *MockIFaceRepository
}

// NewMockIFaceRepository creates a new mock instance.
func NewMockIFaceRepository(ctrl *gomock.Controller) *MockIFaceRepository {
	mock := &MockIFaceRepository{ctrl: ctrl}
	mock.recorder = &MockIFaceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFaceRepository) EXPECT() *MockIFaceRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockIFaceRepository) CreateUser(ctx context.Context, data *model.User, db *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, data, db)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIFaceRepositoryMockRecorder) CreateUser(ctx, data, db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIFaceRepository)(nil).CreateUser), ctx, data, db)
}

// FindOneUser mocks base method.
func (m *MockIFaceRepository) FindOneUser(ctx context.Context, query ...interface{}) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range query {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOneUser", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneUser indicates an expected call of FindOneUser.
func (mr *MockIFaceRepositoryMockRecorder) FindOneUser(ctx interface{}, query ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, query...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneUser", reflect.TypeOf((*MockIFaceRepository)(nil).FindOneUser), varargs...)
}

// Ping mocks base method.
func (m *MockIFaceRepository) Ping(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockIFaceRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockIFaceRepository)(nil).Ping), ctx)
}
