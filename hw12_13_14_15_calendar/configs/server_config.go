package configs

type ServerConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	PortGRPC string `toml:"port_grpc"`
}
