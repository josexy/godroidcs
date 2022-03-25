// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: proto/Ping.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_Ping_proto protoreflect.FileDescriptor

var file_proto_Ping_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x50, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x13, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x32, 0x36, 0x0a, 0x08, 0x50, 0x69, 0x6e, 0x67, 0x54, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a,
	0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x3a, 0x0a, 0x1b, 0x63, 0x6f, 0x6d,
	0x2e, 0x6a, 0x6f, 0x78, 0x72, 0x61, 0x79, 0x73, 0x2e, 0x67, 0x6f, 0x64, 0x72, 0x6f, 0x69, 0x64,
	0x73, 0x76, 0x72, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x42, 0x0d, 0x50, 0x69, 0x6e, 0x67, 0x54, 0x65,
	0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x50, 0x01, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_Ping_proto_goTypes = []interface{}{
	(*Empty)(nil), // 0: protobuf.Empty
}
var file_proto_Ping_proto_depIdxs = []int32{
	0, // 0: protobuf.PingTest.Ping:input_type -> protobuf.Empty
	0, // 1: protobuf.PingTest.Ping:output_type -> protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_Ping_proto_init() }
func file_proto_Ping_proto_init() {
	if File_proto_Ping_proto != nil {
		return
	}
	file_proto_Message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_Ping_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_Ping_proto_goTypes,
		DependencyIndexes: file_proto_Ping_proto_depIdxs,
	}.Build()
	File_proto_Ping_proto = out.File
	file_proto_Ping_proto_rawDesc = nil
	file_proto_Ping_proto_goTypes = nil
	file_proto_Ping_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PingTestClient is the client API for PingTest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PingTestClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type pingTestClient struct {
	cc grpc.ClientConnInterface
}

func NewPingTestClient(cc grpc.ClientConnInterface) PingTestClient {
	return &pingTestClient{cc}
}

func (c *pingTestClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/protobuf.PingTest/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingTestServer is the server API for PingTest service.
type PingTestServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
}

// UnimplementedPingTestServer can be embedded to have forward compatible implementations.
type UnimplementedPingTestServer struct {
}

func (*UnimplementedPingTestServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}

func RegisterPingTestServer(s *grpc.Server, srv PingTestServer) {
	s.RegisterService(&_PingTest_serviceDesc, srv)
}

func _PingTest_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingTestServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.PingTest/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingTestServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _PingTest_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.PingTest",
	HandlerType: (*PingTestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _PingTest_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/Ping.proto",
}
