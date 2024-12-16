package boilerplate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ BoilerplateServer = &boilerplate{}

type boilerplate struct {
	config              BoilerplateConfig
	tracer              trace.Tracer
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

func (s *boilerplate) Run(ctx context.Context) error {

	if s.config.Otel.Enabled {
		shutdown, err := setupOtel(ctx, s.config.Otel, s.config.ServiceName)
		if err != nil {
			return err
		}
		defer shutdown(ctx)
	}

	tp := otel.GetTracerProvider()
	s.tracer = tp.Tracer(s.config.TracerName)

	errChan := make(chan error)

	// if grpc is off, we can have no gateway either
	if s.config.Grpc.Disabled {
		return nil
	}

	go func() {
		errChan <- s.runGrpc()
	}()

	if !s.config.Gateway.Disabled {
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

	if s.config.Grpc.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(s.config.Grpc.TLS.Cert, s.config.Grpc.TLS.Key)
		if err != nil {
			return err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if s.config.Grpc.TLS.Mutual {
			ca := x509.NewCertPool()
			caBytes, err := os.ReadFile(s.config.Grpc.TLS.Ca)
			if err != nil {
				return err
			}
			if ok := ca.AppendCertsFromPEM(caBytes); !ok {
				return errors.New("could not load ca cert")
			}

			tlsConfig.ClientCAs = ca
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		creds := credentials.NewTLS(tlsConfig)
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

	if !s.config.Grpc.TLS.Enabled {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		ca := x509.NewCertPool()
		caBytes, err := os.ReadFile(s.config.Gateway.TLS.Ca)
		if err != nil {
			return err
		}
		if ok := ca.AppendCertsFromPEM(caBytes); !ok {
			return errors.New("could not load ca cert")
		}

		tlsConfig := &tls.Config{
			RootCAs:    ca,
			ServerName: "joe.mama",
		}
		if s.config.Grpc.TLS.Mutual {
			cert, err := tls.LoadX509KeyPair(s.config.Gateway.TLS.Cert, s.config.Gateway.TLS.Key)
			if err != nil {
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		creds := credentials.NewTLS(tlsConfig)
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	}

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

func (s *boilerplate) Tracer() trace.Tracer {
	return s.tracer
}
