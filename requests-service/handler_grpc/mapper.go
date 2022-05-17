package handler_grpc

import (
	"requests-service/model"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
)

func mapRequestToProto(request model.Request) *pb.Request {
	requestPb := &pb.Request{
		SenderID:   int64(request.SenderID),
		RecieverID: int64(request.ReceiverID),
	}
	return requestPb
}
