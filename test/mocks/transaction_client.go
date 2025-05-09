// Code generated by MockGen. DO NOT EDIT.
// Source: ../gen/go/transactionpb/transaction_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=../gen/go/transactionpb/transaction_grpc.pb.go -destination=../test/mocks/transaction_client.go -package=mocks TransactionServiceClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	transactionpb "github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockTransactionServiceClient is a mock of TransactionServiceClient interface.
type MockTransactionServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionServiceClientMockRecorder
	isgomock struct{}
}

// MockTransactionServiceClientMockRecorder is the mock recorder for MockTransactionServiceClient.
type MockTransactionServiceClientMockRecorder struct {
	mock *MockTransactionServiceClient
}

// NewMockTransactionServiceClient creates a new mock instance.
func NewMockTransactionServiceClient(ctrl *gomock.Controller) *MockTransactionServiceClient {
	mock := &MockTransactionServiceClient{ctrl: ctrl}
	mock.recorder = &MockTransactionServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionServiceClient) EXPECT() *MockTransactionServiceClientMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *MockTransactionServiceClient) CreateTransaction(ctx context.Context, in *transactionpb.CreateTransactionRequest, opts ...grpc.CallOption) (*transactionpb.CreateTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateTransaction", varargs...)
	ret0, _ := ret[0].(*transactionpb.CreateTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionServiceClientMockRecorder) CreateTransaction(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionServiceClient)(nil).CreateTransaction), varargs...)
}

// DeleteTransaction mocks base method.
func (m *MockTransactionServiceClient) DeleteTransaction(ctx context.Context, in *transactionpb.DeleteTransactionRequest, opts ...grpc.CallOption) (*transactionpb.DeleteTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTransaction", varargs...)
	ret0, _ := ret[0].(*transactionpb.DeleteTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockTransactionServiceClientMockRecorder) DeleteTransaction(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockTransactionServiceClient)(nil).DeleteTransaction), varargs...)
}

// DeleteTransactionByBroker mocks base method.
func (m *MockTransactionServiceClient) DeleteTransactionByBroker(ctx context.Context, in *transactionpb.DeleteTransactionByBrokerRequest, opts ...grpc.CallOption) (*transactionpb.DeleteTransactionByBrokerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTransactionByBroker", varargs...)
	ret0, _ := ret[0].(*transactionpb.DeleteTransactionByBrokerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTransactionByBroker indicates an expected call of DeleteTransactionByBroker.
func (mr *MockTransactionServiceClientMockRecorder) DeleteTransactionByBroker(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransactionByBroker", reflect.TypeOf((*MockTransactionServiceClient)(nil).DeleteTransactionByBroker), varargs...)
}

// GetTransaction mocks base method.
func (m *MockTransactionServiceClient) GetTransaction(ctx context.Context, in *transactionpb.GetTransactionRequest, opts ...grpc.CallOption) (*transactionpb.GetTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTransaction", varargs...)
	ret0, _ := ret[0].(*transactionpb.GetTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockTransactionServiceClientMockRecorder) GetTransaction(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionServiceClient)(nil).GetTransaction), varargs...)
}

// ListTransactions mocks base method.
func (m *MockTransactionServiceClient) ListTransactions(ctx context.Context, in *transactionpb.ListTransactionsRequest, opts ...grpc.CallOption) (*transactionpb.ListTransactionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListTransactions", varargs...)
	ret0, _ := ret[0].(*transactionpb.ListTransactionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransactions indicates an expected call of ListTransactions.
func (mr *MockTransactionServiceClientMockRecorder) ListTransactions(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransactions", reflect.TypeOf((*MockTransactionServiceClient)(nil).ListTransactions), varargs...)
}

// UpdateTransaction mocks base method.
func (m *MockTransactionServiceClient) UpdateTransaction(ctx context.Context, in *transactionpb.UpdateTransactionRequest, opts ...grpc.CallOption) (*transactionpb.UpdateTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateTransaction", varargs...)
	ret0, _ := ret[0].(*transactionpb.UpdateTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockTransactionServiceClientMockRecorder) UpdateTransaction(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionServiceClient)(nil).UpdateTransaction), varargs...)
}

// MockTransactionServiceServer is a mock of TransactionServiceServer interface.
type MockTransactionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionServiceServerMockRecorder
	isgomock struct{}
}

// MockTransactionServiceServerMockRecorder is the mock recorder for MockTransactionServiceServer.
type MockTransactionServiceServerMockRecorder struct {
	mock *MockTransactionServiceServer
}

// NewMockTransactionServiceServer creates a new mock instance.
func NewMockTransactionServiceServer(ctrl *gomock.Controller) *MockTransactionServiceServer {
	mock := &MockTransactionServiceServer{ctrl: ctrl}
	mock.recorder = &MockTransactionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionServiceServer) EXPECT() *MockTransactionServiceServerMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *MockTransactionServiceServer) CreateTransaction(arg0 context.Context, arg1 *transactionpb.CreateTransactionRequest) (*transactionpb.CreateTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.CreateTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionServiceServerMockRecorder) CreateTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionServiceServer)(nil).CreateTransaction), arg0, arg1)
}

// DeleteTransaction mocks base method.
func (m *MockTransactionServiceServer) DeleteTransaction(arg0 context.Context, arg1 *transactionpb.DeleteTransactionRequest) (*transactionpb.DeleteTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.DeleteTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockTransactionServiceServerMockRecorder) DeleteTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockTransactionServiceServer)(nil).DeleteTransaction), arg0, arg1)
}

// DeleteTransactionByBroker mocks base method.
func (m *MockTransactionServiceServer) DeleteTransactionByBroker(arg0 context.Context, arg1 *transactionpb.DeleteTransactionByBrokerRequest) (*transactionpb.DeleteTransactionByBrokerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransactionByBroker", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.DeleteTransactionByBrokerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTransactionByBroker indicates an expected call of DeleteTransactionByBroker.
func (mr *MockTransactionServiceServerMockRecorder) DeleteTransactionByBroker(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransactionByBroker", reflect.TypeOf((*MockTransactionServiceServer)(nil).DeleteTransactionByBroker), arg0, arg1)
}

// GetTransaction mocks base method.
func (m *MockTransactionServiceServer) GetTransaction(arg0 context.Context, arg1 *transactionpb.GetTransactionRequest) (*transactionpb.GetTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.GetTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockTransactionServiceServerMockRecorder) GetTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionServiceServer)(nil).GetTransaction), arg0, arg1)
}

// ListTransactions mocks base method.
func (m *MockTransactionServiceServer) ListTransactions(arg0 context.Context, arg1 *transactionpb.ListTransactionsRequest) (*transactionpb.ListTransactionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransactions", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.ListTransactionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransactions indicates an expected call of ListTransactions.
func (mr *MockTransactionServiceServerMockRecorder) ListTransactions(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransactions", reflect.TypeOf((*MockTransactionServiceServer)(nil).ListTransactions), arg0, arg1)
}

// UpdateTransaction mocks base method.
func (m *MockTransactionServiceServer) UpdateTransaction(arg0 context.Context, arg1 *transactionpb.UpdateTransactionRequest) (*transactionpb.UpdateTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1)
	ret0, _ := ret[0].(*transactionpb.UpdateTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockTransactionServiceServerMockRecorder) UpdateTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionServiceServer)(nil).UpdateTransaction), arg0, arg1)
}

// mustEmbedUnimplementedTransactionServiceServer mocks base method.
func (m *MockTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTransactionServiceServer")
}

// mustEmbedUnimplementedTransactionServiceServer indicates an expected call of mustEmbedUnimplementedTransactionServiceServer.
func (mr *MockTransactionServiceServerMockRecorder) mustEmbedUnimplementedTransactionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTransactionServiceServer", reflect.TypeOf((*MockTransactionServiceServer)(nil).mustEmbedUnimplementedTransactionServiceServer))
}

// MockUnsafeTransactionServiceServer is a mock of UnsafeTransactionServiceServer interface.
type MockUnsafeTransactionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeTransactionServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeTransactionServiceServerMockRecorder is the mock recorder for MockUnsafeTransactionServiceServer.
type MockUnsafeTransactionServiceServerMockRecorder struct {
	mock *MockUnsafeTransactionServiceServer
}

// NewMockUnsafeTransactionServiceServer creates a new mock instance.
func NewMockUnsafeTransactionServiceServer(ctrl *gomock.Controller) *MockUnsafeTransactionServiceServer {
	mock := &MockUnsafeTransactionServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeTransactionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeTransactionServiceServer) EXPECT() *MockUnsafeTransactionServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedTransactionServiceServer mocks base method.
func (m *MockUnsafeTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTransactionServiceServer")
}

// mustEmbedUnimplementedTransactionServiceServer indicates an expected call of mustEmbedUnimplementedTransactionServiceServer.
func (mr *MockUnsafeTransactionServiceServerMockRecorder) mustEmbedUnimplementedTransactionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTransactionServiceServer", reflect.TypeOf((*MockUnsafeTransactionServiceServer)(nil).mustEmbedUnimplementedTransactionServiceServer))
}
