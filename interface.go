package boilerplate

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type BoilerplateServer interface {
	WithServiceName(string) *boilerplate
	WithConfig(BoilerplateConfig) *boilerplate
	WithGrpcHost(string) *boilerplate
	WithGrpcPort(uint) *boilerplate
	WithGatewayHost(string) *boilerplate
	WithGatewayPort(uint) *boilerplate
	WithGrpcRegisterFunc(GrpcRegisterFunc) *boilerplate
	WithGatewayRegisterFunc(GatewayRegisterFunc) *boilerplate
	WithTracer(string) *boilerplate
	WithJwks(jwksUrls []string) *boilerplate
	RegisterGateway(GatewayRegisterFunc)
	RegisterGrpc(GrpcRegisterFunc)
	Run(context.Context) error
	Tracer() trace.Tracer
}
