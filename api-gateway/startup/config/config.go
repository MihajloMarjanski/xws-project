package config

import "os"

type Config struct {
	Port     string
	UserHost string
	UserPort string
}

func NewConfig() *Config {
	os.Setenv("GATEWAY_PORT", "8100")
	os.Setenv("USER_SERVICE_HOST", "DESKTOP-NJH4ABT")
	os.Setenv("USER_SERVICE_PORT", "8000")
	return &Config{
		// err := godotenv.Load("dev.env")
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		Port:     os.Getenv("GATEWAY_PORT"),
		UserHost: os.Getenv("USER_SERVICE_HOST"),
		UserPort: os.Getenv("USER_SERVICE_PORT"),
	}
}
