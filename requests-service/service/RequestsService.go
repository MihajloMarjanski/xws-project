package service

import (
	"io"
	"os"
	"requests-service/model"
	"requests-service/repo"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	pbUser "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RequestsService struct {
	reqRepo *repo.RequestsRepository
}

func init() {

	f := &lumberjack.Logger{
		Filename:   "./logs/testlogrus.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)
}

func New() (*RequestsService, error) {
	reqRepo, err := repo.New()
	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "NewRequestService"}).Error("Error creating Request Repository.")
		return nil, err
	}

	log.WithFields(log.Fields{"service_name": "request-service", "method_name": "NewRequestService"}).Info("Successfully created Request Service.")
	return &RequestsService{
		reqRepo: reqRepo,
	}, nil
}

func (s *RequestsService) CloseDB() error {
	return s.reqRepo.Close()
}

func (s *RequestsService) GetAllByRecieverId(rid uint) []*pb.UsernameWithRequestId {
	var users []*pb.UsernameWithRequestId

	conn, err := grpc.Dial("user-service:8100", grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "GetAllByRecieverId", "reciever_id": rid}).Error("Error establishing connection with user-service via grpc.")
		panic(err)
	}
	defer conn.Close()
	client := pbUser.NewUserServiceClient(conn)

	for _, request := range s.reqRepo.GetAllByRecieverId(rid) {
		response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(request.SenderID)})
		if err != nil {
			log.WithFields(log.Fields{"service_name": "request-service", "method_name": "GetAllByRecieverId", "reciever_id": rid}).Warn("Error getting user.")
			panic(err)
		}
		users = append(users, mapUserToUsernameReq(response.User, request.ReceiverID, request.SenderID))
	}

	return users
}

func (s *RequestsService) AcceptRequest(sid, rid uint) {
	s.reqRepo.AcceptRequest(sid, rid)
}

func (s *RequestsService) DeclineRequest(sid, rid uint) {
	s.reqRepo.DeclineRequest(sid, rid)
}

func (s *RequestsService) SendRequest(sid, rid uint) {
	s.reqRepo.SendRequest(sid, rid)
	s.SendNotification(sid, rid, "You have received a new connection request from user '")
}

func (s *RequestsService) SendMessage(senderID, receiverID uint, message string) {
	if s.reqRepo.AreConnected(senderID, receiverID) == true {
		s.reqRepo.SendMessage(senderID, receiverID, message)
		s.SendNotification(senderID, receiverID, "You have received a new message from user '")
	}
	return
}

func (s *RequestsService) SendNotification(senderID, receiverID uint, message string) {
	conn, err := grpc.Dial("user-service:8100", grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendNotification", "sender_id": senderID, "receiver_id": receiverID}).Error("Error establishing connection with user-service via grpc.")
		panic(err)
	}
	defer conn.Close()
	client := pbUser.NewUserServiceClient(conn)
	response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(senderID)})
	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendNotification", "sender_id": senderID, "receiver_id": receiverID}).Warn("Error getting user")
		panic(err)
	}

	s.reqRepo.SendNotification(receiverID, message+response.User.Username+"'")
}

func (s *RequestsService) AreConnected(id1 int64, id2 int64) bool {
	return s.reqRepo.AreConnected(uint(id1), uint(id2)) || s.reqRepo.AreRequested(uint(id1), uint(id2))
}

func (s *RequestsService) FindConnections(id int64) []model.User {
	ids1, ids2 := s.reqRepo.GetAllConnections(uint(id))
	var res []model.User

	conn, err := grpc.Dial("user-service:8100", grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"service_name": "request-service", "method_name": "FindConnections", "user_id": id}).Error("Error establishing connection with user-service via grpc.")
		panic(err)
	}
	defer conn.Close()
	client := pbUser.NewUserServiceClient(conn)

	for _, connection := range ids1 {
		response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(connection.UserOne)})
		if err != nil {
			log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendNotification", "user_id": id}).Warn("Error getting user")
			panic(err)
		}
		res = append(res, mapUser(response.User))
	}
	for _, connection := range ids2 {
		response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(connection.UserTwo)})
		if err != nil {
			log.WithFields(log.Fields{"service_name": "request-service", "method_name": "SendNotification", "user_id": id}).Warn("Error getting user")
			panic(err)
		}
		res = append(res, mapUser(response.User))
	}

	return res
}

func (s *RequestsService) FindMessages(id1 int64, id2 int64) []model.Message {
	return s.reqRepo.FindMessages(id1, id2)
}

func (s *RequestsService) DeleteConnection(id1 int64, id2 int64) {
	s.reqRepo.DeleteConnection(id1, id2)
	return
}

func (s *RequestsService) GetNotifications(id int64) []model.Notification {
	return s.reqRepo.GetNotifications(id)
}

func mapUser(user *pbUser.User) model.User {
	res := model.User{
		ID:        uint(user.Id),
		UserName:  user.Username,
		Biography: user.Biography,
		Name:      user.Name,
	}
	return res
}

func mapUserToUsernameReq(user *pbUser.User, receiverId, senderId uint) *pb.UsernameWithRequestId {
	res := &pb.UsernameWithRequestId{
		ReceiverId: int64(receiverId),
		SenderId:   int64(senderId),
		Username:   user.Username,
	}
	return res
}
