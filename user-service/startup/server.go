package startup

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-service/handler_grpc"
	"user-service/service"
	"user-service/startup/config"

	user "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
	"github.com/MihajloMarjanski/xws-project/common/saga/messaging/nats"
	"google.golang.org/grpc"
)

const (
	QueueGroup = "user_service"
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

	requestService, _ := service.New()
	commandSubscriber := server.initSubscriber("block.user.command", QueueGroup)
	replyPublisher := server.initPublisher("block.user.reply")
	server.initBlockUserHandler(requestService, replyPublisher, commandSubscriber)

	server.startGrpcServer(userHandler)
}

func (server *Server) initBlockUserHandler(service *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
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
	const servicePath = "/user.UserService/"
	return map[string][]string{
		servicePath + "UpdateUser":        {"ROLE_USER"},
		servicePath + "AddExperience":     {"ROLE_USER"},
		servicePath + "RemoveExperience":  {"ROLE_USER"},
		servicePath + "AddInterest":       {"ROLE_USER"},
		servicePath + "RemoveInterest":    {"ROLE_USER"},
		servicePath + "BlockUser":         {"ROLE_USER"},
		servicePath + "GetUserByUsername": {"ROLE_USER"},
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
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
