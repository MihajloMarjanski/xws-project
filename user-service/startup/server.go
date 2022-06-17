package startup

import (
	"fmt"
	"google.golang.org/grpc/credentials"
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
		servicePath + "UpdateUser":        {"UpdateUser"},
		servicePath + "AddExperience":     {"AddExperience"},
		servicePath + "RemoveExperience":  {"RemoveExperience"},
		servicePath + "AddInterest":       {"AddInterest"},
		servicePath + "RemoveInterest":    {"RemoveInterest"},
		servicePath + "BlockUser":         {"BlockUser"},
		servicePath + "GetUserByUsername": {"GetUserByUsername"},
		//servicePath + "SearchOffers":      {"ROLE_USER"},
	}
}

func (server *Server) startGrpcServer(userHandler *handler_grpc.UserHandler) {
	log.Println(os.Hostname())
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptor := NewAuthInterceptor(accessibleRoles())

	creds, err := credentials.NewServerTLSFromFile("startup/certTLS/service.pem", "startup/certTLS/service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
		grpc.Creds(creds),
	)
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
