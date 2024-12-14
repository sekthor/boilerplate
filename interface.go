package boilerplate

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type BoilerplateServer interface {
	WithConfig(BoilerplateConfig) *boilerplate
	WithGrpcHost(string) *boilerplate
	WithGrpcPort(uint) *boilerplate
	WithTracer(string) *boilerplate
	RegisterGateway(GatewayRegisterFunc)
	RegisterGrpc(GrpcRegisterFunc)
	Run(context.Context) error
	Tracer() trace.Tracer
}
