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
				log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Verify"}).Warn("Unexpected token signing method.")
				return nil, fmt.Errorf("unexpected token signing method")
			}

			log.WithFields(log.Fields{"service_name": "user-service","method_name": "Verify"}).Info("Token successfully verified.")
			return []byte("tajni_kljuc_za_jwt_hash"), nil
		},
	)

	if err != nil {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Verify"}).Warn("Invalid token.")
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*service.Claims)
	if !ok {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Verify"}).Warn("Invalid token claims.")
		return nil, fmt.Errorf("invalid token claims")
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Verify"}).Info("Token successfully verified.")
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
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "NewUserHandler"}).Error("Error creating user service.")
		return nil, err
	}

	return &UserHandler{
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "NewUserHandler"}).Info("Successfully created user handler.")
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
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetUser"}).Warn("User with that id does not exist.")
		return nil, err
	}
	userPb := mapUserDtoToProto(user)
	response := &pb.GetUserResponse{
		User: userPb,
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetUser"}).Info("User successfully retrieved.")
	return response, nil
}

func (handler *UserHandler) GetUserByUsername(ctx context.Context, request *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	username := request.Username
	user := handler.userService.GetByUsername(username)
	if user.ID == 0 {
		err := status.Error(codes.NotFound, "User with that username does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetUserByUsername"}).Warn("User with that username does not exist.")
		return nil, err
	}
	userPb := mapUserToProto(user)
	response := &pb.GetUserByUsernameResponse{
		User: userPb,
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetUserByUsername"}).Info("User successfully retrieved.")
	return response, nil
}

func (handler *UserHandler) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := mapProtoToUser(request.User)
	user.ID = GetUserID(ctx)
	if handler.userService.GetByID(int(user.ID)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "UpdateUser"}).Warn("User with that id does not exist..")
		return nil, err
	}
	id := handler.userService.UpdateUser(user.ID, user.Name, user.Email, user.Password, user.UserName, user.Gender, user.PhoneNumber, user.DateOfBirth, user.Biography, user.IsPrivate)
	response := &pb.UpdateUserResponse{
		Id: int64(id),
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "UpdateUser"}).Info("User successfully updated.")
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
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "CreateUser"}).Warn("User with same username or email already exists.")
		return nil, err
	}
	response := &pb.CreateUserResponse{
		Id: int64(id),
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "CreateUser"}).Info("User successfully created.")
	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.AddExperienceRequest) (*pb.AddExperienceResponse, error) {
	experience := mapProtoToExperience(request.Experience)
	experience.UserID = GetUserID(ctx)
	id := handler.userService.AddExperience(experience.Company, experience.Position, experience.From, experience.Until, experience.UserID)
	response := &pb.AddExperienceResponse{
		Id: int64(id),
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "AddExperience"}).Info("Experience successfully added.")
	return response, nil
}

func (handler *UserHandler) AddInterest(ctx context.Context, request *pb.AddInterestRequest) (*pb.AddInterestResponse, error) {
	interest := mapProtoToInterest(request.Interest)
	interest.UserID = GetUserID(ctx)
	id := handler.userService.AddInterest(interest.Interest, interest.UserID)
	response := &pb.AddInterestResponse{
		Id: int64(id),
	}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "AddInterest"}).Info("Interest successfully added.")
	return response, nil
}

func (handler *UserHandler) RemoveExperience(ctx context.Context, request *pb.RemoveExperienceRequest) (*pb.RemoveExperienceResponse, error) {
	id := request.Id
	handler.userService.RemoveExperience(int(id))
	response := &pb.RemoveExperienceResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "RemoveExperience"}).Info("Experience successfully removed.")
	return response, nil
}

func (handler *UserHandler) RemoveInterest(ctx context.Context, request *pb.RemoveInterestRequest) (*pb.RemoveInterestResponse, error) {
	id := request.Id
	handler.userService.RemoveInterest(int(id))
	response := &pb.RemoveInterestResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "RemoveInterest"}).Info("Interest successfully removed.")
	return response, nil
}
func (handler *UserHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	creds := request.Credentials
	token, ok := handler.userService.Login(creds.Username, creds.Password, creds.Pin)
	if !ok {
		err := status.Error(codes.NotFound, token)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login"}).Warn("Failed to login.")
		return nil, err
	}
	response := &pb.LoginResponse{
		Token: token,
	}

	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login"}).Info("Login successful.")
	return response, nil
}
func (handler *UserHandler) BlockUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	userId := GetUserID(ctx)
	blockedUserId := request.BlockedUserId
	if handler.userService.GetByID(int(userId)).ID == 0 || handler.userService.GetByID(int(blockedUserId)).ID == 0 {
		err := status.Error(codes.NotFound, "User with that id does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "BlockUser"}).Warn("User with that id does not exist.")
		return nil, err
	}
	handler.userService.BlockUser(int(userId), int(blockedUserId))
	response := &pb.BlockUserResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "BlockUser"}).Info("User successfully blocked.")
	return response, nil
}
func (handler *UserHandler) GetApiKey(ctx context.Context, request *pb.ApiKeyRequest) (*pb.ApiKeyResponse, error) {
	key := handler.userService.GetApiKeyForUserCredentials(request.Username, request.Password)
	if key == "" {
		err := status.Error(codes.NotFound, "User with that username and password does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKey"}).Warn("User with that username and password does not exist.")
		return nil, err
	}
	response := &pb.ApiKeyResponse{
		ApiKey: key,
	}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKey"}).Info("Api key successfully retrieved.")
	return response, nil
}
func (handler *UserHandler) GetApiKeyForUsername(ctx context.Context, request *pb.GetApiKeyForUsernameRequest) (*pb.GetApiKeyForUsernameResponse, error) {
	key := handler.userService.GetApiKeyForUsername(request.Username)
	if key == "" {
		err := status.Error(codes.NotFound, "User with that username does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKeyForUsername"}).Warn("User with that username does not exist.")
		return nil, err
	}
	response := &pb.GetApiKeyForUsernameResponse{
		ApiKey: key,
	}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKeyForUsername"}).Info("Api key successfully retrieved.")
	return response, nil
}
func (handler *UserHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {
	handler.userService.CreateJobOffer(int(request.Offer.Id), request.Offer.JobInfo, request.Offer.JobPosition, request.Offer.CompanyName, request.Offer.Qualifications, request.Offer.ApiKey)
	response := &pb.CreateJobOfferResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKeyForUsername"}).Info("Api key successfully retrieved.")
	return response, nil
}

func (handler *UserHandler) ActivateAccount(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	handler.userService.ActivateAccount(request.Token.Value)
	response := &pb.ActivateAccountResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKeyForUsername"}).Info("Api key successfully retrieved.")
	return response, nil
}

func (handler *UserHandler) GetPrivateStatusForUserId(ctx context.Context, request *pb.GetPrivateStatusForUserIdRequest) (*pb.GetPrivateStatusForUserIdResponse, error) {
	status := handler.userService.GetPrivateStatusForUserId(request.Id)
	response := &pb.GetPrivateStatusForUserIdResponse{
		IsPrivate: status,
	}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "GetApiKeyForUsername"}).Info("Api key successfully retrieved.")
	return response, nil
}

func (handler *UserHandler) ForgotPassword(ctx context.Context, request *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	i := handler.userService.ForgotPassword(request.Username)
	if i == 0 {
		err := status.Error(codes.InvalidArgument, "User with that username does not exist.")
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "ForgotPassword"}).Warn("User with that username does not exist.")
		return nil, err
	}
	response := &pb.ForgotPasswordResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "ForgotPassword"}).Info("Temporary password successfully created.")
	return response, nil
}

func (handler *UserHandler) SendPinFor2Auth(ctx context.Context, request *pb.SendPinFor2AuthRequest) (*pb.SendPinFor2AuthResponse, error) {
	message := handler.userService.SendPinFor2Auth(request.Credentials.Username, request.Credentials.Password)
	if message != "" {
		err := status.Error(codes.InvalidArgument, message)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth"}).Warn("Invalid pin for 2FA.")
		return nil, err
	}
	response := &pb.SendPinFor2AuthResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth"}).Info("Pin for 2FA successfully generated.")
	return response, nil
}

func (handler *UserHandler) SendPasswordlessToken(ctx context.Context, request *pb.SendPasswordlessTokenRequest) (*pb.SendPasswordlessTokenResponse, error) {
	message := handler.userService.SendPasswordlessToken(request.Username)
	if message != "" {
		err := status.Error(codes.InvalidArgument, message)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPasswordlessToken"}).Warn("Invalid passwordless token.")
		return nil, err
	}
	response := &pb.SendPasswordlessTokenResponse{}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPasswordlessToken"}).Info("Passwordless token successfully generated.")
	return response, nil
}

func (handler *UserHandler) LoginPasswordless(ctx context.Context, request *pb.LoginPasswordlessRequest) (*pb.LoginPasswordlessResponse, error) {
	id := int64(GetUserID(ctx))
	message, ok := handler.userService.LoginPasswordless(int(id))
	if !ok {
		err := status.Error(codes.NotFound, message)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "LoginPasswordless"}).Warn("Failed passwordless login.")
		return nil, err
	}
	response := &pb.LoginPasswordlessResponse{
		Jwt: message,
	}
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "LoginPasswordless"}).Warn("Passwordless login successfull.")
	return response, nil
}
