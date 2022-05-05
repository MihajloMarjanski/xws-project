package handler_grpc

import (
	"context"
	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"mime"
	"net/http"
	"strconv"
	"user-service/model"
	"user-service/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

func New() (*UserHandler, error) {

	userService, err := service.New()
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		userService: userService,
	}, nil
}

func (handler *UserHandler) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := request.Id
	user := handler.userService.GetByID(int(id))
	userPb := mapUser(user)
	response := &pb.GetUserResponse{
		User: userPb,
	}
	return response, nil
}