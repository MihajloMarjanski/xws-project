package handler_grpc

import (
	"user-service/model"
	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
)

func mapUser(user *model.User) *pb.User {
	userPb := &pb.User{
		Id:           int64(user.ID),
		Name:          user.Name,
		Username:	 user.UserName,
		Email: user.Email,
		Gender: string(user.Gender),
		Phone: user.PhoneNumber,
		Date: user.DateOfBirth,
		Biography: user.Biography,
	}
	// for _, interest := range user.Interests {
	// 	userPb.Interests = append(userPb.Interests, &pb.Interest{
	// 		ID: interest.ID,
	// 		Interest: interest.Interest,
	// 		UserID: interest.UserID,
	// 	})
	// }
	// for _, experience := range user.Experiences {
	// 	userPb.Experiences = append(userPb.Experiences, &pb.Experience{
	// 		ID: experience.ID,
	// 		Company: experience.Company,
	// 		Position: experience.Position,
	// 		From: experience.From,
	// 		Until: experience.Until,
	// 		UserID: experience.UserID,
	// 	})
	// }
	return userPb
}