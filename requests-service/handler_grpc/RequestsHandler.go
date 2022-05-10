package handler_grpc

import (
	"context"
	"requests-service/service"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
)

type RequestsHandler struct {
	pb.UnimplementedRequestsServiceServer
	requestsService *service.RequestsService
}

func New() (*RequestsHandler, error) {

	requestsService, err := service.New()
	if err != nil {
		return nil, err
	}

	return &RequestsHandler{
		requestsService: requestsService,
	}, nil
}

func (requestsHandler *RequestsHandler) CloseDB() error {

	return requestsHandler.requestsService.CloseDB()
}

func (handler *RequestsHandler) GetAllByRecieverId(ctx context.Context, request *pb.GetAllByRecieverIdRequest) (*pb.GetAllByRecieverIdResponse, error) {
	id := request.RecieverID
	var requests []*pb.Request
	for _, request := range handler.requestsService.GetAllByRecieverId(uint(id)) {
		requests = append(requests, mapRequestToProto(request))
	}
	response := &pb.GetAllByRecieverIdResponse{
		Requests: requests,
	}
	return response, nil
}

func (handler *RequestsHandler) AcceptRequest(ctx context.Context, request *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
	senderID := request.SenderID
	recieverID := request.RecieverID
	handler.requestsService.AcceptRequest(uint(senderID), uint(recieverID))
	response := &pb.AcceptRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) DeclineRequest(ctx context.Context, request *pb.DeclineRequestRequest) (*pb.DeclineRequestResponse, error) {
	senderID := request.SenderID
	recieverID := request.RecieverID
	handler.requestsService.DeclineRequest(uint(senderID), uint(recieverID))
	response := &pb.DeclineRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) SendRequest(ctx context.Context, request *pb.SendRequestRequest) (*pb.SendRequestResponse, error) {
	senderID := request.SenderID
	recieverID := request.RecieverID
	handler.requestsService.SendRequest(uint(senderID), uint(recieverID))
	response := &pb.SendRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	senderID := request.SenderID
	receiverID := request.RecieverID
	message := request.Message.Text
	handler.requestsService.SendMessage(uint(senderID), uint(receiverID), message)
	response := &pb.SendMessageResponse{}
	return response, nil
}
