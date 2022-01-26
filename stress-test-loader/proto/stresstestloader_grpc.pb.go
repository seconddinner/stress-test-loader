// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package stresstestloader

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

// VerifierClient is the client API for Verifier service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VerifierClient interface {
	// Sends a list of proxy public ip to varify
	Verify(ctx context.Context, in *MonitorServerConfig, opts ...grpc.CallOption) (*VerifyReply, error)
}

type verifierClient struct {
	cc grpc.ClientConnInterface
}

func NewVerifierClient(cc grpc.ClientConnInterface) VerifierClient {
	return &verifierClient{cc}
}

func (c *verifierClient) Verify(ctx context.Context, in *MonitorServerConfig, opts ...grpc.CallOption) (*VerifyReply, error) {
	out := new(VerifyReply)
	err := c.cc.Invoke(ctx, "/stresstestloader.Verifier/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerifierServer is the server API for Verifier service.
// All implementations must embed UnimplementedVerifierServer
// for forward compatibility
type VerifierServer interface {
	// Sends a list of proxy public ip to varify
	Verify(context.Context, *MonitorServerConfig) (*VerifyReply, error)
	mustEmbedUnimplementedVerifierServer()
}

// UnimplementedVerifierServer must be embedded to have forward compatible implementations.
type UnimplementedVerifierServer struct {
}

func (UnimplementedVerifierServer) Verify(context.Context, *MonitorServerConfig) (*VerifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedVerifierServer) mustEmbedUnimplementedVerifierServer() {}

// UnsafeVerifierServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VerifierServer will
// result in compilation errors.
type UnsafeVerifierServer interface {
	mustEmbedUnimplementedVerifierServer()
}

func RegisterVerifierServer(s grpc.ServiceRegistrar, srv VerifierServer) {
	s.RegisterService(&Verifier_ServiceDesc, srv)
}

func _Verifier_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorServerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerifierServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stresstestloader.Verifier/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerifierServer).Verify(ctx, req.(*MonitorServerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

// Verifier_ServiceDesc is the grpc.ServiceDesc for Verifier service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Verifier_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stresstestloader.Verifier",
	HandlerType: (*VerifierServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Verify",
			Handler:    _Verifier_Verify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/stresstestloader.proto",
}
