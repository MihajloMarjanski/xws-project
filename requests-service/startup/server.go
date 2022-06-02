package startup

import (
	"fmt"
	"log"
	"net"
	"os"
	"requests-service/handler_grpc"
	config "requests-service/startup/config"

	requestProto "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"

	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	requestProto.UnimplementedRequestsServiceServer
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

func accessibleRoles() map[string][]string {
	const servicePath = "/requests.RequestsService/"
	return map[string][]string{
		servicePath + "GetAllByRecieverId": {"ROLE_USER"},
		servicePath + "AcceptRequest":      {"ROLE_USER"},
		servicePath + "DeclineRequest":     {"ROLE_USER"},
		servicePath + "SendRequest":        {"ROLE_USER"},
		servicePath + "SendMessage":        {"ROLE_USER"},
	}
}

func (server *Server) startGrpcServer(reqHandler *handler_grpc.RequestsHandler) {
	log.Println(os.Hostname())
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	interceptor := NewAuthInterceptor(accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	requestProto.RegisterRequestsServiceServer(grpcServer, reqHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
