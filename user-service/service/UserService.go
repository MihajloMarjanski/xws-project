package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	pbConn "github.com/MihajloMarjanski/xws-project/common/proto/connection_service"
	pbReq "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
	"github.com/MihajloMarjanski/xws-project/common/saga/messaging/nats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"user-service/model"
	"user-service/repo"

	"github.com/dgrijalva/jwt-go"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	QueueGroup = "user_service"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
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

type Claims struct {
	Username    string   `json:"username"`
	Id          string   `json:"id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

var jwtKey = []byte("tajni_kljuc_za_jwt_hash")

type UserService struct {
	userRepo     *repo.UserRepository
	orchestrator *BlockUserOrchestrator
}

func initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}
func initCreateOrderOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *BlockUserOrchestrator {
	orchestrator, err := NewBlockUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func New() (*UserService, error) {

	userRepo, err := repo.New()
	if err != nil {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "NewUserService"}).Error("Error creating User Repository.")
		return nil, err
	}

	commandPublisher := initPublisher("block.user.command")
	replySubscriber := initSubscriber("block.user.reply", QueueGroup)
	blockUserOrchestrator := initCreateOrderOrchestrator(commandPublisher, replySubscriber)
	log.WithFields(log.Fields{"service_name": "user-service", "method_name": "NewUserService"}).Info("Successfully created User Service.")
	return &UserService{
		userRepo:     userRepo,
		orchestrator: blockUserOrchestrator,
	}, nil
}

func (s *UserService) SearchUsers(username string, id uint) []model.User {
	users := s.userRepo.SearchUsers(username)
	if id == 0 {
		return users
	}
	blockedIds := s.userRepo.FindBlockedForUserId(id)
	i := 0
	for _, user := range users {
		if !s.userRepo.Contains(blockedIds, user.ID) && user.ID != id {
			users[i] = user
			i++
		}
	}
	users = users[:i]
	return users
}

func (s *UserService) GetByID(id int) model.User {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetByUsername(username string) model.User {
	return s.userRepo.GetByUsername(username)
}

func (s *UserService) GetMe(id int) model.User {
	return s.userRepo.GetMe(id)
}

func (s *UserService) RemoveExperience(id int) int {
	return s.userRepo.RemoveExperience(id)
}

func (s *UserService) RemoveInterest(id int) int {
	return s.userRepo.RemoveInterest(id)
}

func (s *UserService) CloseDB() error {
	return s.userRepo.Close()
}

func (s *UserService) CreateUser(name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string) int {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	apiKey, _ := bcrypt.GenerateFromPassword([]byte(GenerateRandomString(10)), 8)
	SendActivationMail(email, name, string(apiKey))
	return s.userRepo.CreateUser(name, email, string(hashedPassword), username, gender, phonenumber, dateofbirth, biography, string(apiKey))
}

func SendActivationMail(email string, name string, key string) {
	data := map[string]string{
		"email":  email,
		"name":   name,
		"apiKey": key,
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//_, err = http.Post("https://localhost:8600/email/activation", "application/json", bytes.NewBuffer(json_data))
	//_, err = http.Post("https://host.docker.internal:8600/email/activation", "application/json", bytes.NewBuffer(json_data))
	_, err = http.Post("https://agent-application:8600/email/activation", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Print(err.Error())
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendActivationMail"}).Error("Error sending activation mail.")
		os.Exit(1)
	}
}

func (s *UserService) AddInterest(interest string, userId uint) int {

	return s.userRepo.AddInterest(interest, userId)
}

func (s *UserService) AddExperience(company string, position string, from time.Time, until time.Time, userId uint) int {

	return s.userRepo.AddExperience(company, position, from, until, userId)
}

func (s *UserService) UpdateUser(id uint, name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string, isPrivate bool) int {
	if password != s.GetByID(int(id)).Password {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
		return s.userRepo.UpdateUser(id, name, email, string(hashedPassword), username, gender, phonenumber, dateofbirth, biography, isPrivate, true)
	}
	return s.userRepo.UpdateUser(id, name, email, password, username, gender, phonenumber, dateofbirth, biography, isPrivate, false)
}

func (s *UserService) Login(username string, password, pin string) (string, bool) {
	user := s.GetByUsername(username)
	if user.ID == 0 {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Warn("Username dosn't exist.")
		return "Wrong username", false
	} else if s.IsBlocked(user) {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Info("Account is currently blocked.")
		return "Your account is currently blocked. Try next day again.", false
	} else if !user.IsActivated {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Info("Account is not activated.")
		return "Your have to activate your profile first.", false
	}
	if user.Forgotten == 1 {
		user.Forgotten = 2
		s.userRepo.Save(user)
	} else if user.Forgotten == 2 {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Info("Password not changed after being forgotten.")
		return "You did not changed password first time. If you want to log in, refresh again your password.", false
	}

	expectedPassword := user.Password
	id := strconv.FormatUint(uint64(user.ID), 10)
	err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
	if err != nil {
		s.IncreaseMissedPasswordCounter(user)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Warn("Incorrect password.")
		return "Wrong password", false
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(pin))
	if err1 != nil {
		s.IncreaseMissedPasswordCounter(user)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Warn("Incorrect pin.")
		return "Wrong pin", false
	}

	s.RefreshMissedPasswordCounter(user)
	expirationTime := time.Now().Add(60 * time.Minute)

	permissions := []string{"GetAllByRecieverId", "AcceptRequest", "DeclineRequest", "SendRequest", "SendMessage",
		"FindMessages", "GetNotifications", "UpdateUser", "AddExperience", "RemoveExperience", "AddInterest", "RemoveInterest",
		"BlockUser", "GetUserByUsername"}

	claims := &Claims{
		Username:    username,
		Id:          id,
		Role:        "ROLE_USER",
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "Login", "username": username}).Warn("Invalid token.")
		return "", false
	}
	return tokenString, true
}

func (s *UserService) IsBlocked(user model.User) bool {
	if user.BlockedDate.IsZero() {
		return false
	}
	if user.IsBlocked && time.Now().Before(user.BlockedDate.AddDate(0, 0, 1)) {
		return true
	} else if user.IsBlocked && time.Now().After(user.BlockedDate.AddDate(0, 0, 1)) {
		user.IsBlocked = false
		user.MissedPasswordCounter = 0
		s.userRepo.Save(user)
		return false
	}
	return false
}

func (s *UserService) BlockUser(userId int, blockedUserId int) {
	user := s.GetByID(userId)
	block := s.GetByID(blockedUserId)
	if userId == blockedUserId || user.ID == 0 || block.ID == 0 {
		return
	}
	s.userRepo.BlockUser(userId, blockedUserId)

	err := s.orchestrator.Start(uint(userId), uint(blockedUserId))
	if err != nil {
		s.userRepo.UnblockUser(userId, blockedUserId)
		return
	}
	//DeleteConnection(userId, blockedUserId)
	return

}

func DeleteConnection(userId, id int) {
	crtTlsPath, _ := filepath.Abs("./service.pem")

	creds, err6 := credentials.NewClientTLSFromFile(crtTlsPath, "")
	if err6 != nil {
		//log.Fatalf("could not process the credentials: %v", err6)
	}
	conn, err := grpc.Dial("request-service:8200", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbReq.NewRequestsServiceClient(conn)
	_, err = client.DeleteConnection(context.Background(), &pbReq.DeleteConnectionRequest{Id1: int64(userId), Id2: int64(id)})
	if err != nil {
		panic(err)
	}
}

func (s *UserService) GetApiKeyForUserCredentials(username string, password string) string {
	user := s.GetByUsername(username)
	expectedPassword := user.Password
	err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
	if err != nil {
		return ""
	}
	return user.ApiKey
}

func (s *UserService) CreateJobOffer(id int, info string, position string, companyName string, qualifications string, key string) {
	s.userRepo.CreateJobOffer(id, info, position, companyName, qualifications, key)
	return
}

func (s *UserService) ActivateAccount(token string) {
	s.userRepo.ActivateAccount(token)
	return
}

func (s *UserService) GetApiKeyForUsername(username string) string {
	return s.GetByUsername(username).ApiKey
}

func (s *UserService) GetPrivateStatusForUserId(id int64) bool {
	return s.GetByID(int(id)).IsPrivate
}

func (s *UserService) SearchOffers(text string) []model.JobOffer {
	if text == "" {
		return s.userRepo.GetAllOffers()
	}
	return s.userRepo.SearchOffers(text)
}

func (s *UserService) ForgotPassword(username string) int {
	user := s.GetByUsername(username)
	if user.ID == 0 {
		return 0
	}
	newPass := GenerateRandomString(10)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPass), 8)
	user.Password = string(hashedPassword)
	user.Forgotten = 1
	s.userRepo.Save(user)

	SendActivationMail(user.Email, newPass, "")
	return 1
}

func (s *UserService) IncreaseMissedPasswordCounter(user model.User) {
	user.MissedPasswordCounter++
	if user.MissedPasswordCounter > 5 {
		user.IsBlocked = true
		user.BlockedDate = time.Now()
	}
	s.userRepo.Save(user)
}

func (s *UserService) RefreshMissedPasswordCounter(user model.User) {
	user.MissedPasswordCounter = 0
	s.userRepo.Save(user)
}

func (s *UserService) SendPinFor2Auth(username string, password string) string {
	user := s.GetByUsername(username)
	if user.ID == 0 {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth", "username": username}).Warn("Username dosn't exist.")
		return "Wrong username"
	} else if s.IsBlocked(user) {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth", "username": username}).Info("Account is blocked.")
		return "Your account is currently blocked. Try next day again."
	} else if !user.IsActivated {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth", "username": username}).Info("Account isn't activated.")
		return "Your have to activate your profile first."
	}

	expectedPassword := user.Password
	err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
	if err != nil {
		s.IncreaseMissedPasswordCounter(user)
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPinFor2Auth", "username": username}).Warn("Incorrect password.")
		return "Wrong password"
	}

	pin := GenerateRandomNumber(4)
	hashedPin, _ := bcrypt.GenerateFromPassword([]byte(pin), 8)
	user.Pin = string(hashedPin)
	user.PinCreatedDate = time.Now()
	SendActivationMail(user.Email, "", pin)
	s.userRepo.Save(user)
	return ""
}

func (s *UserService) SendPasswordlessToken(username string) string {
	user := s.GetByUsername(username)
	if user.ID == 0 {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPasswordlessToken", "username": username}).Warn("Username dosn't exist.")
		return "Wrong username"
	}

	expirationTime := time.Now().Add(3 * time.Minute)
	claims := &Claims{
		Username: username,
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Role:     "ROLE_PASSWORDLESS",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "SendPasswordlessToken", "username": username}).Warn("Token error.")
		return "token error"
	}

	SendActivationMail(user.Email, "token", tokenString)

	return ""
}

func (s *UserService) LoginPasswordless(id int) (string, bool) {
	user := s.GetByID(id)
	if s.IsBlocked(user) {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "LoginPasswordless", "user_id": id}).Info("Account is blocked.")
		return "Your account is currently blocked. Try next day again.", false
	} else if !user.IsActivated {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "LoginPasswordless", "user_id": id}).Info("Account isn't activated.")
		return "Your have to activate your profile first.", false
	}

	s.RefreshMissedPasswordCounter(user)
	expirationTime := time.Now().Add(60 * time.Minute)

	permissions := []string{"GetAllByRecieverId", "AcceptRequest", "DeclineRequest", "SendRequest", "SendMessage",
		"FindMessages", "GetNotifications", "UpdateUser", "AddExperience", "RemoveExperience", "AddInterest", "RemoveInterest",
		"BlockUser", "GetUserByUsername"}

	claims := &Claims{
		Username:    user.UserName,
		Id:          strconv.FormatUint(uint64(user.ID), 10),
		Role:        "ROLE_USER",
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.WithFields(log.Fields{"service_name": "user-service", "method_name": "LoginPasswordless", "user_id": id}).Info("Invald token.")
		return "", false
	}
	return tokenString, true
}

func (s *UserService) GetRecommendedConnections(id int64) []model.User {
	var users []model.User

	conn, err := grpc.Dial("user-service:8700", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbConn.NewConnectionServiceClient(conn)
	response, err := client.GetRecommendedConnections(context.Background(), &pbConn.GetById{Id: uint64(id)})
	if err != nil {
		panic(err)
	}

	for _, id := range response.GetIds() {
		users = append(users, s.GetByID(int(id)))
	}

	return users
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "1234567890"

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateRandomNumber(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}
