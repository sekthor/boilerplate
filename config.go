package boilerplate

import "fmt"

const (
	DEFAULT_GRPC_PORT    = 50001
	DEFAULT_GATEWAY_PORT = 50002
	DEFAULT_HOST         = "0.0.0.0"
	DEFAULT_OTEL_PORT    = 4317
)

var defaultConfig = BoilerplateConfig{
	Grpc: ServerConfig{
		Port:            DEFAULT_GRPC_PORT,
		Host:            DEFAULT_HOST,
		Enabled:         true,
		Insecure:        false,
		TlsServerCert:   "server_cert.pem",
		TlsServerKey:    "server_key.pem",
		TlsServerCaCert: "ca_cert.pem",
	},
	Gateway: ServerConfig{
		Port:    DEFAULT_GATEWAY_PORT,
		Host:    DEFAULT_HOST,
		Enabled: true,
	},
	Otel: OtelConfig{
		OtelExporter: OtelExporter{
			Port: DEFAULT_OTEL_PORT,
			Host: "127.0.0.1",
		},
	},
}

type BoilerplateConfig struct {
	Grpc    ServerConfig
	Gateway ServerConfig
	Otel    OtelConfig
}

type ServerConfig struct {
	Enabled         bool
	Host            string
	Port            uint
	Insecure        bool
	TlsServerCert   string
	TlsServerKey    string
	TlsServerCaCert string
}

type OtelConfig struct {
	OtelExporter
	Tracing OtelExporter
	Metrics OtelExporter
}

type OtelExporter struct {
	Enabled  bool
	Host     string
	Port     uint
	Interval uint
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (e OtelExporter) Addr() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
