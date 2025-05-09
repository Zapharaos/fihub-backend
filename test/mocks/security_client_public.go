// Code generated by MockGen. DO NOT EDIT.
// Source: ../gen/go/securitypb/security_public_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=../gen/go/securitypb/security_public_grpc.pb.go -destination=../test/mocks/security_client_public.go -package=mocks PublicSecurityServiceClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	securitypb "github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockPublicSecurityServiceClient is a mock of PublicSecurityServiceClient interface.
type MockPublicSecurityServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockPublicSecurityServiceClientMockRecorder
	isgomock struct{}
}

// MockPublicSecurityServiceClientMockRecorder is the mock recorder for MockPublicSecurityServiceClient.
type MockPublicSecurityServiceClientMockRecorder struct {
	mock *MockPublicSecurityServiceClient
}

// NewMockPublicSecurityServiceClient creates a new mock instance.
func NewMockPublicSecurityServiceClient(ctrl *gomock.Controller) *MockPublicSecurityServiceClient {
	mock := &MockPublicSecurityServiceClient{ctrl: ctrl}
	mock.recorder = &MockPublicSecurityServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublicSecurityServiceClient) EXPECT() *MockPublicSecurityServiceClientMockRecorder {
	return m.recorder
}

// CheckPermission mocks base method.
func (m *MockPublicSecurityServiceClient) CheckPermission(ctx context.Context, in *securitypb.CheckPermissionRequest, opts ...grpc.CallOption) (*securitypb.CheckPermissionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckPermission", varargs...)
	ret0, _ := ret[0].(*securitypb.CheckPermissionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPermission indicates an expected call of CheckPermission.
func (mr *MockPublicSecurityServiceClientMockRecorder) CheckPermission(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermission", reflect.TypeOf((*MockPublicSecurityServiceClient)(nil).CheckPermission), varargs...)
}

// MockPublicSecurityServiceServer is a mock of PublicSecurityServiceServer interface.
type MockPublicSecurityServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockPublicSecurityServiceServerMockRecorder
	isgomock struct{}
}

// MockPublicSecurityServiceServerMockRecorder is the mock recorder for MockPublicSecurityServiceServer.
type MockPublicSecurityServiceServerMockRecorder struct {
	mock *MockPublicSecurityServiceServer
}

// NewMockPublicSecurityServiceServer creates a new mock instance.
func NewMockPublicSecurityServiceServer(ctrl *gomock.Controller) *MockPublicSecurityServiceServer {
	mock := &MockPublicSecurityServiceServer{ctrl: ctrl}
	mock.recorder = &MockPublicSecurityServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublicSecurityServiceServer) EXPECT() *MockPublicSecurityServiceServerMockRecorder {
	return m.recorder
}

// CheckPermission mocks base method.
func (m *MockPublicSecurityServiceServer) CheckPermission(arg0 context.Context, arg1 *securitypb.CheckPermissionRequest) (*securitypb.CheckPermissionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPermission", arg0, arg1)
	ret0, _ := ret[0].(*securitypb.CheckPermissionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPermission indicates an expected call of CheckPermission.
func (mr *MockPublicSecurityServiceServerMockRecorder) CheckPermission(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermission", reflect.TypeOf((*MockPublicSecurityServiceServer)(nil).CheckPermission), arg0, arg1)
}

// mustEmbedUnimplementedPublicSecurityServiceServer mocks base method.
func (m *MockPublicSecurityServiceServer) mustEmbedUnimplementedPublicSecurityServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedPublicSecurityServiceServer")
}

// mustEmbedUnimplementedPublicSecurityServiceServer indicates an expected call of mustEmbedUnimplementedPublicSecurityServiceServer.
func (mr *MockPublicSecurityServiceServerMockRecorder) mustEmbedUnimplementedPublicSecurityServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedPublicSecurityServiceServer", reflect.TypeOf((*MockPublicSecurityServiceServer)(nil).mustEmbedUnimplementedPublicSecurityServiceServer))
}

// MockUnsafePublicSecurityServiceServer is a mock of UnsafePublicSecurityServiceServer interface.
type MockUnsafePublicSecurityServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafePublicSecurityServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafePublicSecurityServiceServerMockRecorder is the mock recorder for MockUnsafePublicSecurityServiceServer.
type MockUnsafePublicSecurityServiceServerMockRecorder struct {
	mock *MockUnsafePublicSecurityServiceServer
}

// NewMockUnsafePublicSecurityServiceServer creates a new mock instance.
func NewMockUnsafePublicSecurityServiceServer(ctrl *gomock.Controller) *MockUnsafePublicSecurityServiceServer {
	mock := &MockUnsafePublicSecurityServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafePublicSecurityServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafePublicSecurityServiceServer) EXPECT() *MockUnsafePublicSecurityServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedPublicSecurityServiceServer mocks base method.
func (m *MockUnsafePublicSecurityServiceServer) mustEmbedUnimplementedPublicSecurityServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedPublicSecurityServiceServer")
}

// mustEmbedUnimplementedPublicSecurityServiceServer indicates an expected call of mustEmbedUnimplementedPublicSecurityServiceServer.
func (mr *MockUnsafePublicSecurityServiceServerMockRecorder) mustEmbedUnimplementedPublicSecurityServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedPublicSecurityServiceServer", reflect.TypeOf((*MockUnsafePublicSecurityServiceServer)(nil).mustEmbedUnimplementedPublicSecurityServiceServer))
}
