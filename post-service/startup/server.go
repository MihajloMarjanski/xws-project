package startup

import (
	"fmt"
	"log"
	"net"
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

func accessibleRoles() map[string][]string {
	const servicePath = "/post.PostService/"
	return map[string][]string{
		servicePath + "CreatePost": {"user"},
		servicePath + "AddComment": {"user"},
		servicePath + "AddLike":    {"user"},
		servicePath + "AddDislike": {"user"},
	}
}

func (server *Server) startGrpcServer(postHandler *handler_grpc.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptor := NewAuthInterceptor(accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	post.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
