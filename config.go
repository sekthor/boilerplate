package boilerplate

import (
	"time"
)

const (
	DEFAULT_GRPC_ADDR     = ":50001"
	DEFAULT_GATEWAY_ADDR  = ":50002"
	DEFAULT_OTEL_ADDR     = "127.0.0.1:4317"
	DEFAULT_OTEL_INTERVAL = 5
)

var defaultConfig = BoilerplateConfig{
	ServiceName: "UnnamedBoilerplateService",
	Grpc: ServerConfig{
		Addr: DEFAULT_GRPC_ADDR,
		TLS: TlsConfig{
			Enabled: false,
			Mutual:  true,
			Cert:    "certs/server_cert.pem",
			Key:     "certs/server_key.pem",
			Ca:      "certs/client_ca_cert.pem",
		},
	},
	Gateway: ServerConfig{
		Addr: DEFAULT_GATEWAY_ADDR,
		TLS: TlsConfig{
			Mutual: true,
			Cert:   "certs/client_cert.pem",
			Key:    "certs/client_key.pem",
			Ca:     "certs/server_ca_cert.pem",
		},
	},
	Otel: OtelConfig{
		TracerName: "github.com/sekthor/boilerplate",
		LoggerName: "github.com/sekthor/boilerplate",
		OtelExporterConfig: OtelExporterConfig{
			Enabled:  true,
			Addr:     DEFAULT_OTEL_ADDR,
			Protocol: "grpc",
			Interval: 5,
		},
		Logging: OtelExporterConfig{
			Enabled: true,
		},
	},
}

type BoilerplateConfig struct {
	ServiceName string
	Grpc        ServerConfig
	Gateway     ServerConfig
	Otel        OtelConfig
	JwkUrls     []string
}

type ServerConfig struct {
	Disabled bool
	Addr     string
	TLS      TlsConfig
}

type OtelConfig struct {
	OtelExporterConfig
	Tracing    OtelExporterConfig
	Metrics    OtelExporterConfig
	Logging    OtelExporterConfig
	TracerName string
	LoggerName string
}

type OtelExporterConfig struct {
	Enabled  bool
	Addr     string
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

func (c OtelConfig) TracingAddr() string {
	if c.Tracing.Addr != "" {
		return c.Tracing.Addr
	}
	if c.Addr != "" {
		return c.Addr
	}
	return DEFAULT_OTEL_ADDR
}

func (c OtelConfig) MetricsAddr() string {
	if c.Metrics.Addr != "" {
		return c.Metrics.Addr
	}
	if c.Addr != "" {
		return c.Addr
	}
	return DEFAULT_OTEL_ADDR
}

func (c OtelConfig) LoggingAddr() string {
	if c.Logging.Addr != "" {
		return c.Logging.Addr
	}
	if c.Addr != "" {
		return c.Addr
	}
	return DEFAULT_OTEL_ADDR
}

func (c OtelConfig) TracingProtocol() string {
	if c.Tracing.Protocol != "" {
		return c.Tracing.Protocol
	}
	return c.Protocol
}

func (c OtelConfig) MetricsProtocol() string {
	if c.Metrics.Protocol != "" {
		return c.Metrics.Protocol
	}
	return c.Protocol
}

func (c OtelConfig) LoggingProtocol() string {
	if c.Logging.Protocol != "" {
		return c.Logging.Protocol
	}
	return c.Protocol
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

func (c OtelConfig) LoggingInterval() time.Duration {
	if c.Logging.Interval != 0 {
		return time.Second * time.Duration(c.Logging.Interval)
	}

	if c.Interval != 0 {
		return time.Second * time.Duration(c.Interval)
	}
	return DEFAULT_OTEL_INTERVAL * time.Second
}

func (c OtelConfig) TracingInsecure() bool {
	return c.Insecure || c.Tracing.Insecure
}

func (c OtelConfig) MetricsInsecure() bool {
	return c.Insecure || c.Metrics.Insecure
}

func (c OtelConfig) LoggingInsecure() bool {
	return c.Insecure || c.Logging.Insecure
}
