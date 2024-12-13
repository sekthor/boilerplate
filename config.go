package boilerplate

import "fmt"

const (
	DEFAULT_GRPC_PORT    = 50001
	DEFAULT_GATEWAY_PORT = 50002
	DEFAULT_HOST         = "0.0.0.0"
)

var defaultConfig = BoilerplateConfig{
	Grpc: Serverconfig{
		Port:    DEFAULT_GRPC_PORT,
		Host:    DEFAULT_HOST,
		Enabled: true,
	},
	Gateway: Serverconfig{
		Port:    DEFAULT_GATEWAY_PORT,
		Host:    DEFAULT_HOST,
		Enabled: true,
	},
}

type BoilerplateConfig struct {
	Grpc    Serverconfig
	Gateway Serverconfig
}

type Serverconfig struct {
	Enabled bool
	Host    string
	Port    uint
}

func (s Serverconfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
