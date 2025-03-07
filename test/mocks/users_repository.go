// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../../test/mocks/users_repository.go --package=mocks -mock_names=Repository=UsersRepository Repository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	users "github.com/Zapharaos/fihub-backend/internal/auth/users"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// UsersRepository is a mock of Repository interface.
type UsersRepository struct {
	ctrl     *gomock.Controller
	recorder *UsersRepositoryMockRecorder
	isgomock struct{}
}

// UsersRepositoryMockRecorder is the mock recorder for UsersRepository.
type UsersRepositoryMockRecorder struct {
	mock *UsersRepository
}

// NewUsersRepository creates a new mock instance.
func NewUsersRepository(ctrl *gomock.Controller) *UsersRepository {
	mock := &UsersRepository{ctrl: ctrl}
	mock.recorder = &UsersRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *UsersRepository) EXPECT() *UsersRepositoryMockRecorder {
	return m.recorder
}

// AddUsersRole mocks base method.
func (m *UsersRepository) AddUsersRole(userUUIDs []uuid.UUID, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUsersRole", userUUIDs, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUsersRole indicates an expected call of AddUsersRole.
func (mr *UsersRepositoryMockRecorder) AddUsersRole(userUUIDs, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUsersRole", reflect.TypeOf((*UsersRepository)(nil).AddUsersRole), userUUIDs, id)
}

// Authenticate mocks base method.
func (m *UsersRepository) Authenticate(email, password string) (users.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", email, password)
	ret0, _ := ret[0].(users.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Authenticate indicates an expected call of Authenticate.
func (mr *UsersRepositoryMockRecorder) Authenticate(email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*UsersRepository)(nil).Authenticate), email, password)
}

// Create mocks base method.
func (m *UsersRepository) Create(user users.UserWithPassword) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *UsersRepositoryMockRecorder) Create(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*UsersRepository)(nil).Create), user)
}

// Delete mocks base method.
func (m *UsersRepository) Delete(userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *UsersRepositoryMockRecorder) Delete(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*UsersRepository)(nil).Delete), userID)
}

// Exists mocks base method.
func (m *UsersRepository) Exists(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *UsersRepositoryMockRecorder) Exists(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*UsersRepository)(nil).Exists), email)
}

// Get mocks base method.
func (m *UsersRepository) Get(userID uuid.UUID) (users.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", userID)
	ret0, _ := ret[0].(users.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *UsersRepositoryMockRecorder) Get(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*UsersRepository)(nil).Get), userID)
}

// GetAllWithRoles mocks base method.
func (m *UsersRepository) GetAllWithRoles() ([]users.UserWithRoles, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithRoles")
	ret0, _ := ret[0].([]users.UserWithRoles)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithRoles indicates an expected call of GetAllWithRoles.
func (mr *UsersRepositoryMockRecorder) GetAllWithRoles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithRoles", reflect.TypeOf((*UsersRepository)(nil).GetAllWithRoles))
}

// GetByEmail mocks base method.
func (m *UsersRepository) GetByEmail(email string) (users.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(users.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *UsersRepositoryMockRecorder) GetByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*UsersRepository)(nil).GetByEmail), email)
}

// GetUsersByRoleID mocks base method.
func (m *UsersRepository) GetUsersByRoleID(roleUUID uuid.UUID) ([]users.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByRoleID", roleUUID)
	ret0, _ := ret[0].([]users.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByRoleID indicates an expected call of GetUsersByRoleID.
func (mr *UsersRepositoryMockRecorder) GetUsersByRoleID(roleUUID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByRoleID", reflect.TypeOf((*UsersRepository)(nil).GetUsersByRoleID), roleUUID)
}

// GetWithRoles mocks base method.
func (m *UsersRepository) GetWithRoles(userID uuid.UUID) (users.UserWithRoles, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithRoles", userID)
	ret0, _ := ret[0].(users.UserWithRoles)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithRoles indicates an expected call of GetWithRoles.
func (mr *UsersRepositoryMockRecorder) GetWithRoles(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithRoles", reflect.TypeOf((*UsersRepository)(nil).GetWithRoles), userID)
}

// RemoveUsersRole mocks base method.
func (m *UsersRepository) RemoveUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUsersRole", userUUIDs, roleUUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUsersRole indicates an expected call of RemoveUsersRole.
func (mr *UsersRepositoryMockRecorder) RemoveUsersRole(userUUIDs, roleUUID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUsersRole", reflect.TypeOf((*UsersRepository)(nil).RemoveUsersRole), userUUIDs, roleUUID)
}

// SetUserRoles mocks base method.
func (m *UsersRepository) SetUserRoles(userUUID uuid.UUID, roleUUIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserRoles", userUUID, roleUUIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserRoles indicates an expected call of SetUserRoles.
func (mr *UsersRepositoryMockRecorder) SetUserRoles(userUUID, roleUUIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserRoles", reflect.TypeOf((*UsersRepository)(nil).SetUserRoles), userUUID, roleUUIDs)
}

// Update mocks base method.
func (m *UsersRepository) Update(user users.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *UsersRepositoryMockRecorder) Update(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*UsersRepository)(nil).Update), user)
}

// UpdateWithPassword mocks base method.
func (m *UsersRepository) UpdateWithPassword(user users.UserWithPassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithPassword", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWithPassword indicates an expected call of UpdateWithPassword.
func (mr *UsersRepositoryMockRecorder) UpdateWithPassword(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithPassword", reflect.TypeOf((*UsersRepository)(nil).UpdateWithPassword), user)
}

// UpdateWithRoles mocks base method.
func (m *UsersRepository) UpdateWithRoles(user users.UserWithRoles, roleUUIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithRoles", user, roleUUIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWithRoles indicates an expected call of UpdateWithRoles.
func (mr *UsersRepositoryMockRecorder) UpdateWithRoles(user, roleUUIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithRoles", reflect.TypeOf((*UsersRepository)(nil).UpdateWithRoles), user, roleUUIDs)
}
