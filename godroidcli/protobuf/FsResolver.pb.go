// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: proto/FsResolver.proto

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

var File_proto_FsResolver_proto protoreflect.FileDescriptor

var file_proto_FsResolver_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x46, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x1a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x83, 0x06, 0x0a, 0x0a, 0x46, 0x73, 0x52, 0x65,
	0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x12, 0x3c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x54, 0x75, 0x70, 0x6c, 0x65,
	0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73, 0x1a,
	0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x28, 0x01, 0x12, 0x3c, 0x0a, 0x13, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x39, 0x0a, 0x07, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x72, 0x12, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x50, 0x61, 0x69, 0x72, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x32,
	0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0x00, 0x12, 0x32, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x05, 0x4d, 0x6b, 0x44, 0x69, 0x72, 0x12,
	0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x05, 0x52, 0x6d, 0x44, 0x69, 0x72, 0x12, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61,
	0x69, 0x72, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x52, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04, 0x43, 0x6f,
	0x70, 0x79, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x08,
	0x52, 0x65, 0x61, 0x64, 0x54, 0x65, 0x78, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x35,
	0x0a, 0x09, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x69,
	0x72, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x54,
	0x65, 0x78, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x42, 0x40, 0x0a,
	0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x6a, 0x6f, 0x78, 0x72, 0x61, 0x79, 0x73, 0x2e, 0x67, 0x6f, 0x64,
	0x72, 0x6f, 0x69, 0x64, 0x73, 0x76, 0x72, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72,
	0x42, 0x0f, 0x46, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x50, 0x01, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_FsResolver_proto_goTypes = []interface{}{
	(*StringTuple)(nil),  // 0: protobuf.StringTuple
	(*ParamBytes)(nil),   // 1: protobuf.ParamBytes
	(*String)(nil),       // 2: protobuf.String
	(*StringPair)(nil),   // 3: protobuf.StringPair
	(*Status)(nil),       // 4: protobuf.Status
	(*Bytes)(nil),        // 5: protobuf.Bytes
	(*FileInfoList)(nil), // 6: protobuf.FileInfoList
}
var file_proto_FsResolver_proto_depIdxs = []int32{
	0,  // 0: protobuf.FsResolver.GetBaseFileTree:input_type -> protobuf.StringTuple
	1,  // 1: protobuf.FsResolver.UploadGeneralFile:input_type -> protobuf.ParamBytes
	2,  // 2: protobuf.FsResolver.DownloadGeneralFile:input_type -> protobuf.String
	3,  // 3: protobuf.FsResolver.ListDir:input_type -> protobuf.StringPair
	2,  // 4: protobuf.FsResolver.DeleteFile:input_type -> protobuf.String
	2,  // 5: protobuf.FsResolver.CreateFile:input_type -> protobuf.String
	2,  // 6: protobuf.FsResolver.MkDir:input_type -> protobuf.String
	2,  // 7: protobuf.FsResolver.RmDir:input_type -> protobuf.String
	3,  // 8: protobuf.FsResolver.Move:input_type -> protobuf.StringPair
	3,  // 9: protobuf.FsResolver.Rename:input_type -> protobuf.StringPair
	3,  // 10: protobuf.FsResolver.Copy:input_type -> protobuf.StringPair
	2,  // 11: protobuf.FsResolver.ReadText:input_type -> protobuf.String
	3,  // 12: protobuf.FsResolver.WriteText:input_type -> protobuf.StringPair
	3,  // 13: protobuf.FsResolver.AppendText:input_type -> protobuf.StringPair
	2,  // 14: protobuf.FsResolver.GetBaseFileTree:output_type -> protobuf.String
	4,  // 15: protobuf.FsResolver.UploadGeneralFile:output_type -> protobuf.Status
	5,  // 16: protobuf.FsResolver.DownloadGeneralFile:output_type -> protobuf.Bytes
	6,  // 17: protobuf.FsResolver.ListDir:output_type -> protobuf.FileInfoList
	4,  // 18: protobuf.FsResolver.DeleteFile:output_type -> protobuf.Status
	4,  // 19: protobuf.FsResolver.CreateFile:output_type -> protobuf.Status
	4,  // 20: protobuf.FsResolver.MkDir:output_type -> protobuf.Status
	4,  // 21: protobuf.FsResolver.RmDir:output_type -> protobuf.Status
	4,  // 22: protobuf.FsResolver.Move:output_type -> protobuf.Status
	4,  // 23: protobuf.FsResolver.Rename:output_type -> protobuf.Status
	4,  // 24: protobuf.FsResolver.Copy:output_type -> protobuf.Status
	4,  // 25: protobuf.FsResolver.ReadText:output_type -> protobuf.Status
	4,  // 26: protobuf.FsResolver.WriteText:output_type -> protobuf.Status
	4,  // 27: protobuf.FsResolver.AppendText:output_type -> protobuf.Status
	14, // [14:28] is the sub-list for method output_type
	0,  // [0:14] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_proto_FsResolver_proto_init() }
func file_proto_FsResolver_proto_init() {
	if File_proto_FsResolver_proto != nil {
		return
	}
	file_proto_Message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_FsResolver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_FsResolver_proto_goTypes,
		DependencyIndexes: file_proto_FsResolver_proto_depIdxs,
	}.Build()
	File_proto_FsResolver_proto = out.File
	file_proto_FsResolver_proto_rawDesc = nil
	file_proto_FsResolver_proto_goTypes = nil
	file_proto_FsResolver_proto_depIdxs = nil
}
