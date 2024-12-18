package boilerplate

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type BoilerplateServer interface {
	WithServiceName(string) *boilerplate
	WithConfig(BoilerplateConfig) *boilerplate
	WithGrpcAddr(string) *boilerplate
	WithGatewayAddr(string) *boilerplate
	WithGrpcRegisterFunc(GrpcRegisterFunc) *boilerplate
	WithGatewayRegisterFunc(GatewayRegisterFunc) *boilerplate
	WithTracer(string) *boilerplate
	WithJwks(jwksUrls []string) *boilerplate
	RegisterGateway(GatewayRegisterFunc)
	RegisterGrpc(GrpcRegisterFunc)
	Run(context.Context) error
	Tracer() trace.Tracer
}
