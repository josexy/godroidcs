// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: proto/Ping.proto

package protobuf

import (
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
