// Code generated by MockGen. DO NOT EDIT.
// Source: utils.go
//
// Generated by this command:
//
//	mockgen -source=utils.go -destination=../../test/mocks/handlers_utils.go --package=mocks Utils
//

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	users "github.com/Zapharaos/fihub-backend/internal/auth/users"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
	language "golang.org/x/text/language"
)

// MockUtils is a mock of Utils interface.
type MockUtils struct {
	ctrl     *gomock.Controller
	recorder *MockUtilsMockRecorder
	isgomock struct{}
}

// MockUtilsMockRecorder is the mock recorder for MockUtils.
type MockUtilsMockRecorder struct {
	mock *MockUtils
}

// NewMockUtils creates a new mock instance.
func NewMockUtils(ctrl *gomock.Controller) *MockUtils {
	mock := &MockUtils{ctrl: ctrl}
	mock.recorder = &MockUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUtils) EXPECT() *MockUtilsMockRecorder {
	return m.recorder
}

// CheckPermission mocks base method.
func (m *MockUtils) CheckPermission(w http.ResponseWriter, r *http.Request, permission string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPermission", w, r, permission)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckPermission indicates an expected call of CheckPermission.
func (mr *MockUtilsMockRecorder) CheckPermission(w, r, permission any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermission", reflect.TypeOf((*MockUtils)(nil).CheckPermission), w, r, permission)
}

// GetUserFromContext mocks base method.
func (m *MockUtils) GetUserFromContext(r *http.Request) (users.UserWithRoles, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromContext", r)
	ret0, _ := ret[0].(users.UserWithRoles)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetUserFromContext indicates an expected call of GetUserFromContext.
func (mr *MockUtilsMockRecorder) GetUserFromContext(r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromContext", reflect.TypeOf((*MockUtils)(nil).GetUserFromContext), r)
}

// ParseParamBool mocks base method.
func (m *MockUtils) ParseParamBool(w http.ResponseWriter, r *http.Request, key string) (bool, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseParamBool", w, r, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ParseParamBool indicates an expected call of ParseParamBool.
func (mr *MockUtilsMockRecorder) ParseParamBool(w, r, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseParamBool", reflect.TypeOf((*MockUtils)(nil).ParseParamBool), w, r, key)
}

// ParseParamLanguage mocks base method.
func (m *MockUtils) ParseParamLanguage(w http.ResponseWriter, r *http.Request) language.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseParamLanguage", w, r)
	ret0, _ := ret[0].(language.Tag)
	return ret0
}

// ParseParamLanguage indicates an expected call of ParseParamLanguage.
func (mr *MockUtilsMockRecorder) ParseParamLanguage(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseParamLanguage", reflect.TypeOf((*MockUtils)(nil).ParseParamLanguage), w, r)
}

// ParseParamString mocks base method.
func (m *MockUtils) ParseParamString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseParamString", w, r, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ParseParamString indicates an expected call of ParseParamString.
func (mr *MockUtilsMockRecorder) ParseParamString(w, r, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseParamString", reflect.TypeOf((*MockUtils)(nil).ParseParamString), w, r, key)
}

// ParseParamUUID mocks base method.
func (m *MockUtils) ParseParamUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseParamUUID", w, r, key)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ParseParamUUID indicates an expected call of ParseParamUUID.
func (mr *MockUtilsMockRecorder) ParseParamUUID(w, r, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseParamUUID", reflect.TypeOf((*MockUtils)(nil).ParseParamUUID), w, r, key)
}

// ParseUUIDPair mocks base method.
func (m *MockUtils) ParseUUIDPair(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, uuid.UUID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseUUIDPair", w, r, key)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(uuid.UUID)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// ParseUUIDPair indicates an expected call of ParseUUIDPair.
func (mr *MockUtilsMockRecorder) ParseUUIDPair(w, r, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseUUIDPair", reflect.TypeOf((*MockUtils)(nil).ParseUUIDPair), w, r, key)
}

// ReadImage mocks base method.
func (m *MockUtils) ReadImage(w http.ResponseWriter, r *http.Request) ([]byte, string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadImage", w, r)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// ReadImage indicates an expected call of ReadImage.
func (mr *MockUtilsMockRecorder) ReadImage(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadImage", reflect.TypeOf((*MockUtils)(nil).ReadImage), w, r)
}
