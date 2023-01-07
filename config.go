package tinyurl

type ServerConfig struct {
	Host string
	Port string
}

func NewConfig() *ServerConfig {
	return &ServerConfig{}
}
