package handler_grpc

import (
	"context"
	"fmt"
	"io"
	"os"
	"requests-service/model"
	"requests-service/service"
	"strconv"
	"strings"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	"github.com/dgrijalva/jwt-go"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type RequestsHandler struct {
	pb.UnimplementedRequestsServiceServer
	requestsService *service.RequestsService
}

func init() {

	f := &lumberjack.Logger{
		Filename:   "./testlogrus.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)
}

func Verify(accessToken string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&model.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.WithFields(log.Fields{"service_name": "request-service", "method_name": "Verify"}).Error("Unexpected token signing method.")
				return nil, fmt.Errorf("unexpected token signing method")
			}

			log.WithFields(log.Fields{"service_name": "request-service", "method_name": "Verify"}).Info("Token successfully verified.")
			return []byte("tajni_kljuc_za_jwt_hash"), nil
		},
	)

	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "Verify"}).Warn("Invalid token.")
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "Verify"}).Warn("Invalid token claims.")
		return nil, fmt.Errorf("invalid token claims")
	}

	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "Verify"}).Info("Token successfully verified.")
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
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "NewRequestsHandler"}).Error("Error creating request service.")
		return nil, err
	}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "NewRequestsHandler"}).Info("Successfully created request handler.")
	return &RequestsHandler{
		requestsService: requestsService,
	}, nil
}

func (requestsHandler *RequestsHandler) CloseDB() error {

	return requestsHandler.requestsService.CloseDB()
}

func (handler *RequestsHandler) GetAllByRecieverId(ctx context.Context, request *pb.GetAllByRecieverIdRequest) (*pb.GetAllByRecieverIdResponse, error) {
	idReceived := request.ReceiverId
	id := int64(GetUserID(ctx))
	if idReceived != id {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "NewRequestsHandler"}).Warn("Someone tried to pose as different user.")
	}
	users := handler.requestsService.GetAllByRecieverId(uint(id))

	response := &pb.GetAllByRecieverIdResponse{
		Users: users,
	}
	return response, nil
}

func (handler *RequestsHandler) AcceptRequest(ctx context.Context, request *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
	senderID := request.SenderId
	recieverID := request.ReceiverId
	recieverID = int64(GetUserID(ctx))
	handler.requestsService.AcceptRequest(uint(senderID), uint(recieverID))
	response := &pb.AcceptRequestResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "AcceptRequest"}).Info("Successfully accepted request.")
	return response, nil
}

func (handler *RequestsHandler) DeclineRequest(ctx context.Context, request *pb.DeclineRequestRequest) (*pb.DeclineRequestResponse, error) {
	senderID := request.SenderId
	recieverID := request.ReceiverId
	recieverID = int64(GetUserID(ctx))
	handler.requestsService.DeclineRequest(uint(senderID), uint(recieverID))
	response := &pb.DeclineRequestResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "DeclineRequest"}).Info("Successfully denied request.")
	return response, nil
}

func (handler *RequestsHandler) SendRequest(ctx context.Context, request *pb.SendRequestRequest) (*pb.SendRequestResponse, error) {
	senderID := request.SenderId
	recieverID := request.ReceiverId
	senderID = int64(GetUserID(ctx))
	handler.requestsService.SendRequest(uint(senderID), uint(recieverID))
	response := &pb.SendRequestResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendRequest"}).Info("Successfully sent request.")
	return response, nil
}

func (handler *RequestsHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	senderID := request.SenderId
	receiverID := request.ReceiverId
	senderID = int64(GetUserID(ctx))
	message := request.Message.Text
	handler.requestsService.SendMessage(uint(senderID), uint(receiverID), message)
	response := &pb.SendMessageResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendMessage"}).Info("Successfully sent message.")
	return response, nil
}

func (handler *RequestsHandler) AreConnected(ctx context.Context, request *pb.AreConnectedRequest) (*pb.AreConnectedResponse, error) {
	res := handler.requestsService.AreConnected(request.FirstId, request.SecondId)
	response := &pb.AreConnectedResponse{
		AreConnected: res,
	}
	return response, nil
}

func (handler *RequestsHandler) FindConnections(ctx context.Context, request *pb.FindConnectionsRequest) (*pb.FindConnectionsResponse, error) {
	var users []*pb.User
	for _, user := range handler.requestsService.FindConnections(request.Id) {
		users = append(users, mapUserToProto(user))
	}
	response := &pb.FindConnectionsResponse{
		Users: users,
	}
	return response, nil
}

func (handler *RequestsHandler) FindMessages(ctx context.Context, request *pb.FindMessagesRequest) (*pb.FindMessagesResponse, error) {
	var messages []*pb.Message
	for _, message := range handler.requestsService.FindMessages(request.Id1, request.Id2) {
		messages = append(messages, mapMessageToProto(message))
	}
	response := &pb.FindMessagesResponse{
		Messages: messages,
	}
	return response, nil
}

func (handler *RequestsHandler) DeleteConnection(ctx context.Context, request *pb.DeleteConnectionRequest) (*pb.DeleteConnectionResponse, error) {
	handler.requestsService.DeleteConnection(request.Id1, request.Id2)
	response := &pb.DeleteConnectionResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "DeleteConnection"}).Info("Successfully deleted connection.")
	return response, nil
}

func (handler *RequestsHandler) GetNotifications(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	var notifications []*pb.Notification
	for _, notification := range handler.requestsService.GetNotifications(request.Id) {
		notifications = append(notifications, mapNotificationToProto(notification))
	}
	response := &pb.GetNotificationsResponse{
		Notifications: notifications,
	}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "GetNotifications"}).Info("Successfully gotten notifications.")
	return response, nil
}

func (handler *RequestsHandler) SendNotification(ctx context.Context, request *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	senderID := request.SenderId
	receiverID := request.ReceiverId
	message := request.Message
	handler.requestsService.SendNotification(uint(senderID), uint(receiverID), message)
	response := &pb.SendNotificationResponse{}
	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendNotification"}).Info("Successfully sent notification.")
	return response, nil
}
