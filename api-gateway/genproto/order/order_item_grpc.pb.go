// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: order_item.proto

package order

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

const (
	OrderItemService_Create_FullMethodName  = "/order.OrderItemService/Create"
	OrderItemService_GetById_FullMethodName = "/order.OrderItemService/GetById"
	OrderItemService_GetAll_FullMethodName  = "/order.OrderItemService/GetAll"
)

// OrderItemServiceClient is the client API for OrderItemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderItemServiceClient interface {
	Create(ctx context.Context, in *OrderItemCreateReq, opts ...grpc.CallOption) (*Void, error)
	GetById(ctx context.Context, in *ById, opts ...grpc.CallOption) (*OrderItemGetByIdRes, error)
	GetAll(ctx context.Context, in *OrderItemGetAllReq, opts ...grpc.CallOption) (*OrderItemGetAllRes, error)
}

type orderItemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderItemServiceClient(cc grpc.ClientConnInterface) OrderItemServiceClient {
	return &orderItemServiceClient{cc}
}

func (c *orderItemServiceClient) Create(ctx context.Context, in *OrderItemCreateReq, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, OrderItemService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderItemServiceClient) GetById(ctx context.Context, in *ById, opts ...grpc.CallOption) (*OrderItemGetByIdRes, error) {
	out := new(OrderItemGetByIdRes)
	err := c.cc.Invoke(ctx, OrderItemService_GetById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderItemServiceClient) GetAll(ctx context.Context, in *OrderItemGetAllReq, opts ...grpc.CallOption) (*OrderItemGetAllRes, error) {
	out := new(OrderItemGetAllRes)
	err := c.cc.Invoke(ctx, OrderItemService_GetAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderItemServiceServer is the server API for OrderItemService service.
// All implementations must embed UnimplementedOrderItemServiceServer
// for forward compatibility
type OrderItemServiceServer interface {
	Create(context.Context, *OrderItemCreateReq) (*Void, error)
	GetById(context.Context, *ById) (*OrderItemGetByIdRes, error)
	GetAll(context.Context, *OrderItemGetAllReq) (*OrderItemGetAllRes, error)
	mustEmbedUnimplementedOrderItemServiceServer()
}

// UnimplementedOrderItemServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrderItemServiceServer struct {
}

func (UnimplementedOrderItemServiceServer) Create(context.Context, *OrderItemCreateReq) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedOrderItemServiceServer) GetById(context.Context, *ById) (*OrderItemGetByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedOrderItemServiceServer) GetAll(context.Context, *OrderItemGetAllReq) (*OrderItemGetAllRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedOrderItemServiceServer) mustEmbedUnimplementedOrderItemServiceServer() {}

// UnsafeOrderItemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderItemServiceServer will
// result in compilation errors.
type UnsafeOrderItemServiceServer interface {
	mustEmbedUnimplementedOrderItemServiceServer()
}

func RegisterOrderItemServiceServer(s grpc.ServiceRegistrar, srv OrderItemServiceServer) {
	s.RegisterService(&OrderItemService_ServiceDesc, srv)
}

func _OrderItemService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderItemCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).Create(ctx, req.(*OrderItemCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderItemService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).GetById(ctx, req.(*ById))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderItemService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderItemGetAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_GetAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).GetAll(ctx, req.(*OrderItemGetAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderItemService_ServiceDesc is the grpc.ServiceDesc for OrderItemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderItemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderItemService",
	HandlerType: (*OrderItemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _OrderItemService_Create_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _OrderItemService_GetById_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _OrderItemService_GetAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order_item.proto",
}
