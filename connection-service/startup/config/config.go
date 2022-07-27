package config

import "os"

//import "github.com/MihajloMarjanski/xws-project"

type Config struct {
	Port string
	Host string
}

func NewConfig() *Config {
	os.Setenv("CONNECTION_SERVICE_PORT", "8700")
	os.Setenv("CONNECTION_SERVICE_HOST", "connection-service")
	return &Config{
		// err := godotenv.Load("dev.env")
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		Port: os.Getenv("CONNECTION_SERVICE_PORT"),
		Host: os.Getenv("CONNECTION_SERVICE_HOST"),
	}
}
