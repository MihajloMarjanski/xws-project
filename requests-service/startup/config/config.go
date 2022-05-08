package congig

import "os"

type Config struct {
	Port string
	Host string
}

func NewConfig() *Config {
	os.Setenv("REQUEST_SERVICE_PORT", "8200")
	os.Setenv("REQUEST_SERVICE_HOST", "requests-service")
	return &Config{
		Port: os.Getenv("REQUEST_SERVICE_PORT"),
		Host: os.Getenv("REQUEST_SERVICE_HOST"),
	}
}
