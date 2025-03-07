// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../../test/mocks/users_password_repository.go --package=mocks -mock_names=Repository=UsersPasswordRepository Repository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	password "github.com/Zapharaos/fihub-backend/internal/auth/password"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// UsersPasswordRepository is a mock of Repository interface.
type UsersPasswordRepository struct {
	ctrl     *gomock.Controller
	recorder *UsersPasswordRepositoryMockRecorder
	isgomock struct{}
}

// UsersPasswordRepositoryMockRecorder is the mock recorder for UsersPasswordRepository.
type UsersPasswordRepositoryMockRecorder struct {
	mock *UsersPasswordRepository
}

// NewUsersPasswordRepository creates a new mock instance.
func NewUsersPasswordRepository(ctrl *gomock.Controller) *UsersPasswordRepository {
	mock := &UsersPasswordRepository{ctrl: ctrl}
	mock.recorder = &UsersPasswordRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *UsersPasswordRepository) EXPECT() *UsersPasswordRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *UsersPasswordRepository) Create(request password.Request) (password.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", request)
	ret0, _ := ret[0].(password.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *UsersPasswordRepositoryMockRecorder) Create(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*UsersPasswordRepository)(nil).Create), request)
}

// Delete mocks base method.
func (m *UsersPasswordRepository) Delete(requestID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", requestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *UsersPasswordRepositoryMockRecorder) Delete(requestID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*UsersPasswordRepository)(nil).Delete), requestID)
}

// GetExpiresAt mocks base method.
func (m *UsersPasswordRepository) GetExpiresAt(userID uuid.UUID) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiresAt", userID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpiresAt indicates an expected call of GetExpiresAt.
func (mr *UsersPasswordRepositoryMockRecorder) GetExpiresAt(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiresAt", reflect.TypeOf((*UsersPasswordRepository)(nil).GetExpiresAt), userID)
}

// GetRequestID mocks base method.
func (m *UsersPasswordRepository) GetRequestID(userID uuid.UUID, token string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequestID", userID, token)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRequestID indicates an expected call of GetRequestID.
func (mr *UsersPasswordRepositoryMockRecorder) GetRequestID(userID, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestID", reflect.TypeOf((*UsersPasswordRepository)(nil).GetRequestID), userID, token)
}

// Valid mocks base method.
func (m *UsersPasswordRepository) Valid(userID, requestID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Valid", userID, requestID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Valid indicates an expected call of Valid.
func (mr *UsersPasswordRepositoryMockRecorder) Valid(userID, requestID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Valid", reflect.TypeOf((*UsersPasswordRepository)(nil).Valid), userID, requestID)
}

// ValidForUser mocks base method.
func (m *UsersPasswordRepository) ValidForUser(userID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidForUser", userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidForUser indicates an expected call of ValidForUser.
func (mr *UsersPasswordRepositoryMockRecorder) ValidForUser(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidForUser", reflect.TypeOf((*UsersPasswordRepository)(nil).ValidForUser), userID)
}
