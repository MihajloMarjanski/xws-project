package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	pbReq "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"user-service/model"
	"user-service/repo"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("tajni_kljuc_za_jwt_hash")

type UserService struct {
	userRepo *repo.UserRepository
}

func New() (*UserService, error) {

	userRepo, err := repo.New()
	if err != nil {
		return nil, err
	}

	return &UserService{
		userRepo: userRepo,
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
	_, err = http.Post("https://localhost:8600/email/activation", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Print(err.Error())
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
		return s.userRepo.UpdateUser(id, name, email, string(hashedPassword), username, gender, phonenumber, dateofbirth, biography, isPrivate)
	}
	return s.userRepo.UpdateUser(id, name, email, password, username, gender, phonenumber, dateofbirth, biography, isPrivate)
}

func (s *UserService) Login(username string, password string) string {
	expectedPassword := s.GetByUsername(username).Password
	id := strconv.FormatUint(uint64(s.GetByUsername(username).ID), 10)
	//hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
	if err != nil {
		return ""
	}
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &Claims{
		Username: username,
		Id:       id,
		Role:     "ROLE_USER",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}
	return tokenString
}

func (s *UserService) BlockUser(userId int, blockedUserId int) {
	user := s.GetByID(userId)
	block := s.GetByID(blockedUserId)
	if userId == blockedUserId || user.ID == 0 || block.ID == 0 {
		return
	}
	s.userRepo.BlockUser(userId, blockedUserId)
	DeleteConnection(userId, blockedUserId)
	return
}

func DeleteConnection(userId, id int) {
	conn, err := grpc.Dial("request-service:8200", grpc.WithInsecure())
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
