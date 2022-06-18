// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ContactResolverClient is the client API for ContactResolver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContactResolverClient interface {
	GetAllContactInfo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ContactMetaInfoList, error)
	GetContactInfo(ctx context.Context, in *String, opts ...grpc.CallOption) (*ContactInfo, error)
	DeleteContact(ctx context.Context, in *String, opts ...grpc.CallOption) (*Empty, error)
	AddContact(ctx context.Context, in *ContactInfo, opts ...grpc.CallOption) (*Empty, error)
}

type contactResolverClient struct {
	cc grpc.ClientConnInterface
}

func NewContactResolverClient(cc grpc.ClientConnInterface) ContactResolverClient {
	return &contactResolverClient{cc}
}

func (c *contactResolverClient) GetAllContactInfo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ContactMetaInfoList, error) {
	out := new(ContactMetaInfoList)
	err := c.cc.Invoke(ctx, "/protobuf.ContactResolver/GetAllContactInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contactResolverClient) GetContactInfo(ctx context.Context, in *String, opts ...grpc.CallOption) (*ContactInfo, error) {
	out := new(ContactInfo)
	err := c.cc.Invoke(ctx, "/protobuf.ContactResolver/GetContactInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contactResolverClient) DeleteContact(ctx context.Context, in *String, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/protobuf.ContactResolver/DeleteContact", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contactResolverClient) AddContact(ctx context.Context, in *ContactInfo, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/protobuf.ContactResolver/AddContact", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContactResolverServer is the server API for ContactResolver service.
// All implementations must embed UnimplementedContactResolverServer
// for forward compatibility
type ContactResolverServer interface {
	GetAllContactInfo(context.Context, *Empty) (*ContactMetaInfoList, error)
	GetContactInfo(context.Context, *String) (*ContactInfo, error)
	DeleteContact(context.Context, *String) (*Empty, error)
	AddContact(context.Context, *ContactInfo) (*Empty, error)
	mustEmbedUnimplementedContactResolverServer()
}

// UnimplementedContactResolverServer must be embedded to have forward compatible implementations.
type UnimplementedContactResolverServer struct {
}

func (UnimplementedContactResolverServer) GetAllContactInfo(context.Context, *Empty) (*ContactMetaInfoList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllContactInfo not implemented")
}
func (UnimplementedContactResolverServer) GetContactInfo(context.Context, *String) (*ContactInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContactInfo not implemented")
}
func (UnimplementedContactResolverServer) DeleteContact(context.Context, *String) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteContact not implemented")
}
func (UnimplementedContactResolverServer) AddContact(context.Context, *ContactInfo) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddContact not implemented")
}
func (UnimplementedContactResolverServer) mustEmbedUnimplementedContactResolverServer() {}

// UnsafeContactResolverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContactResolverServer will
// result in compilation errors.
type UnsafeContactResolverServer interface {
	mustEmbedUnimplementedContactResolverServer()
}

func RegisterContactResolverServer(s grpc.ServiceRegistrar, srv ContactResolverServer) {
	s.RegisterService(&ContactResolver_ServiceDesc, srv)
}

func _ContactResolver_GetAllContactInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactResolverServer).GetAllContactInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ContactResolver/GetAllContactInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactResolverServer).GetAllContactInfo(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContactResolver_GetContactInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactResolverServer).GetContactInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ContactResolver/GetContactInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactResolverServer).GetContactInfo(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContactResolver_DeleteContact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactResolverServer).DeleteContact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ContactResolver/DeleteContact",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactResolverServer).DeleteContact(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContactResolver_AddContact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContactInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactResolverServer).AddContact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ContactResolver/AddContact",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactResolverServer).AddContact(ctx, req.(*ContactInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// ContactResolver_ServiceDesc is the grpc.ServiceDesc for ContactResolver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContactResolver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.ContactResolver",
	HandlerType: (*ContactResolverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllContactInfo",
			Handler:    _ContactResolver_GetAllContactInfo_Handler,
		},
		{
			MethodName: "GetContactInfo",
			Handler:    _ContactResolver_GetContactInfo_Handler,
		},
		{
			MethodName: "DeleteContact",
			Handler:    _ContactResolver_DeleteContact_Handler,
		},
		{
			MethodName: "AddContact",
			Handler:    _ContactResolver_AddContact_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ContactResolver.proto",
}
