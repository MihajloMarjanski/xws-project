package handler_grpc

import (
	"requests-service/model"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
)

func mapRequestToProto(request model.Request) *pb.Request {
	requestPb := &pb.Request{
		SenderId:   int64(request.SenderID),
		ReceiverId: int64(request.ReceiverID),
	}
	return requestPb
}

func mapUserToProto(user model.User) *pb.User {
	userPb := &pb.User{
		Id:        int64(user.ID),
		Username:  user.UserName,
		Biography: user.Biography,
		Name:      user.Name,
	}
	return userPb
}

func mapMessageToProto(message model.Message) *pb.Message {
	messagePb := &pb.Message{
		Text:       message.Text,
		ReceiverId: int64(message.ReceiverId),
		SenderId:   int64(message.SenderId),
	}
	return messagePb
}
