// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: ads.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AdService_Create_FullMethodName           = "/ads.AdService/Create"
	AdService_ChangeStatus_FullMethodName     = "/ads.AdService/ChangeStatus"
	AdService_Update_FullMethodName           = "/ads.AdService/Update"
	AdService_Filter_FullMethodName           = "/ads.AdService/Filter"
	AdService_GetByID_FullMethodName          = "/ads.AdService/GetByID"
	AdService_GetOnlyPublished_FullMethodName = "/ads.AdService/GetOnlyPublished"
	AdService_Delete_FullMethodName           = "/ads.AdService/Delete"
)

// AdServiceClient is the client API for AdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdServiceClient interface {
	Create(ctx context.Context, in *CreateAdRequest, opts ...grpc.CallOption) (*AdResponse, error)
	ChangeStatus(ctx context.Context, in *ChangeAdStatusRequest, opts ...grpc.CallOption) (*AdResponse, error)
	Update(ctx context.Context, in *UpdateAdRequest, opts ...grpc.CallOption) (*AdResponse, error)
	Filter(ctx context.Context, in *FilterAdsRequest, opts ...grpc.CallOption) (*AdsResponse, error)
	GetByID(ctx context.Context, in *GetAdByIDRequest, opts ...grpc.CallOption) (*AdResponse, error)
	GetOnlyPublished(ctx context.Context, in *AdIDsRequest, opts ...grpc.CallOption) (*AdsResponse, error)
	Delete(ctx context.Context, in *DeleteAdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type adServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdServiceClient(cc grpc.ClientConnInterface) AdServiceClient {
	return &adServiceClient{cc}
}

func (c *adServiceClient) Create(ctx context.Context, in *CreateAdRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) ChangeStatus(ctx context.Context, in *ChangeAdStatusRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_ChangeStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) Update(ctx context.Context, in *UpdateAdRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) Filter(ctx context.Context, in *FilterAdsRequest, opts ...grpc.CallOption) (*AdsResponse, error) {
	out := new(AdsResponse)
	err := c.cc.Invoke(ctx, AdService_Filter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) GetByID(ctx context.Context, in *GetAdByIDRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_GetByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) GetOnlyPublished(ctx context.Context, in *AdIDsRequest, opts ...grpc.CallOption) (*AdsResponse, error) {
	out := new(AdsResponse)
	err := c.cc.Invoke(ctx, AdService_GetOnlyPublished_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) Delete(ctx context.Context, in *DeleteAdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AdService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdServiceServer is the server API for AdService service.
// All implementations should embed UnimplementedAdServiceServer
// for forward compatibility
type AdServiceServer interface {
	Create(context.Context, *CreateAdRequest) (*AdResponse, error)
	ChangeStatus(context.Context, *ChangeAdStatusRequest) (*AdResponse, error)
	Update(context.Context, *UpdateAdRequest) (*AdResponse, error)
	Filter(context.Context, *FilterAdsRequest) (*AdsResponse, error)
	GetByID(context.Context, *GetAdByIDRequest) (*AdResponse, error)
	GetOnlyPublished(context.Context, *AdIDsRequest) (*AdsResponse, error)
	Delete(context.Context, *DeleteAdRequest) (*emptypb.Empty, error)
}

// UnimplementedAdServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAdServiceServer struct {
}

func (UnimplementedAdServiceServer) Create(context.Context, *CreateAdRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAdServiceServer) ChangeStatus(context.Context, *ChangeAdStatusRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeStatus not implemented")
}
func (UnimplementedAdServiceServer) Update(context.Context, *UpdateAdRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAdServiceServer) Filter(context.Context, *FilterAdsRequest) (*AdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Filter not implemented")
}
func (UnimplementedAdServiceServer) GetByID(context.Context, *GetAdByIDRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedAdServiceServer) GetOnlyPublished(context.Context, *AdIDsRequest) (*AdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOnlyPublished not implemented")
}
func (UnimplementedAdServiceServer) Delete(context.Context, *DeleteAdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeAdServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdServiceServer will
// result in compilation errors.
type UnsafeAdServiceServer interface {
	mustEmbedUnimplementedAdServiceServer()
}

func RegisterAdServiceServer(s grpc.ServiceRegistrar, srv AdServiceServer) {
	s.RegisterService(&AdService_ServiceDesc, srv)
}

func _AdService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).Create(ctx, req.(*CreateAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_ChangeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeAdStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).ChangeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_ChangeStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).ChangeStatus(ctx, req.(*ChangeAdStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).Update(ctx, req.(*UpdateAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_Filter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilterAdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).Filter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_Filter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).Filter(ctx, req.(*FilterAdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_GetByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).GetByID(ctx, req.(*GetAdByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_GetOnlyPublished_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).GetOnlyPublished(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_GetOnlyPublished_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).GetOnlyPublished(ctx, req.(*AdIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).Delete(ctx, req.(*DeleteAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdService_ServiceDesc is the grpc.ServiceDesc for AdService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ads.AdService",
	HandlerType: (*AdServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AdService_Create_Handler,
		},
		{
			MethodName: "ChangeStatus",
			Handler:    _AdService_ChangeStatus_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _AdService_Update_Handler,
		},
		{
			MethodName: "Filter",
			Handler:    _AdService_Filter_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _AdService_GetByID_Handler,
		},
		{
			MethodName: "GetOnlyPublished",
			Handler:    _AdService_GetOnlyPublished_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _AdService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ads.proto",
}
