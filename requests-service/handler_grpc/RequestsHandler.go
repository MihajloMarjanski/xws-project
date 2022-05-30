package handler_grpc

import (
	"context"
	"fmt"
	"requests-service/model"
	"requests-service/service"
	"strconv"
	"strings"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

type RequestsHandler struct {
	pb.UnimplementedRequestsServiceServer
	requestsService *service.RequestsService
}

func Verify(accessToken string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&model.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte("tajni_kljuc_za_jwt_hash"), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GetUserID(ctx context.Context) uint {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md["authorization"]
	accessToken := values[0]
	words := strings.Fields(accessToken)

	claims, _ := Verify(words[1])
	id, _ := strconv.ParseUint(claims.Id, 10, 64)
	return uint(id)
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
	id = int64(GetUserID(ctx))
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
	senderID = int64(GetUserID(ctx))
	handler.requestsService.AcceptRequest(uint(senderID), uint(recieverID))
	response := &pb.AcceptRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) DeclineRequest(ctx context.Context, request *pb.DeclineRequestRequest) (*pb.DeclineRequestResponse, error) {
	senderID := request.SenderID
	recieverID := request.RecieverID
	senderID = int64(GetUserID(ctx))
	handler.requestsService.DeclineRequest(uint(senderID), uint(recieverID))
	response := &pb.DeclineRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) SendRequest(ctx context.Context, request *pb.SendRequestRequest) (*pb.SendRequestResponse, error) {
	senderID := request.SenderID
	recieverID := request.RecieverID
	senderID = int64(GetUserID(ctx))
	handler.requestsService.SendRequest(uint(senderID), uint(recieverID))
	response := &pb.SendRequestResponse{}
	return response, nil
}

func (handler *RequestsHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	senderID := request.SenderID
	receiverID := request.RecieverID
	senderID = int64(GetUserID(ctx))
	message := request.Message.Text
	handler.requestsService.SendMessage(uint(senderID), uint(receiverID), message)
	response := &pb.SendMessageResponse{}
	return response, nil
}
