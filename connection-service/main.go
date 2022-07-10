package main

import (
	"connection-service/startup"
	cfg "connection-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
