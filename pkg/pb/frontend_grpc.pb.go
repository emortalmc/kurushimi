// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: frontend.proto

package pb

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

// FrontendClient is the client API for Frontend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontendClient interface {
	CreateTicket(ctx context.Context, in *CreateTicketRequest, opts ...grpc.CallOption) (*Ticket, error)
	DeleteTicket(ctx context.Context, in *DeleteTicketRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetTicket(ctx context.Context, in *GetTicketRequest, opts ...grpc.CallOption) (*Ticket, error)
	WatchTicketCountdown(ctx context.Context, in *WatchCountdownRequest, opts ...grpc.CallOption) (Frontend_WatchTicketCountdownClient, error)
	WatchTicketAssignment(ctx context.Context, in *WatchAssignmentRequest, opts ...grpc.CallOption) (Frontend_WatchTicketAssignmentClient, error)
}

type frontendClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontendClient(cc grpc.ClientConnInterface) FrontendClient {
	return &frontendClient{cc}
}

func (c *frontendClient) CreateTicket(ctx context.Context, in *CreateTicketRequest, opts ...grpc.CallOption) (*Ticket, error) {
	out := new(Ticket)
	err := c.cc.Invoke(ctx, "/dev.emortal.kurushimi.Frontend/CreateTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) DeleteTicket(ctx context.Context, in *DeleteTicketRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/dev.emortal.kurushimi.Frontend/DeleteTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) GetTicket(ctx context.Context, in *GetTicketRequest, opts ...grpc.CallOption) (*Ticket, error) {
	out := new(Ticket)
	err := c.cc.Invoke(ctx, "/dev.emortal.kurushimi.Frontend/GetTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) WatchTicketCountdown(ctx context.Context, in *WatchCountdownRequest, opts ...grpc.CallOption) (Frontend_WatchTicketCountdownClient, error) {
	stream, err := c.cc.NewStream(ctx, &Frontend_ServiceDesc.Streams[0], "/dev.emortal.kurushimi.Frontend/WatchTicketCountdown", opts...)
	if err != nil {
		return nil, err
	}
	x := &frontendWatchTicketCountdownClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Frontend_WatchTicketCountdownClient interface {
	Recv() (*WatchCountdownResponse, error)
	grpc.ClientStream
}

type frontendWatchTicketCountdownClient struct {
	grpc.ClientStream
}

func (x *frontendWatchTicketCountdownClient) Recv() (*WatchCountdownResponse, error) {
	m := new(WatchCountdownResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *frontendClient) WatchTicketAssignment(ctx context.Context, in *WatchAssignmentRequest, opts ...grpc.CallOption) (Frontend_WatchTicketAssignmentClient, error) {
	stream, err := c.cc.NewStream(ctx, &Frontend_ServiceDesc.Streams[1], "/dev.emortal.kurushimi.Frontend/WatchTicketAssignment", opts...)
	if err != nil {
		return nil, err
	}
	x := &frontendWatchTicketAssignmentClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Frontend_WatchTicketAssignmentClient interface {
	Recv() (*WatchAssignmentResponse, error)
	grpc.ClientStream
}

type frontendWatchTicketAssignmentClient struct {
	grpc.ClientStream
}

func (x *frontendWatchTicketAssignmentClient) Recv() (*WatchAssignmentResponse, error) {
	m := new(WatchAssignmentResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FrontendServer is the server API for Frontend service.
// All implementations must embed UnimplementedFrontendServer
// for forward compatibility
type FrontendServer interface {
	CreateTicket(context.Context, *CreateTicketRequest) (*Ticket, error)
	DeleteTicket(context.Context, *DeleteTicketRequest) (*emptypb.Empty, error)
	GetTicket(context.Context, *GetTicketRequest) (*Ticket, error)
	WatchTicketCountdown(*WatchCountdownRequest, Frontend_WatchTicketCountdownServer) error
	WatchTicketAssignment(*WatchAssignmentRequest, Frontend_WatchTicketAssignmentServer) error
	mustEmbedUnimplementedFrontendServer()
}

// UnimplementedFrontendServer must be embedded to have forward compatible implementations.
type UnimplementedFrontendServer struct {
}

func (UnimplementedFrontendServer) CreateTicket(context.Context, *CreateTicketRequest) (*Ticket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTicket not implemented")
}
func (UnimplementedFrontendServer) DeleteTicket(context.Context, *DeleteTicketRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTicket not implemented")
}
func (UnimplementedFrontendServer) GetTicket(context.Context, *GetTicketRequest) (*Ticket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTicket not implemented")
}
func (UnimplementedFrontendServer) WatchTicketCountdown(*WatchCountdownRequest, Frontend_WatchTicketCountdownServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchTicketCountdown not implemented")
}
func (UnimplementedFrontendServer) WatchTicketAssignment(*WatchAssignmentRequest, Frontend_WatchTicketAssignmentServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchTicketAssignment not implemented")
}
func (UnimplementedFrontendServer) mustEmbedUnimplementedFrontendServer() {}

// UnsafeFrontendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontendServer will
// result in compilation errors.
type UnsafeFrontendServer interface {
	mustEmbedUnimplementedFrontendServer()
}

func RegisterFrontendServer(s grpc.ServiceRegistrar, srv FrontendServer) {
	s.RegisterService(&Frontend_ServiceDesc, srv)
}

func _Frontend_CreateTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).CreateTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.emortal.kurushimi.Frontend/CreateTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).CreateTicket(ctx, req.(*CreateTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_DeleteTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).DeleteTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.emortal.kurushimi.Frontend/DeleteTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).DeleteTicket(ctx, req.(*DeleteTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_GetTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).GetTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.emortal.kurushimi.Frontend/GetTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).GetTicket(ctx, req.(*GetTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_WatchTicketCountdown_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchCountdownRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FrontendServer).WatchTicketCountdown(m, &frontendWatchTicketCountdownServer{stream})
}

type Frontend_WatchTicketCountdownServer interface {
	Send(*WatchCountdownResponse) error
	grpc.ServerStream
}

type frontendWatchTicketCountdownServer struct {
	grpc.ServerStream
}

func (x *frontendWatchTicketCountdownServer) Send(m *WatchCountdownResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Frontend_WatchTicketAssignment_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchAssignmentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FrontendServer).WatchTicketAssignment(m, &frontendWatchTicketAssignmentServer{stream})
}

type Frontend_WatchTicketAssignmentServer interface {
	Send(*WatchAssignmentResponse) error
	grpc.ServerStream
}

type frontendWatchTicketAssignmentServer struct {
	grpc.ServerStream
}

func (x *frontendWatchTicketAssignmentServer) Send(m *WatchAssignmentResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Frontend_ServiceDesc is the grpc.ServiceDesc for Frontend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Frontend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dev.emortal.kurushimi.Frontend",
	HandlerType: (*FrontendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTicket",
			Handler:    _Frontend_CreateTicket_Handler,
		},
		{
			MethodName: "DeleteTicket",
			Handler:    _Frontend_DeleteTicket_Handler,
		},
		{
			MethodName: "GetTicket",
			Handler:    _Frontend_GetTicket_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchTicketCountdown",
			Handler:       _Frontend_WatchTicketCountdown_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "WatchTicketAssignment",
			Handler:       _Frontend_WatchTicketAssignment_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "frontend.proto",
}
