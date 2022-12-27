package server

type ServerConfig struct {
	Port string
}

func NewServerConfig(port string) ServerConfig {
	return ServerConfig{
		Port: port,
	}
}
