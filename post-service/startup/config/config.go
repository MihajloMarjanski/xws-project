package config

import "os"

//import "github.com/MihajloMarjanski/xws-project"

type Config struct {
	Port string
	Host string
}

func NewConfig() *Config {
	os.Setenv("POST_SERVICE_PORT", "8001")
	os.Setenv("POST_SERVICE_HOST", "post-service")
	return &Config{
		// err := godotenv.Load("dev.env")
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		Port: os.Getenv("POST_SERVICE_PORT"),
		Host: os.Getenv("POST_SERVICE_HOST"),
	}
}
