package config

import "os"

type Config struct {
	Port        string
	UserHost    string
	UserPort    string
	RequestHost string
	RequestPort string
	PostHost    string
	PostPort    string
}

func NewConfig() *Config {
	os.Setenv("GATEWAY_PORT", "8000")
	os.Setenv("USER_SERVICE_HOST", "user-service")
	os.Setenv("USER_SERVICE_PORT", "8100")
	os.Setenv("REQUEST_SERVICE_HOST", "request-service")
	os.Setenv("REQUEST_SERVICE_PORT", "8200")
	os.Setenv("POST_SERVICE_HOST", "post-service")
	os.Setenv("POST_SERVICE_PORT", "8300")
	return &Config{
		// err := godotenv.Load("dev.env")
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		Port:        os.Getenv("GATEWAY_PORT"),
		UserHost:    os.Getenv("USER_SERVICE_HOST"),
		UserPort:    os.Getenv("USER_SERVICE_PORT"),
		RequestHost: os.Getenv("REQUEST_SERVICE_HOST"),
		RequestPort: os.Getenv("REQUEST_SERVICE_PORT"),
		PostHost:    os.Getenv("POST_SERVICE_HOST"),
		PostPort:    os.Getenv("POST_SERVICE_PORT"),
	}
}
