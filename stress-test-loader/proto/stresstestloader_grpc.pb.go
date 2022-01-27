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

// LoadTestLoaderClient is the client API for LoadTestLoader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoadTestLoaderClient interface {
	// Sends a list of stresstest public ip to varify
	StartLoadTest(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestReply, error)
}

type loadTestLoaderClient struct {
	cc grpc.ClientConnInterface
}

func NewLoadTestLoaderClient(cc grpc.ClientConnInterface) LoadTestLoaderClient {
	return &loadTestLoaderClient{cc}
}

func (c *loadTestLoaderClient) StartLoadTest(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestReply, error) {
	out := new(TestReply)
	err := c.cc.Invoke(ctx, "/stresstestloader.LoadTestLoader/StartLoadTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoadTestLoaderServer is the server API for LoadTestLoader service.
// All implementations must embed UnimplementedLoadTestLoaderServer
// for forward compatibility
type LoadTestLoaderServer interface {
	// Sends a list of stresstest public ip to varify
	StartLoadTest(context.Context, *TestRequest) (*TestReply, error)
	mustEmbedUnimplementedLoadTestLoaderServer()
}

// UnimplementedLoadTestLoaderServer must be embedded to have forward compatible implementations.
type UnimplementedLoadTestLoaderServer struct {
}

func (UnimplementedLoadTestLoaderServer) StartLoadTest(context.Context, *TestRequest) (*TestReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartLoadTest not implemented")
}
func (UnimplementedLoadTestLoaderServer) mustEmbedUnimplementedLoadTestLoaderServer() {}

// UnsafeLoadTestLoaderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoadTestLoaderServer will
// result in compilation errors.
type UnsafeLoadTestLoaderServer interface {
	mustEmbedUnimplementedLoadTestLoaderServer()
}

func RegisterLoadTestLoaderServer(s grpc.ServiceRegistrar, srv LoadTestLoaderServer) {
	s.RegisterService(&LoadTestLoader_ServiceDesc, srv)
}

func _LoadTestLoader_StartLoadTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadTestLoaderServer).StartLoadTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stresstestloader.LoadTestLoader/StartLoadTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadTestLoaderServer).StartLoadTest(ctx, req.(*TestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LoadTestLoader_ServiceDesc is the grpc.ServiceDesc for LoadTestLoader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoadTestLoader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stresstestloader.LoadTestLoader",
	HandlerType: (*LoadTestLoaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartLoadTest",
			Handler:    _LoadTestLoader_StartLoadTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/stresstestloader.proto",
}
