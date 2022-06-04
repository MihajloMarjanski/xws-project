package startup

import (
	"fmt"
	"log"
	"net"
	"os"
	"requests-service/handler_grpc"
	"requests-service/service"
	config "requests-service/startup/config"

	requestProto "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
	"github.com/MihajloMarjanski/xws-project/common/saga/messaging/nats"

	"google.golang.org/grpc"
)

const (
	QueueGroup = "request_service"
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
	requestService, _ := service.New()
	commandSubscriber := server.initSubscriber("block.user.command", QueueGroup)
	replyPublisher := server.initPublisher("block.user.reply")
	server.initBlockUserHandler(requestService, replyPublisher, commandSubscriber)

	server.startGrpcServer(requestHandler)
}

func (server *Server) initBlockUserHandler(service *service.RequestsService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handler_grpc.NewBlockUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func accessibleRoles() map[string][]string {
	const servicePath = "/requests.RequestsService/"
	return map[string][]string{
		servicePath + "GetAllByRecieverId": {"ROLE_USER"},
		servicePath + "AcceptRequest":      {"ROLE_USER"},
		servicePath + "DeclineRequest":     {"ROLE_USER"},
		servicePath + "SendRequest":        {"ROLE_USER"},
		servicePath + "SendMessage":        {"ROLE_USER"},
		servicePath + "FindMessages":       {"ROLE_USER"},
		servicePath + "FindConnections":    {"ROLE_USER"},
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
