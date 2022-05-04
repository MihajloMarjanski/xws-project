package handler

import (
	"mime"
	"net/http"
	"strconv"
	"user-service/model"
	"user-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"encoding/json"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *service.UserService
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct{
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("tajni_kljuc_za_jwt_hash")

func (userHandler *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeUserBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rt.Password), 8)

	id := userHandler.userService.CreateUser(rt.Name, rt.Email, string(hashedPassword), rt.UserName, rt.Gender, rt.PhoneNumber, rt.DateOfBirth, rt.Biography)
	if id == 0 {
		http.Error(w, "Couldn't create a user with given values", http.StatusBadRequest)
		return
	}
	renderJSON(w, model.ResponseId{Id: id})
}

func (userHandler *UserHandler) LoginUser(w http.ResponseWriter, req * http.Request) {
	var creds Credentials

	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword := userHandler.userService.GetByUsername(creds.Username).Password

	if err = bcrypt.CompareHashAndPassword([]byte(expectedPassword),[]byte(creds.Password)); err != nil{
		w.WriteHeader(http.StatusUnauthorized)
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func (userHandler *UserHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeUserBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := userHandler.userService.UpdateUser(rt.ID, rt.Name, rt.Email, rt.Password, rt.UserName, rt.Gender, rt.PhoneNumber, rt.DateOfBirth, rt.Biography)
	if id == 0 {
		http.Error(w, "Couldn't update user with given values", http.StatusBadRequest)
		return
	}
	renderJSON(w, model.ResponseId{Id: id})
}

func (userHandler *UserHandler) SearchUsers(w http.ResponseWriter, req *http.Request) {

	username, _ := mux.Vars(req)["username"]
	users := userHandler.userService.SearchUsers(username)
	renderJSON(w, users)
}

func (userHandler *UserHandler) GetUser(w http.ResponseWriter, req *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	user := userHandler.userService.GetByID(id)
	renderJSON(w, user)
}

func (userHandler *UserHandler) GetMe(w http.ResponseWriter, req *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	user := userHandler.userService.GetMe(id)
	renderJSON(w, user)
}

func (userHandler *UserHandler) RemoveExperience(w http.ResponseWriter, req *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	removed := userHandler.userService.RemoveExperience(id)
	renderJSON(w, removed)
}

func (userHandler *UserHandler) RemoveInterest(w http.ResponseWriter, req *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	removed := userHandler.userService.RemoveInterest(id)
	renderJSON(w, removed)
}

func (userHandler *UserHandler) AddInterest(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeInterestBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := userHandler.userService.AddInterest(rt.Interest, rt.UserID)
	if id == 0 {
		http.Error(w, "Couldn't add interest with given values", http.StatusBadRequest)
		return
	}
	renderJSON(w, model.ResponseId{Id: id})
}

func (userHandler *UserHandler) AddExperience(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeExperienceBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := userHandler.userService.AddExperience(rt.Company, rt.Position, rt.From, rt.Until, rt.UserID)
	if id == 0 {
		http.Error(w, "Couldn't add interest with given values", http.StatusBadRequest)
		return
	}
	renderJSON(w, model.ResponseId{Id: id})
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
