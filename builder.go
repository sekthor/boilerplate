package boilerplate

func (s *boilerplate) WithConfig(conf BoilerplateConfig) *boilerplate {
	s.config = conf
	return s
}

func (s *boilerplate) WithServiceName(name string) *boilerplate {
	s.config.ServiceName = name
	return s
}

func (s *boilerplate) WithGrpcPort(port uint) *boilerplate {
	s.config.Grpc.Port = port
	return s
}

func (s *boilerplate) WithGrpcHost(host string) *boilerplate {
	// TODO: validate host format
	s.config.Grpc.Host = host
	return s
}

func (s *boilerplate) WithGatewayPort(port uint) *boilerplate {
	s.config.Gateway.Port = port
	return s
}

func (s *boilerplate) WithGatewayHost(host string) *boilerplate {
	// TODO: validate host format
	s.config.Gateway.Host = host
	return s
}

func (s *boilerplate) WithTracer(name string) *boilerplate {
	s.config.Otel.Enabled = true
	s.config.Otel.Tracing.Enabled = true
	s.config.TracerName = name
	s.config.Otel.Tracing.Interval = 5
	return s
}

func (s *boilerplate) WithGrpcRegisterFunc(f GrpcRegisterFunc) *boilerplate {
	s.grpcRegisterFunc = f
	return s
}

func (s *boilerplate) WithGatewayRegisterFunc(f GatewayRegisterFunc) *boilerplate {
	s.gatewayRegisterFunc = f
	return s
}

func (s *boilerplate) WithJwks(jwksUrls []string) *boilerplate {
	s.config.JwkUrls = jwksUrls
	return s
}
