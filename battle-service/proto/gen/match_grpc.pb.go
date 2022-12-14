// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: match.proto

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

// MatchServiceClient is the client API for MatchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchServiceClient interface {
	Match(ctx context.Context, opts ...grpc.CallOption) (MatchService_MatchClient, error)
}

type matchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchServiceClient(cc grpc.ClientConnInterface) MatchServiceClient {
	return &matchServiceClient{cc}
}

func (c *matchServiceClient) Match(ctx context.Context, opts ...grpc.CallOption) (MatchService_MatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &MatchService_ServiceDesc.Streams[0], "/proto.MatchService/Match", opts...)
	if err != nil {
		return nil, err
	}
	x := &matchServiceMatchClient{stream}
	return x, nil
}

type MatchService_MatchClient interface {
	Send(*MatchRequest) error
	Recv() (*MatchResponse, error)
	grpc.ClientStream
}

type matchServiceMatchClient struct {
	grpc.ClientStream
}

func (x *matchServiceMatchClient) Send(m *MatchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *matchServiceMatchClient) Recv() (*MatchResponse, error) {
	m := new(MatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MatchServiceServer is the server API for MatchService service.
// All implementations should embed UnimplementedMatchServiceServer
// for forward compatibility
type MatchServiceServer interface {
	Match(MatchService_MatchServer) error
}

// UnimplementedMatchServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMatchServiceServer struct {
}

func (UnimplementedMatchServiceServer) Match(MatchService_MatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Match not implemented")
}

// UnsafeMatchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchServiceServer will
// result in compilation errors.
type UnsafeMatchServiceServer interface {
	mustEmbedUnimplementedMatchServiceServer()
}

func RegisterMatchServiceServer(s grpc.ServiceRegistrar, srv MatchServiceServer) {
	s.RegisterService(&MatchService_ServiceDesc, srv)
}

func _MatchService_Match_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MatchServiceServer).Match(&matchServiceMatchServer{stream})
}

type MatchService_MatchServer interface {
	Send(*MatchResponse) error
	Recv() (*MatchRequest, error)
	grpc.ServerStream
}

type matchServiceMatchServer struct {
	grpc.ServerStream
}

func (x *matchServiceMatchServer) Send(m *MatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *matchServiceMatchServer) Recv() (*MatchRequest, error) {
	m := new(MatchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MatchService_ServiceDesc is the grpc.ServiceDesc for MatchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MatchService",
	HandlerType: (*MatchServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Match",
			Handler:       _MatchService_Match_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "match.proto",
}
