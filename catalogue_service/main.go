package main

import (
	"github.com/tamararankovic/microservices_demo/catalogue_service/startup"
	cfg "github.com/tamararankovic/microservices_demo/catalogue_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
