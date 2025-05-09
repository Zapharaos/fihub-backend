// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: security_public.proto

package securitypb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CheckPermissionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Permission    string                 `protobuf:"bytes,2,opt,name=permission,proto3" json:"permission,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckPermissionRequest) Reset() {
	*x = CheckPermissionRequest{}
	mi := &file_security_public_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckPermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPermissionRequest) ProtoMessage() {}

func (x *CheckPermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_security_public_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPermissionRequest.ProtoReflect.Descriptor instead.
func (*CheckPermissionRequest) Descriptor() ([]byte, []int) {
	return file_security_public_proto_rawDescGZIP(), []int{0}
}

func (x *CheckPermissionRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CheckPermissionRequest) GetPermission() string {
	if x != nil {
		return x.Permission
	}
	return ""
}

type CheckPermissionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	HasPermission bool                   `protobuf:"varint,1,opt,name=has_permission,json=hasPermission,proto3" json:"has_permission,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckPermissionResponse) Reset() {
	*x = CheckPermissionResponse{}
	mi := &file_security_public_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckPermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPermissionResponse) ProtoMessage() {}

func (x *CheckPermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_security_public_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPermissionResponse.ProtoReflect.Descriptor instead.
func (*CheckPermissionResponse) Descriptor() ([]byte, []int) {
	return file_security_public_proto_rawDescGZIP(), []int{1}
}

func (x *CheckPermissionResponse) GetHasPermission() bool {
	if x != nil {
		return x.HasPermission
	}
	return false
}

var File_security_public_proto protoreflect.FileDescriptor

const file_security_public_proto_rawDesc = "" +
	"\n" +
	"\x15security_public.proto\x12\bsecurity\"Q\n" +
	"\x16CheckPermissionRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1e\n" +
	"\n" +
	"permission\x18\x02 \x01(\tR\n" +
	"permission\"@\n" +
	"\x17CheckPermissionResponse\x12%\n" +
	"\x0ehas_permission\x18\x01 \x01(\bR\rhasPermission2o\n" +
	"\x15PublicSecurityService\x12V\n" +
	"\x0fCheckPermission\x12 .security.CheckPermissionRequest\x1a!.security.CheckPermissionResponseB\x0eZ\f./securitypbb\x06proto3"

var (
	file_security_public_proto_rawDescOnce sync.Once
	file_security_public_proto_rawDescData []byte
)

func file_security_public_proto_rawDescGZIP() []byte {
	file_security_public_proto_rawDescOnce.Do(func() {
		file_security_public_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_security_public_proto_rawDesc), len(file_security_public_proto_rawDesc)))
	})
	return file_security_public_proto_rawDescData
}

var file_security_public_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_security_public_proto_goTypes = []any{
	(*CheckPermissionRequest)(nil),  // 0: security.CheckPermissionRequest
	(*CheckPermissionResponse)(nil), // 1: security.CheckPermissionResponse
}
var file_security_public_proto_depIdxs = []int32{
	0, // 0: security.PublicSecurityService.CheckPermission:input_type -> security.CheckPermissionRequest
	1, // 1: security.PublicSecurityService.CheckPermission:output_type -> security.CheckPermissionResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_security_public_proto_init() }
func file_security_public_proto_init() {
	if File_security_public_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_security_public_proto_rawDesc), len(file_security_public_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_security_public_proto_goTypes,
		DependencyIndexes: file_security_public_proto_depIdxs,
		MessageInfos:      file_security_public_proto_msgTypes,
	}.Build()
	File_security_public_proto = out.File
	file_security_public_proto_goTypes = nil
	file_security_public_proto_depIdxs = nil
}
