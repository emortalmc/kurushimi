// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: grpc.proto

package pb

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

// MatchmakerClient is the client API for Matchmaker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchmakerClient interface {
	QueueByPlayer(ctx context.Context, in *QueueByPlayerRequest, opts ...grpc.CallOption) (*QueueByPlayerResponse, error)
	// QueueInitialLobbyByPlayer is used to queue a player for a lobby match before they've technically joined the server.
	// This is used by the proxy as the player doesn't yet have a party.
	// Note this method does much less validation and error handling than QueueByPlayer due to its intended use case.
	// This method will create a ticket with auto_teleport as false and default to the 'lobby' game mode.
	QueueInitialLobbyByPlayer(ctx context.Context, in *QueueInitialLobbyByPlayerRequest, opts ...grpc.CallOption) (*QueueInitialLobbyByPlayerResponse, error)
	DequeueByPlayer(ctx context.Context, in *DequeueByPlayerRequest, opts ...grpc.CallOption) (*DequeueByPlayerResponse, error)
	ChangePlayerMapVote(ctx context.Context, in *ChangePlayerMapVoteRequest, opts ...grpc.CallOption) (*ChangePlayerMapVoteResponse, error)
	GetPlayerQueueInfo(ctx context.Context, in *GetPlayerQueueInfoRequest, opts ...grpc.CallOption) (*GetPlayerQueueInfoResponse, error)
}

type matchmakerClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchmakerClient(cc grpc.ClientConnInterface) MatchmakerClient {
	return &matchmakerClient{cc}
}

func (c *matchmakerClient) QueueByPlayer(ctx context.Context, in *QueueByPlayerRequest, opts ...grpc.CallOption) (*QueueByPlayerResponse, error) {
	out := new(QueueByPlayerResponse)
	err := c.cc.Invoke(ctx, "/emortal.kurushimi.grpc.Matchmaker/QueueByPlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakerClient) QueueInitialLobbyByPlayer(ctx context.Context, in *QueueInitialLobbyByPlayerRequest, opts ...grpc.CallOption) (*QueueInitialLobbyByPlayerResponse, error) {
	out := new(QueueInitialLobbyByPlayerResponse)
	err := c.cc.Invoke(ctx, "/emortal.kurushimi.grpc.Matchmaker/QueueInitialLobbyByPlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakerClient) DequeueByPlayer(ctx context.Context, in *DequeueByPlayerRequest, opts ...grpc.CallOption) (*DequeueByPlayerResponse, error) {
	out := new(DequeueByPlayerResponse)
	err := c.cc.Invoke(ctx, "/emortal.kurushimi.grpc.Matchmaker/DequeueByPlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakerClient) ChangePlayerMapVote(ctx context.Context, in *ChangePlayerMapVoteRequest, opts ...grpc.CallOption) (*ChangePlayerMapVoteResponse, error) {
	out := new(ChangePlayerMapVoteResponse)
	err := c.cc.Invoke(ctx, "/emortal.kurushimi.grpc.Matchmaker/ChangePlayerMapVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakerClient) GetPlayerQueueInfo(ctx context.Context, in *GetPlayerQueueInfoRequest, opts ...grpc.CallOption) (*GetPlayerQueueInfoResponse, error) {
	out := new(GetPlayerQueueInfoResponse)
	err := c.cc.Invoke(ctx, "/emortal.kurushimi.grpc.Matchmaker/GetPlayerQueueInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchmakerServer is the server API for Matchmaker service.
// All implementations must embed UnimplementedMatchmakerServer
// for forward compatibility
type MatchmakerServer interface {
	QueueByPlayer(context.Context, *QueueByPlayerRequest) (*QueueByPlayerResponse, error)
	// QueueInitialLobbyByPlayer is used to queue a player for a lobby match before they've technically joined the server.
	// This is used by the proxy as the player doesn't yet have a party.
	// Note this method does much less validation and error handling than QueueByPlayer due to its intended use case.
	// This method will create a ticket with auto_teleport as false and default to the 'lobby' game mode.
	QueueInitialLobbyByPlayer(context.Context, *QueueInitialLobbyByPlayerRequest) (*QueueInitialLobbyByPlayerResponse, error)
	DequeueByPlayer(context.Context, *DequeueByPlayerRequest) (*DequeueByPlayerResponse, error)
	ChangePlayerMapVote(context.Context, *ChangePlayerMapVoteRequest) (*ChangePlayerMapVoteResponse, error)
	GetPlayerQueueInfo(context.Context, *GetPlayerQueueInfoRequest) (*GetPlayerQueueInfoResponse, error)
	mustEmbedUnimplementedMatchmakerServer()
}

// UnimplementedMatchmakerServer must be embedded to have forward compatible implementations.
type UnimplementedMatchmakerServer struct {
}

func (UnimplementedMatchmakerServer) QueueByPlayer(context.Context, *QueueByPlayerRequest) (*QueueByPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueueByPlayer not implemented")
}
func (UnimplementedMatchmakerServer) QueueInitialLobbyByPlayer(context.Context, *QueueInitialLobbyByPlayerRequest) (*QueueInitialLobbyByPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueueInitialLobbyByPlayer not implemented")
}
func (UnimplementedMatchmakerServer) DequeueByPlayer(context.Context, *DequeueByPlayerRequest) (*DequeueByPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DequeueByPlayer not implemented")
}
func (UnimplementedMatchmakerServer) ChangePlayerMapVote(context.Context, *ChangePlayerMapVoteRequest) (*ChangePlayerMapVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePlayerMapVote not implemented")
}
func (UnimplementedMatchmakerServer) GetPlayerQueueInfo(context.Context, *GetPlayerQueueInfoRequest) (*GetPlayerQueueInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayerQueueInfo not implemented")
}
func (UnimplementedMatchmakerServer) mustEmbedUnimplementedMatchmakerServer() {}

// UnsafeMatchmakerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchmakerServer will
// result in compilation errors.
type UnsafeMatchmakerServer interface {
	mustEmbedUnimplementedMatchmakerServer()
}

func RegisterMatchmakerServer(s grpc.ServiceRegistrar, srv MatchmakerServer) {
	s.RegisterService(&Matchmaker_ServiceDesc, srv)
}

func _Matchmaker_QueueByPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueByPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakerServer).QueueByPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emortal.kurushimi.grpc.Matchmaker/QueueByPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakerServer).QueueByPlayer(ctx, req.(*QueueByPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Matchmaker_QueueInitialLobbyByPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueInitialLobbyByPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakerServer).QueueInitialLobbyByPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emortal.kurushimi.grpc.Matchmaker/QueueInitialLobbyByPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakerServer).QueueInitialLobbyByPlayer(ctx, req.(*QueueInitialLobbyByPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Matchmaker_DequeueByPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DequeueByPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakerServer).DequeueByPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emortal.kurushimi.grpc.Matchmaker/DequeueByPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakerServer).DequeueByPlayer(ctx, req.(*DequeueByPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Matchmaker_ChangePlayerMapVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePlayerMapVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakerServer).ChangePlayerMapVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emortal.kurushimi.grpc.Matchmaker/ChangePlayerMapVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakerServer).ChangePlayerMapVote(ctx, req.(*ChangePlayerMapVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Matchmaker_GetPlayerQueueInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayerQueueInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakerServer).GetPlayerQueueInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emortal.kurushimi.grpc.Matchmaker/GetPlayerQueueInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakerServer).GetPlayerQueueInfo(ctx, req.(*GetPlayerQueueInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Matchmaker_ServiceDesc is the grpc.ServiceDesc for Matchmaker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Matchmaker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "emortal.kurushimi.grpc.Matchmaker",
	HandlerType: (*MatchmakerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueueByPlayer",
			Handler:    _Matchmaker_QueueByPlayer_Handler,
		},
		{
			MethodName: "QueueInitialLobbyByPlayer",
			Handler:    _Matchmaker_QueueInitialLobbyByPlayer_Handler,
		},
		{
			MethodName: "DequeueByPlayer",
			Handler:    _Matchmaker_DequeueByPlayer_Handler,
		},
		{
			MethodName: "ChangePlayerMapVote",
			Handler:    _Matchmaker_ChangePlayerMapVote_Handler,
		},
		{
			MethodName: "GetPlayerQueueInfo",
			Handler:    _Matchmaker_GetPlayerQueueInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}
