package startup

import (
	"fmt"
	"log"
	"net"
	"os"
	"post-service/handler_grpc"
	"post-service/startup/config"

	post "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	post.UnimplementedPostServiceServer
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	postHandler, err := handler_grpc.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	server.startGrpcServer(postHandler)
}

func (server *Server) startGrpcServer(postHandler *handler_grpc.PostHandler) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", name)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	post.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
