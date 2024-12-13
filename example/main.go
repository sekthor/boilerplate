package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sekthor/boilerplate"
	greeterv1 "github.com/sekthor/boilerplate/example/greeter/v1"
	"google.golang.org/grpc"
)

type ServiceImplementation struct {
	greeterv1.UnimplementedGreeterServiceServer
	server boilerplate.BoilerplateServer
}

// implement protobuf service
func (i *ServiceImplementation) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	name := req.GetName()

	return &greeterv1.SayHelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}, nil
}

// register service implementation -> grpc Server
func (i *ServiceImplementation) GrpcFunc() boilerplate.GrpcRegisterFunc {
	return func(s *grpc.Server) error {
		greeterv1.RegisterGreeterServiceServer(s, i)
		return nil
	}
}

// register gateway handler
func (i *ServiceImplementation) GatewayFunc() boilerplate.GatewayRegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error {
		return greeterv1.RegisterGreeterServiceHandler(ctx, mux, cc)
	}
}

func main() {
	service := ServiceImplementation{}
	server := boilerplate.Default()
	server.RegisterGrpc(service.GrpcFunc())
	server.RegisterGateway(service.GatewayFunc())
	service.server = server

	ctx := context.Background()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
