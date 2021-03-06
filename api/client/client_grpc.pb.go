// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.6.1
// source: client.proto

package client

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

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	Serve(ctx context.Context, opts ...grpc.CallOption) (Client_ServeClient, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) Serve(ctx context.Context, opts ...grpc.CallOption) (Client_ServeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Client_ServiceDesc.Streams[0], "/Client/Serve", opts...)
	if err != nil {
		return nil, err
	}
	x := &clientServeClient{stream}
	return x, nil
}

type Client_ServeClient interface {
	Send(*CToH) error
	Recv() (*HToC, error)
	grpc.ClientStream
}

type clientServeClient struct {
	grpc.ClientStream
}

func (x *clientServeClient) Send(m *CToH) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clientServeClient) Recv() (*HToC, error) {
	m := new(HToC)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClientServer is the server API for Client service.
// All implementations must embed UnimplementedClientServer
// for forward compatibility
type ClientServer interface {
	Serve(Client_ServeServer) error
	mustEmbedUnimplementedClientServer()
}

// UnimplementedClientServer must be embedded to have forward compatible implementations.
type UnimplementedClientServer struct {
}

func (UnimplementedClientServer) Serve(Client_ServeServer) error {
	return status.Errorf(codes.Unimplemented, "method Serve not implemented")
}
func (UnimplementedClientServer) mustEmbedUnimplementedClientServer() {}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_Serve_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClientServer).Serve(&clientServeServer{stream})
}

type Client_ServeServer interface {
	Send(*HToC) error
	Recv() (*CToH, error)
	grpc.ServerStream
}

type clientServeServer struct {
	grpc.ServerStream
}

func (x *clientServeServer) Send(m *HToC) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clientServeServer) Recv() (*CToH, error) {
	m := new(CToH)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Client",
	HandlerType: (*ClientServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Serve",
			Handler:       _Client_Serve_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "client.proto",
}
