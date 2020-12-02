// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package circle

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// CircleServiceClient is the client API for CircleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CircleServiceClient interface {
	CircleTree(ctx context.Context, in *Circle, opts ...grpc.CallOption) (*CircleTreeResponse, error)
}

type circleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCircleServiceClient(cc grpc.ClientConnInterface) CircleServiceClient {
	return &circleServiceClient{cc}
}

func (c *circleServiceClient) CircleTree(ctx context.Context, in *Circle, opts ...grpc.CallOption) (*CircleTreeResponse, error) {
	out := new(CircleTreeResponse)
	err := c.cc.Invoke(ctx, "/circle.CircleService/CircleTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CircleServiceServer is the server API for CircleService service.
// All implementations must embed UnimplementedCircleServiceServer
// for forward compatibility
type CircleServiceServer interface {
	CircleTree(context.Context, *Circle) (*CircleTreeResponse, error)
	mustEmbedUnimplementedCircleServiceServer()
}

// UnimplementedCircleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCircleServiceServer struct {
}

func (UnimplementedCircleServiceServer) CircleTree(context.Context, *Circle) (*CircleTreeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CircleTree not implemented")
}
func (UnimplementedCircleServiceServer) mustEmbedUnimplementedCircleServiceServer() {}

// UnsafeCircleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CircleServiceServer will
// result in compilation errors.
type UnsafeCircleServiceServer interface {
	mustEmbedUnimplementedCircleServiceServer()
}

func RegisterCircleServiceServer(s grpc.ServiceRegistrar, srv CircleServiceServer) {
	s.RegisterService(&_CircleService_serviceDesc, srv)
}

func _CircleService_CircleTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Circle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CircleServiceServer).CircleTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circle.CircleService/CircleTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CircleServiceServer).CircleTree(ctx, req.(*Circle))
	}
	return interceptor(ctx, in, info, handler)
}

var _CircleService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "circle.CircleService",
	HandlerType: (*CircleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CircleTree",
			Handler:    _CircleService_CircleTree_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/manager/circle/circle.proto",
}