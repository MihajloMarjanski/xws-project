package handler

import (
	"mime"
	"net/http"
	"user-service/model"
	"user-service/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

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
	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := userHandler.userService.CreateUser(rt.Name, rt.Email, rt.Password, rt.UserName, rt.Gender, rt.PhoneNumber, rt.DateOfBirth, rt.Biography)
	if id == 0 {
		http.Error(w, "Couldn't create a user with given values", http.StatusBadRequest)
		return
	}
	renderJSON(w, model.ResponseId{Id: id})
}

func (userHandler *UserHandler) SearchUsers(w http.ResponseWriter, req *http.Request) {

	username, _ := mux.Vars(req)["username"]
	users := userHandler.userService.SearchUsers(username)
	renderJSON(w, users)
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
