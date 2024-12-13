package boilerplate

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ BoilerplateServer = &boilerplate{}

type boilerplate struct {
	config              BoilerplateConfig
	grpcRegisterFunc    GrpcRegisterFunc
	gatewayRegisterFunc GatewayRegisterFunc
}

func New() BoilerplateServer {
	return &boilerplate{}
}

func Default() BoilerplateServer {
	return &boilerplate{
		config: defaultConfig,
	}
}

func (s *boilerplate) WithConfig(conf BoilerplateConfig) {
	s.config = conf
}

func (s *boilerplate) WithGrpcPort(port uint) {
	s.config.Grpc.Port = port
}

func (s *boilerplate) WithGrpcHost(host string) {
	// TODO: validate host format
	s.config.Grpc.Host = host
}

func (s *boilerplate) Run(ctx context.Context) error {
	errChan := make(chan error)

	// if grpc is off, we can have no gateway either
	if !s.config.Grpc.Enabled {
		return nil
	}

	go func() {
		errChan <- s.runGrpc()
	}()

	if s.config.Gateway.Enabled {
		go func() {
			errChan <- s.runGateway(ctx)
		}()
	}

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return errors.New("context deadline exceeded")
	}
}

func (s *boilerplate) runGrpc() error {
	var opts []grpc.ServerOption

	if !s.config.Grpc.Insecure {
		creds, err := credentials.NewServerTLSFromFile(s.config.Grpc.TlsServerCert, s.config.Grpc.TlsServerKey)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	server := grpc.NewServer(opts...)
	err := s.grpcRegisterFunc(server)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", s.config.Grpc.Addr())
	if err != nil {
		return err
	}

	return server.Serve(lis)
}

func (s *boilerplate) runGateway(ctx context.Context) error {

	var dialOptions []grpc.DialOption

	if s.config.Grpc.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds, err := credentials.NewClientTLSFromFile(s.config.Grpc.TlsServerCaCert, "joe.mama")
		if err != nil {
			return err
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	}

	// TODO: set options from config
	//dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))

	conn, err := grpc.NewClient(
		s.config.Grpc.Addr(),
		dialOptions...,
	)

	if err != nil {
		return err
	}

	mux := runtime.NewServeMux()

	err = s.gatewayRegisterFunc(ctx, mux, conn)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    s.config.Gateway.Addr(),
		Handler: mux,
	}
	return server.ListenAndServe()
}
