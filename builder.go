package boilerplate

func (s *boilerplate) WithConfig(conf BoilerplateConfig) *boilerplate {
	s.config = conf
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

func (s *boilerplate) WithTracer(name string) *boilerplate {
	s.config.TracerName = name
	return s
}
