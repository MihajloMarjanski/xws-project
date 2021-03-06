package handler_grpc

import (
	"time"
	"user-service/model"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
)

func mapUserDtoToProto(user model.User) *pb.User {
	userPb := &pb.User{
		Id:        int64(user.ID),
		Name:      user.Name,
		Username:  user.UserName,
		Email:     user.Email,
		Gender:    string(user.Gender),
		Phone:     user.PhoneNumber,
		Date:      user.DateOfBirth.Format("02-Jan-2006"),
		Biography: user.Biography,
		IsPrivate: user.IsPrivate,
	}

	for _, interest := range user.Interests {
		userPb.Interests = append(userPb.Interests, &pb.Interest{
			Id:       int64(interest.ID),
			Interest: interest.Interest,
			UserId:   int64(interest.UserID),
		})
	}
	for _, experience := range user.Experiences {
		userPb.Experience = append(userPb.Experience, &pb.Experience{
			Id:       int64(experience.ID),
			Company:  experience.Company,
			Position: experience.Position,
			From:     experience.From.String(),
			Until:    experience.Until.String(),
			UserId:   int64(experience.UserID),
		})
	}

	return userPb
}

func mapUserOfferProto(offer model.JobOffer) *pb.JobOffer {
	offerPb := &pb.JobOffer{
		Id:             int64(offer.ID),
		JobInfo:        offer.JobInfo,
		ApiKey:         offer.ApiKey,
		Qualifications: offer.Qualifications,
		CompanyName:    offer.CompanyName,
		JobPosition:    offer.JobPosition,
	}
	return offerPb
}

func mapUserToProto(user model.User) *pb.UserWithPass {
	userPb := &pb.UserWithPass{
		Id:        int64(user.ID),
		Name:      user.Name,
		Username:  user.UserName,
		Email:     user.Email,
		Gender:    string(user.Gender),
		Phone:     user.PhoneNumber,
		Date:      user.DateOfBirth.Format("02-Jan-2006"),
		Biography: user.Biography,
		Password:  user.Password,
		IsPrivate: user.IsPrivate,
	}

	for _, interest := range user.Interests {
		userPb.Interests = append(userPb.Interests, &pb.Interest{
			Id:       int64(interest.ID),
			Interest: interest.Interest,
			UserId:   int64(interest.UserID),
		})
	}
	for _, experience := range user.Experiences {
		userPb.Experience = append(userPb.Experience, &pb.Experience{
			Id:       int64(experience.ID),
			Company:  experience.Company,
			Position: experience.Position,
			From:     experience.From.Format("02-Jan-2006"),
			Until:    experience.Until.Format("02-Jan-2006"),
			UserId:   int64(experience.UserID),
		})
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
		IsPrivate:   user.IsPrivate,
	}
	return userPb
}

func mapProtoToExperience(experience *pb.Experience) model.Experience {
	dateFrom, _ := time.Parse(time.RFC3339, experience.From)
	dateUntil, _ := time.Parse(time.RFC3339, experience.Until)
	expPb := model.Experience{
		ID:       uint(experience.Id),
		Company:  experience.Company,
		Position: experience.Position,
		From:     dateFrom,
		Until:    dateUntil,
		UserID:   uint(experience.UserId),
	}
	return expPb
}

func mapProtoToInterest(interest *pb.Interest) model.Interest {
	intPb := model.Interest{
		ID:       uint(interest.Id),
		Interest: interest.Interest,
		UserID:   uint(interest.UserId),
	}
	return intPb
}
