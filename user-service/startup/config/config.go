package config

import "os"

//import "github.com/MihajloMarjanski/xws-project"

type Config struct {
	Port string
}

func NewConfig() *Config {
	os.Setenv("USER_SERVICE_PORT", "8000")
	return &Config{
		// err := godotenv.Load("dev.env")
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		Port: os.Getenv("USER_SERVICE_PORT"),
	}
}
