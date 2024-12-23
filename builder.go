package boilerplate

func (s *boilerplate) WithConfig(conf BoilerplateConfig) *boilerplate {
	s.config = conf
	return s
}

func (s *boilerplate) WithServiceName(name string) *boilerplate {
	s.config.ServiceName = name
	return s
}

func (s *boilerplate) WithGrpcAddr(addr string) *boilerplate {
	// TODO: validate host format
	s.config.Grpc.Addr = addr
	return s
}

func (s *boilerplate) WithGatewayAddr(addr string) *boilerplate {
	s.config.Gateway.Addr = addr
	return s
}

func (s *boilerplate) WithTracer(name string) *boilerplate {
	s.config.Otel.Enabled = true
	s.config.Otel.Tracing.Enabled = true
	s.config.Otel.TracerName = name
	return s
}

func (s *boilerplate) WithLogger(name string) *boilerplate {
	s.config.Otel.Enabled = true
	s.config.Otel.Logging.Enabled = true
	s.config.Otel.LoggerName = name
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

func (s *boilerplate) WithOtlpProtocol(protocol string) *boilerplate {
	s.config.Otel.Protocol = protocol
	return s
}

func (s *boilerplate) WithOtlpInsecure() *boilerplate {
	s.config.Otel.Insecure = true
	return s
}
