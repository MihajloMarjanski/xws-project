package handler_grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-service/service"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func New() (*UserHandler, error) {

	userService, err := service.New()
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		userService: userService,
	}, nil
}

func (userHandler *UserHandler) CloseDB() error {

	return userHandler.userService.CloseDB()
}

func (handler *UserHandler) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := request.Id
	user := handler.userService.GetByID(int(id))
	if user.ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		return nil, err
	}
	userPb := mapUserDtoToProto(user)
	response := &pb.GetUserResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) GetMe(ctx context.Context, request *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	id := request.Id
	user := handler.userService.GetMe(int(id))
	if user.ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		return nil, err
	}
	userPb := mapUserToProto(user)
	response := &pb.GetMeResponse{
		User: userPb,
	}

	//response1, err := http.Get("http://localhost:8600/company/owner/1")
	//if err != nil {
	//	fmt.Print(err.Error())
	//	os.Exit(1)
	//}
	//
	//responseData, err := ioutil.ReadAll(response1.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("dobio: ", string(responseData))
	//log.Println("dobio pre: ", response1)

	return response, nil
}

func (handler *UserHandler) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := mapProtoToUser(request.User)
	if handler.userService.GetByID(int(user.ID)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		return nil, err
	}
	id := handler.userService.UpdateUser(user.ID, user.Name, user.Email, user.Password, user.UserName, user.Gender, user.PhoneNumber, user.DateOfBirth, user.Biography, user.IsPublic)
	response := &pb.UpdateUserResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) SearchUsers(ctx context.Context, request *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	username := request.Username
	var users []*pb.User
	for _, user := range handler.userService.SearchUsers(username) {
		users = append(users, mapUserDtoToProto(user))
	}
	response := &pb.SearchUsersResponse{
		Users: users,
	}
	return response, nil
}

func (handler *UserHandler) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := mapProtoToUser(request.User)
	id := handler.userService.CreateUser(user.Name, user.Email, user.Password, user.UserName, user.Gender, user.PhoneNumber, user.DateOfBirth, user.Biography)
	if id == 0 {
		err := status.Error(codes.AlreadyExists, "User with same username or email already exists.")
		return nil, err
	}
	response := &pb.CreateUserResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.AddExperienceRequest) (*pb.AddExperienceResponse, error) {
	experience := mapProtoToExperience(request.Experience)
	id := handler.userService.AddExperience(experience.Company, experience.Position, experience.From, experience.Until, experience.UserID)
	response := &pb.AddExperienceResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) AddInterest(ctx context.Context, request *pb.AddInterestRequest) (*pb.AddInterestResponse, error) {
	interest := mapProtoToInterest(request.Interest)
	id := handler.userService.AddInterest(interest.Interest, interest.UserID)
	response := &pb.AddInterestResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) RemoveExperience(ctx context.Context, request *pb.RemoveExperienceRequest) (*pb.RemoveExperienceResponse, error) {
	id := request.Id
	handler.userService.RemoveExperience(int(id))
	response := &pb.RemoveExperienceResponse{}
	return response, nil
}

func (handler *UserHandler) RemoveInterest(ctx context.Context, request *pb.RemoveInterestRequest) (*pb.RemoveInterestResponse, error) {
	id := request.Id
	handler.userService.RemoveInterest(int(id))
	response := &pb.RemoveInterestResponse{}
	return response, nil
}
func (handler *UserHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	creds := request.Credentials
	token := handler.userService.Login(creds.Username, creds.Password)
	response := &pb.LoginResponse{
		Token: token,
	}
	return response, nil
}
func (handler *UserHandler) BlockUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	userId := request.UserId
	blockedUserId := request.BlockedUserId
	if handler.userService.GetByID(int(userId)).ID == 0 || handler.userService.GetByID(int(blockedUserId)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		return nil, err
	}
	handler.userService.BlockUser(int(userId), int(blockedUserId))
	response := &pb.BlockUserResponse{}
	return response, nil
}
