// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package lmt

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// LakeMetadataTrackerClient is the client API for LakeMetadataTracker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LakeMetadataTrackerClient interface {
	PutEntry(ctx context.Context, in *PutEntryRequest, opts ...grpc.CallOption) (*PutEntryReply, error)
	GetEntry(ctx context.Context, in *GetEntryRequest, opts ...grpc.CallOption) (*GetEntryReply, error)
	DeleteEntry(ctx context.Context, in *DeleteEntryRequest, opts ...grpc.CallOption) (*DeleteEntryReply, error)
}

type lakeMetadataTrackerClient struct {
	cc grpc.ClientConnInterface
}

func NewLakeMetadataTrackerClient(cc grpc.ClientConnInterface) LakeMetadataTrackerClient {
	return &lakeMetadataTrackerClient{cc}
}

func (c *lakeMetadataTrackerClient) PutEntry(ctx context.Context, in *PutEntryRequest, opts ...grpc.CallOption) (*PutEntryReply, error) {
	out := new(PutEntryReply)
	err := c.cc.Invoke(ctx, "/lakefs.v1.LakeMetadataTracker/PutEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lakeMetadataTrackerClient) GetEntry(ctx context.Context, in *GetEntryRequest, opts ...grpc.CallOption) (*GetEntryReply, error) {
	out := new(GetEntryReply)
	err := c.cc.Invoke(ctx, "/lakefs.v1.LakeMetadataTracker/GetEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lakeMetadataTrackerClient) DeleteEntry(ctx context.Context, in *DeleteEntryRequest, opts ...grpc.CallOption) (*DeleteEntryReply, error) {
	out := new(DeleteEntryReply)
	err := c.cc.Invoke(ctx, "/lakefs.v1.LakeMetadataTracker/DeleteEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LakeMetadataTrackerServer is the server API for LakeMetadataTracker service.
// All implementations must embed UnimplementedLakeMetadataTrackerServer
// for forward compatibility
type LakeMetadataTrackerServer interface {
	PutEntry(context.Context, *PutEntryRequest) (*PutEntryReply, error)
	GetEntry(context.Context, *GetEntryRequest) (*GetEntryReply, error)
	DeleteEntry(context.Context, *DeleteEntryRequest) (*DeleteEntryReply, error)
	mustEmbedUnimplementedLakeMetadataTrackerServer()
}

// UnimplementedLakeMetadataTrackerServer must be embedded to have forward compatible implementations.
type UnimplementedLakeMetadataTrackerServer struct {
}

func (UnimplementedLakeMetadataTrackerServer) PutEntry(context.Context, *PutEntryRequest) (*PutEntryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutEntry not implemented")
}
func (UnimplementedLakeMetadataTrackerServer) GetEntry(context.Context, *GetEntryRequest) (*GetEntryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEntry not implemented")
}
func (UnimplementedLakeMetadataTrackerServer) DeleteEntry(context.Context, *DeleteEntryRequest) (*DeleteEntryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEntry not implemented")
}
func (UnimplementedLakeMetadataTrackerServer) mustEmbedUnimplementedLakeMetadataTrackerServer() {}

// UnsafeLakeMetadataTrackerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LakeMetadataTrackerServer will
// result in compilation errors.
type UnsafeLakeMetadataTrackerServer interface {
	mustEmbedUnimplementedLakeMetadataTrackerServer()
}

func RegisterLakeMetadataTrackerServer(s grpc.ServiceRegistrar, srv LakeMetadataTrackerServer) {
	s.RegisterService(&_LakeMetadataTracker_serviceDesc, srv)
}

func _LakeMetadataTracker_PutEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LakeMetadataTrackerServer).PutEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lakefs.v1.LakeMetadataTracker/PutEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LakeMetadataTrackerServer).PutEntry(ctx, req.(*PutEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LakeMetadataTracker_GetEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LakeMetadataTrackerServer).GetEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lakefs.v1.LakeMetadataTracker/GetEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LakeMetadataTrackerServer).GetEntry(ctx, req.(*GetEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LakeMetadataTracker_DeleteEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LakeMetadataTrackerServer).DeleteEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lakefs.v1.LakeMetadataTracker/DeleteEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LakeMetadataTrackerServer).DeleteEntry(ctx, req.(*DeleteEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LakeMetadataTracker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lakefs.v1.LakeMetadataTracker",
	HandlerType: (*LakeMetadataTrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutEntry",
			Handler:    _LakeMetadataTracker_PutEntry_Handler,
		},
		{
			MethodName: "GetEntry",
			Handler:    _LakeMetadataTracker_GetEntry_Handler,
		},
		{
			MethodName: "DeleteEntry",
			Handler:    _LakeMetadataTracker_DeleteEntry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lmt/lmt.proto",
}
