package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sekthor/boilerplate"
	greeterv1 "github.com/sekthor/boilerplate/example/greeter/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ServiceImplementation struct {
	greeterv1.UnimplementedGreeterServiceServer
	server boilerplate.BoilerplateServer
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Roles []string
}

func customClaims() jwt.Claims { return &CustomClaims{} }

func (i *ServiceImplementation) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	t := i.server.Tracer()
	_, span := t.Start(ctx, "SayHello")
	defer span.End()

	name := req.GetName()

	claims, err := boilerplate.GetClaimsFromContext[*CustomClaims](ctx)
	if err == nil {
		name = name + " (" + claims.Subject + ")"
	}

	logrus.WithContext(ctx).Info("said hello")

	return &greeterv1.SayHelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}, nil
}

func main() {
	ctx := context.Background()

	i := &ServiceImplementation{}
	grpcFunc := func(s *grpc.Server) error {
		greeterv1.RegisterGreeterServiceServer(s, i)
		return nil
	}

	gatewayFunc := func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error {
		return greeterv1.RegisterGreeterServiceHandler(ctx, mux, cc)
	}
	authInterceptor, err := boilerplate.UnaryJwtClaimsInterceptor(
		[]string{"http://localhost:3001/realms/gig/protocol/openid-connect/certs"},
		customClaims,
	)

	if err != nil {
		log.Fatalf("could not create interceptor")
	}

	i.server = boilerplate.New().
		WithGrpcAddr(":50001").
		WithGatewayAddr(":50002").
		WithServiceName("BoilerPlate").
		WithTracer("github.com/sekthor/boilerplate/example/builder").
		WithOtlpProtocol("grpc").
		WithOtlpInsecure().
		WithLogger("github.com/sekthor/boilerplate/example/builder").
		WithGrpcRegisterFunc(grpcFunc).
		WithGatewayRegisterFunc(gatewayFunc).
		AddInterceptor(authInterceptor)

	if err := i.server.Run(ctx); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}
