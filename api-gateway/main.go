package main

import (
	"api-gateway/startup"
	"api-gateway/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
