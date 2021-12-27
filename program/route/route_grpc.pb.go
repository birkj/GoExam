// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package helloworld

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

// RouteClient is the client API for Route service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouteClient interface {
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*Acknowledgement, error)
	SayHello(ctx context.Context, in *RequestText, opts ...grpc.CallOption) (*ReplyText, error)
	BroadcastMessage(ctx context.Context, in *RequestText, opts ...grpc.CallOption) (*GenericText, error)
}

type routeClient struct {
	cc grpc.ClientConnInterface
}

func NewRouteClient(cc grpc.ClientConnInterface) RouteClient {
	return &routeClient{cc}
}

func (c *routeClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*Acknowledgement, error) {
	out := new(Acknowledgement)
	err := c.cc.Invoke(ctx, "/Route/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routeClient) SayHello(ctx context.Context, in *RequestText, opts ...grpc.CallOption) (*ReplyText, error) {
	out := new(ReplyText)
	err := c.cc.Invoke(ctx, "/Route/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routeClient) BroadcastMessage(ctx context.Context, in *RequestText, opts ...grpc.CallOption) (*GenericText, error) {
	out := new(GenericText)
	err := c.cc.Invoke(ctx, "/Route/BroadcastMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RouteServer is the server API for Route service.
// All implementations must embed UnimplementedRouteServer
// for forward compatibility
type RouteServer interface {
	Connect(context.Context, *ConnectRequest) (*Acknowledgement, error)
	SayHello(context.Context, *RequestText) (*ReplyText, error)
	BroadcastMessage(context.Context, *RequestText) (*GenericText, error)
	mustEmbedUnimplementedRouteServer()
}

// UnimplementedRouteServer must be embedded to have forward compatible implementations.
type UnimplementedRouteServer struct {
}

func (UnimplementedRouteServer) Connect(context.Context, *ConnectRequest) (*Acknowledgement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedRouteServer) SayHello(context.Context, *RequestText) (*ReplyText, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedRouteServer) BroadcastMessage(context.Context, *RequestText) (*GenericText, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastMessage not implemented")
}
func (UnimplementedRouteServer) mustEmbedUnimplementedRouteServer() {}

// UnsafeRouteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouteServer will
// result in compilation errors.
type UnsafeRouteServer interface {
	mustEmbedUnimplementedRouteServer()
}

func RegisterRouteServer(s grpc.ServiceRegistrar, srv RouteServer) {
	s.RegisterService(&Route_ServiceDesc, srv)
}

func _Route_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Route/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Route_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestText)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Route/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteServer).SayHello(ctx, req.(*RequestText))
	}
	return interceptor(ctx, in, info, handler)
}

func _Route_BroadcastMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestText)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteServer).BroadcastMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Route/BroadcastMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteServer).BroadcastMessage(ctx, req.(*RequestText))
	}
	return interceptor(ctx, in, info, handler)
}

// Route_ServiceDesc is the grpc.ServiceDesc for Route service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Route_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Route",
	HandlerType: (*RouteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _Route_Connect_Handler,
		},
		{
			MethodName: "SayHello",
			Handler:    _Route_SayHello_Handler,
		},
		{
			MethodName: "BroadcastMessage",
			Handler:    _Route_BroadcastMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "route/route.proto",
}
