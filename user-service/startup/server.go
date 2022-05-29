package startup

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-service/handler_grpc"
	"user-service/startup/config"

	user "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"google.golang.org/grpc"
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

func accessibleRoles() map[string][]string {
	const servicePath = "/user.UserService/"
	return map[string][]string{
		servicePath + "UpdateUser":       {"user"},
		servicePath + "AddExperience":    {"user"},
		servicePath + "RemoveExperience": {"user"},
		servicePath + "AddInterest":      {"user"},
		servicePath + "RemoveInterest":   {"user"},
		servicePath + "BlockUser":        {"user"},
	}
}

func (server *Server) startGrpcServer(userHandler *handler_grpc.UserHandler) {
	log.Println(os.Hostname())
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptor := NewAuthInterceptor(accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
