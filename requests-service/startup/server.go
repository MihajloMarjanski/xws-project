package startup

import (
	"fmt"
	"log"
	"net"
	"requests-service/handler_grpc"
	config "requests-service/startup/config"

	requestProto "github.com/MihajloMarjanski/xws-project/common/proto/request_service"

	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	requestProto.UnimplementedRequestServiceServer
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	requestHandler, err := handler_grpc.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	server.startGrpcServer(requestHandler)
}

func (server *Server) startGrpcServer(reqHandler *handler_grpc.RequestHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	requestProto.RegisterRequestServiceServer(grpcServer, reqHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
