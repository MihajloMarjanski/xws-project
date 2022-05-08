package main

import (
	"requests-service/startup"
	cfg "requests-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
