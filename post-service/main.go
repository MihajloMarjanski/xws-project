package main

import (
	"post-service/startup"
	cfg "post-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
