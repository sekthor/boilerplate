package boilerplate

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GrpcRegisterFunc func(*grpc.Server) error
type GatewayRegisterFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

func (s *boilerplate) RegisterGrpc(grpcRegisterFunc GrpcRegisterFunc) {
	s.grpcRegisterFunc = grpcRegisterFunc
}

func (s *boilerplate) RegisterGateway(gatewayRegisterFunc GatewayRegisterFunc) {
	s.gatewayRegisterFunc = gatewayRegisterFunc
}
