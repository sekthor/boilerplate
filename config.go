package boilerplate

import (
	"fmt"
	"time"
)

const (
	DEFAULT_GRPC_PORT     = 50001
	DEFAULT_GATEWAY_PORT  = 50002
	DEFAULT_HOST          = "0.0.0.0"
	DEFAULT_OTEL_PORT     = 4317
	DEFAULT_OTEL_INTERVAL = 5
)

var defaultConfig = BoilerplateConfig{
	ServiceName: "UnnamedBoilerplateService",
	TracerName:  "github.com/sekthor/boilerplate",
	Grpc: ServerConfig{
		Port: DEFAULT_GRPC_PORT,
		Host: DEFAULT_HOST,

		TLS: TlsConfig{
			Enabled: false,
			Mutual:  true,
			Cert:    "certs/server_cert.pem",
			Key:     "certs/server_key.pem",
			Ca:      "certs/client_ca_cert.pem",
		},
	},
	Gateway: ServerConfig{
		Port: DEFAULT_GATEWAY_PORT,
		Host: DEFAULT_HOST,

		TLS: TlsConfig{
			Mutual: true,
			Cert:   "certs/client_cert.pem",
			Key:    "certs/client_key.pem",
			Ca:     "certs/server_ca_cert.pem",
		},
	},
	Otel: OtelConfig{
		OtelExporterConfig: OtelExporterConfig{
			Enabled: true,
			Port:    DEFAULT_OTEL_PORT,
			Host:    "127.0.0.1",
		},
		Tracing: OtelExporterConfig{
			Enabled:  true,
			Protocol: "grpc",
			Port:     DEFAULT_OTEL_PORT,
			Host:     "127.0.0.1",
			Interval: 5,
			Insecure: true,
		},
	},
}

type BoilerplateConfig struct {
	ServiceName string
	TracerName  string
	Grpc        ServerConfig
	Gateway     ServerConfig
	Otel        OtelConfig
	JwkUrls     []string
}

type ServerConfig struct {
	Disabled bool
	Host     string
	Port     uint
	TLS      TlsConfig
}

type OtelConfig struct {
	OtelExporterConfig
	Tracing OtelExporterConfig
	Metrics OtelExporterConfig
}

type OtelExporterConfig struct {
	Enabled  bool
	Host     string
	Port     uint
	Interval uint
	Protocol string
	Insecure bool
}

type TlsConfig struct {
	Enabled bool
	Mutual  bool
	Key     string
	Cert    string
	Ca      string
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (e OtelExporterConfig) Addr() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (c OtelConfig) TracingPort() uint {
	if c.Tracing.Port != 0 {
		return c.Tracing.Port
	}

	if c.Port != 0 {
		return c.Port
	}

	return DEFAULT_OTEL_PORT
}

func (c OtelConfig) MetricsPort() uint {
	if c.Metrics.Port != 0 {
		return c.Metrics.Port
	}

	if c.Port != 0 {
		return c.Port
	}

	return DEFAULT_OTEL_PORT
}

func (c OtelConfig) TracingHost() string {
	if c.Tracing.Host != "" {
		return c.Tracing.Host
	}

	return c.Host
}

func (c OtelConfig) MetricsHost() string {
	if c.Metrics.Host != "" {
		return c.Metrics.Host
	}

	return c.Host
}

func (c OtelConfig) TracingInterval() time.Duration {
	if c.Tracing.Interval != 0 {
		return time.Second * time.Duration(c.Tracing.Interval)
	}

	if c.Interval != 0 {
		return time.Second * time.Duration(c.Interval)
	}
	return DEFAULT_OTEL_INTERVAL * time.Second
}

func (c OtelConfig) MetricsInterval() time.Duration {
	if c.Metrics.Interval != 0 {
		return time.Second * time.Duration(c.Metrics.Interval)
	}

	if c.Interval != 0 {
		return time.Second * time.Duration(c.Interval)
	}
	return DEFAULT_OTEL_INTERVAL * time.Second
}
