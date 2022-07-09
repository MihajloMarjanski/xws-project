package handler_grpc

import (
	"connection-service/service"
	"context"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/connection_service"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	connectionService *service.ConnectionService
}

func New() (*ConnectionHandler, error) {

	connectionService, _ := service.New()
	return &ConnectionHandler{
		connectionService: connectionService,
	}, nil
}

func (handler *ConnectionHandler) Connect(ctx context.Context, request *pb.UsersConnectionRequest) (*pb.UsersConnectionResponse, error) {
	id1 := request.UserPair.Id1
	id2 := request.UserPair.Id2
	handler.connectionService.Connect(id1, id2)

	response := &pb.UsersConnectionResponse{}
	return response, nil
}
