package main

import {
	
}

func main(){
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}