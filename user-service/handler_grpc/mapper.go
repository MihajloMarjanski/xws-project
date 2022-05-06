package handler_grpc

import (
	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"time"
	"user-service/model"
)

func mapUserDtoToProto(user model.UserDTO) *pb.User {
	userPb := &pb.User{
		Id:        int64(user.ID),
		Name:      user.Name,
		Username:  user.UserName,
		Email:     user.Email,
		Gender:    string(user.Gender),
		Phone:     user.PhoneNumber,
		Date:      user.DateOfBirth.String(),
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

func mapUserToProto(user model.User) *pb.UserWithPass {
	userPb := &pb.UserWithPass{
		Id:        int64(user.ID),
		Name:      user.Name,
		Username:  user.UserName,
		Email:     user.Email,
		Gender:    string(user.Gender),
		Phone:     user.PhoneNumber,
		Date:      user.DateOfBirth.String(),
		Biography: user.Biography,
		Password:  user.Password,
	}
	return userPb
}

func mapProtoToUser(user *pb.UserWithPass) model.User {
	date, _ := time.Parse(time.RFC3339, user.Date)
	userPb := model.User{
		ID:          uint(user.Id),
		Name:        user.Name,
		UserName:    user.Username,
		Email:       user.Email,
		Gender:      model.Gender(user.Gender),
		PhoneNumber: user.Phone,
		DateOfBirth: date,
		Biography:   user.Biography,
		Password:    user.Password,
	}
	return userPb
}
