package service

import (
	pbUser "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"requests-service/model"
	"requests-service/repo"
)

type RequestsService struct {
	reqRepo *repo.RequestsRepository
}

func New() (*RequestsService, error) {
	reqRepo, err := repo.New()
	if err != nil {
		return nil, err
	}

	return &RequestsService{
		reqRepo: reqRepo,
	}, nil
}

func (s *RequestsService) CloseDB() error {
	return s.reqRepo.Close()
}

func (s *RequestsService) GetAllByRecieverId(rid uint) []model.Request {
	return s.reqRepo.GetAllByRecieverId(rid)
}

func (s *RequestsService) AcceptRequest(sid, rid uint) {
	s.reqRepo.AcceptRequest(sid, rid)
}

func (s *RequestsService) DeclineRequest(sid, rid uint) {
	s.reqRepo.DeclineRequest(sid, rid)
}

func (s *RequestsService) SendRequest(sid, rid uint) {
	s.reqRepo.SendRequest(sid, rid)
}

func (s *RequestsService) SendMessage(senderID, receiverID uint, message string) {
	if s.reqRepo.AreConnected(senderID, receiverID) == true {
		s.reqRepo.SendMessage(senderID, receiverID, message)
	}
	return
}

func (s *RequestsService) AreConnected(id1 int64, id2 int64) bool {
	return s.reqRepo.AreConnected(uint(id1), uint(id2)) || s.reqRepo.AreRequested(uint(id1), uint(id2))
}

func (s *RequestsService) FindConnections(id int64) []model.User {
	ids1, ids2 := s.reqRepo.GetAllConnections(uint(id))
	var res []model.User

	conn, err := grpc.Dial("localhost:8100", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbUser.NewUserServiceClient(conn)

	for _, connection := range ids1 {
		response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(connection.UserOne)})
		if err != nil {
			panic(err)
		}
		res = append(res, mapUser(response.User))
	}
	for _, connection := range ids2 {
		response, err := client.GetUser(context.Background(), &pbUser.GetUserRequest{Id: int64(connection.UserTwo)})
		if err != nil {
			panic(err)
		}
		res = append(res, mapUser(response.User))
	}

	return res
}

func (s *RequestsService) FindMessages(id1 int64, id2 int64) []model.Message {
	return s.reqRepo.FindMessages(id1, id2)
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
