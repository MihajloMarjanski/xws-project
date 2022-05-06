package startup

import (
	"fmt"
	user "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-service/handler_grpc"
	"user-service/startup/config"
)

type Server struct {
	config *config.Config
	user.UnimplementedUserServiceServer
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	userHandler, err := handler_grpc.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	server.startGrpcServer(userHandler)
}

func (server *Server) startGrpcServer(userHandler *handler_grpc.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
