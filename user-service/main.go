package main

import (
	"user-service/startup"
	cfg "user-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
