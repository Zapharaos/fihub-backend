// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: security_public.proto

package securitypb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PublicSecurityService_CheckPermission_FullMethodName = "/security.PublicSecurityService/CheckPermission"
)

// PublicSecurityServiceClient is the client API for PublicSecurityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublicSecurityServiceClient interface {
	// Public - Security management
	CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error)
}

type publicSecurityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPublicSecurityServiceClient(cc grpc.ClientConnInterface) PublicSecurityServiceClient {
	return &publicSecurityServiceClient{cc}
}

func (c *publicSecurityServiceClient) CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckPermissionResponse)
	err := c.cc.Invoke(ctx, PublicSecurityService_CheckPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublicSecurityServiceServer is the server API for PublicSecurityService service.
// All implementations must embed UnimplementedPublicSecurityServiceServer
// for forward compatibility.
type PublicSecurityServiceServer interface {
	// Public - Security management
	CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error)
	mustEmbedUnimplementedPublicSecurityServiceServer()
}

// UnimplementedPublicSecurityServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPublicSecurityServiceServer struct{}

func (UnimplementedPublicSecurityServiceServer) CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPermission not implemented")
}
func (UnimplementedPublicSecurityServiceServer) mustEmbedUnimplementedPublicSecurityServiceServer() {}
func (UnimplementedPublicSecurityServiceServer) testEmbeddedByValue()                               {}

// UnsafePublicSecurityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublicSecurityServiceServer will
// result in compilation errors.
type UnsafePublicSecurityServiceServer interface {
	mustEmbedUnimplementedPublicSecurityServiceServer()
}

func RegisterPublicSecurityServiceServer(s grpc.ServiceRegistrar, srv PublicSecurityServiceServer) {
	// If the following call pancis, it indicates UnimplementedPublicSecurityServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PublicSecurityService_ServiceDesc, srv)
}

func _PublicSecurityService_CheckPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicSecurityServiceServer).CheckPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PublicSecurityService_CheckPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicSecurityServiceServer).CheckPermission(ctx, req.(*CheckPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PublicSecurityService_ServiceDesc is the grpc.ServiceDesc for PublicSecurityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PublicSecurityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "security.PublicSecurityService",
	HandlerType: (*PublicSecurityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckPermission",
			Handler:    _PublicSecurityService_CheckPermission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "security_public.proto",
}
