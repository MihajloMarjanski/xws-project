package handler_grpc

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"user-service/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"github.com/dgrijalva/jwt-go"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
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

func Verify(accessToken string) (*service.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&service.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.Error("Unexpected token signing method.")
				return nil, fmt.Errorf("unexpected token signing method")
			}

			log.Info("Token successfully verified.")
			return []byte("tajni_kljuc_za_jwt_hash"), nil
		},
	)

	if err != nil {
		log.Error("Invalid token.")
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*service.Claims)
	if !ok {
		log.Error("Invalid token claims.")
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

func New() (*UserHandler, error) {

	userService, err := service.New()
	if err != nil {
		log.Error("Error creating User Handler.")
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
		log.WithFields(log.Fields{"method_name": "GetUser", "method_type": "GET"}).Warn("User with that id does not exist.")
		return nil, err
	}
	userPb := mapUserDtoToProto(user)
	response := &pb.GetUserResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) GetUserByUsername(ctx context.Context, request *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	username := request.Username
	user := handler.userService.GetByUsername(username)
	if user.ID == 0 {
		err := status.Error(codes.NotFound, "User with that username does not exist.")
		log.Error("User with that username does not exist.")
		return nil, err
	}
	userPb := mapUserToProto(user)
	response := &pb.GetUserByUsernameResponse{
		User: userPb,
	}

	return response, nil
}

func (handler *UserHandler) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := mapProtoToUser(request.User)
	user.ID = GetUserID(ctx)
	if handler.userService.GetByID(int(user.ID)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		log.Error("User with that username does not exist.")
		return nil, err
	}
	id := handler.userService.UpdateUser(user.ID, user.Name, user.Email, user.Password, user.UserName, user.Gender, user.PhoneNumber, user.DateOfBirth, user.Biography, user.IsPrivate)
	response := &pb.UpdateUserResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) SearchUsers(ctx context.Context, request *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	username := request.Username
	loggedUserId := request.LoggedUserId
	var users []*pb.User
	for _, user := range handler.userService.SearchUsers(username, uint(loggedUserId)) {
		users = append(users, mapUserDtoToProto(user))
	}
	response := &pb.SearchUsersResponse{
		Users: users,
	}
	return response, nil
}

func (handler *UserHandler) SearchOffers(ctx context.Context, request *pb.SearchOffersRequest) (*pb.SearchOffersResponse, error) {
	var offers []*pb.JobOffer
	for _, offer := range handler.userService.SearchOffers(request.Text) {
		offers = append(offers, mapUserOfferProto(offer))
	}
	response := &pb.SearchOffersResponse{
		Offers: offers,
	}
	return response, nil
}

func (handler *UserHandler) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := mapProtoToUser(request.User)
	id := handler.userService.CreateUser(user.Name, user.Email, user.Password, user.UserName, user.Gender, user.PhoneNumber, user.DateOfBirth, user.Biography)
	if id == 0 {
		err := status.Error(codes.AlreadyExists, "User with same username or email already exists.")
		log.Error("User with same username or email already exists.")
		return nil, err
	}
	response := &pb.CreateUserResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.AddExperienceRequest) (*pb.AddExperienceResponse, error) {
	experience := mapProtoToExperience(request.Experience)
	experience.UserID = GetUserID(ctx)
	id := handler.userService.AddExperience(experience.Company, experience.Position, experience.From, experience.Until, experience.UserID)
	response := &pb.AddExperienceResponse{
		Id: int64(id),
	}
	return response, nil
}

func (handler *UserHandler) AddInterest(ctx context.Context, request *pb.AddInterestRequest) (*pb.AddInterestResponse, error) {
	interest := mapProtoToInterest(request.Interest)
	interest.UserID = GetUserID(ctx)
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
	userId := GetUserID(ctx)
	blockedUserId := request.BlockedUserId
	if handler.userService.GetByID(int(userId)).ID == 0 || handler.userService.GetByID(int(blockedUserId)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		log.Error("User with that id does not exist.")
		return nil, err
	}
	handler.userService.BlockUser(int(userId), int(blockedUserId))
	response := &pb.BlockUserResponse{}
	return response, nil
}
func (handler *UserHandler) GetApiKey(ctx context.Context, request *pb.ApiKeyRequest) (*pb.ApiKeyResponse, error) {
	key := handler.userService.GetApiKeyForUserCredentials(request.Username, request.Password)
	if key == "" {
		err := status.Error(codes.NotFound, "User with that username and password does not exist.")
		log.Error("User with that username and password does not exist.")
		return nil, err
	}
	response := &pb.ApiKeyResponse{
		ApiKey: key,
	}
	return response, nil
}
func (handler *UserHandler) GetApiKeyForUsername(ctx context.Context, request *pb.GetApiKeyForUsernameRequest) (*pb.GetApiKeyForUsernameResponse, error) {
	key := handler.userService.GetApiKeyForUsername(request.Username)
	if key == "" {
		err := status.Error(codes.NotFound, "User with that username does not exist.")
		log.Error("User with that username does not exist.")
		return nil, err
	}
	response := &pb.GetApiKeyForUsernameResponse{
		ApiKey: key,
	}
	return response, nil
}
func (handler *UserHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {
	handler.userService.CreateJobOffer(int(request.Offer.Id), request.Offer.JobInfo, request.Offer.JobPosition, request.Offer.CompanyName, request.Offer.Qualifications, request.Offer.ApiKey)
	response := &pb.CreateJobOfferResponse{}
	return response, nil
}

func (handler *UserHandler) ActivateAccount(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	handler.userService.ActivateAccount(request.Token.Value)
	response := &pb.ActivateAccountResponse{}
	return response, nil
}

func (handler *UserHandler) GetPrivateStatusForUserId(ctx context.Context, request *pb.GetPrivateStatusForUserIdRequest) (*pb.GetPrivateStatusForUserIdResponse, error) {
	status := handler.userService.GetPrivateStatusForUserId(request.Id)
	response := &pb.GetPrivateStatusForUserIdResponse{
		IsPrivate: status,
	}
	return response, nil
}
