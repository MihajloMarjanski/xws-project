package startup

import (
	"connection-service/handler_grpc"
	"connection-service/startup/config"
	"fmt"
	"log"
	"net"
	"os"

	connection "github.com/MihajloMarjanski/xws-project/common/proto/connection_service"
	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	connection.UnimplementedConnectionServiceServer
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	connectionHandler, err := handler_grpc.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	server.startGrpcServer(connectionHandler)
}

func (server *Server) startGrpcServer(connectionHandler *handler_grpc.ConnectionHandler) {
	log.Println(os.Hostname())
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	grpcServer := grpc.NewServer()
	connection.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
