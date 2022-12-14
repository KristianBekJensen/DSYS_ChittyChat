// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/chat.proto

package proto

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

// ServicesClient is the client API for Services service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesClient interface {
	ChatService(ctx context.Context, opts ...grpc.CallOption) (Services_ChatServiceClient, error)
	ClientGreeting(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Id, error)
}

type servicesClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesClient(cc grpc.ClientConnInterface) ServicesClient {
	return &servicesClient{cc}
}

func (c *servicesClient) ChatService(ctx context.Context, opts ...grpc.CallOption) (Services_ChatServiceClient, error) {
	stream, err := c.cc.NewStream(ctx, &Services_ServiceDesc.Streams[0], "/proto.Services/ChatService", opts...)
	if err != nil {
		return nil, err
	}
	x := &servicesChatServiceClient{stream}
	return x, nil
}

type Services_ChatServiceClient interface {
	Send(*ClientMessage) error
	Recv() (*ServerMessage, error)
	grpc.ClientStream
}

type servicesChatServiceClient struct {
	grpc.ClientStream
}

func (x *servicesChatServiceClient) Send(m *ClientMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *servicesChatServiceClient) Recv() (*ServerMessage, error) {
	m := new(ServerMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *servicesClient) ClientGreeting(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Id, error) {
	out := new(Id)
	err := c.cc.Invoke(ctx, "/proto.Services/ClientGreeting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesServer is the server API for Services service.
// All implementations must embed UnimplementedServicesServer
// for forward compatibility
type ServicesServer interface {
	ChatService(Services_ChatServiceServer) error
	ClientGreeting(context.Context, *Id) (*Id, error)
	mustEmbedUnimplementedServicesServer()
}

// UnimplementedServicesServer must be embedded to have forward compatible implementations.
type UnimplementedServicesServer struct {
}

func (UnimplementedServicesServer) ChatService(Services_ChatServiceServer) error {
	return status.Errorf(codes.Unimplemented, "method ChatService not implemented")
}
func (UnimplementedServicesServer) ClientGreeting(context.Context, *Id) (*Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientGreeting not implemented")
}
func (UnimplementedServicesServer) mustEmbedUnimplementedServicesServer() {}

// UnsafeServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesServer will
// result in compilation errors.
type UnsafeServicesServer interface {
	mustEmbedUnimplementedServicesServer()
}

func RegisterServicesServer(s grpc.ServiceRegistrar, srv ServicesServer) {
	s.RegisterService(&Services_ServiceDesc, srv)
}

func _Services_ChatService_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServicesServer).ChatService(&servicesChatServiceServer{stream})
}

type Services_ChatServiceServer interface {
	Send(*ServerMessage) error
	Recv() (*ClientMessage, error)
	grpc.ServerStream
}

type servicesChatServiceServer struct {
	grpc.ServerStream
}

func (x *servicesChatServiceServer) Send(m *ServerMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *servicesChatServiceServer) Recv() (*ClientMessage, error) {
	m := new(ClientMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Services_ClientGreeting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).ClientGreeting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Services/ClientGreeting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).ClientGreeting(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

// Services_ServiceDesc is the grpc.ServiceDesc for Services service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Services_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Services",
	HandlerType: (*ServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClientGreeting",
			Handler:    _Services_ClientGreeting_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ChatService",
			Handler:       _Services_ChatService_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/chat.proto",
}
