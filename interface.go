package boilerplate

import "context"

type BoilerplateServer interface {
	WithConfig(BoilerplateConfig)
	WithGrpcHost(string)
	WithGrpcPort(uint)
	RegisterGateway(GatewayRegisterFunc)
	RegisterGrpc(GrpcRegisterFunc)
	Run(context.Context) error
}
