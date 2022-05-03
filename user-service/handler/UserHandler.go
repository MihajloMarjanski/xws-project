package handler

import (
	"mime"
	"net/http"
	"strconv"
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
	rt, err := decodeUserBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := userHandler.userService.CreateUser(rt.Name, rt.Email, rt.Password, rt.UserName, rt.Gender, rt.PhoneNumber, rt.DateOfBirth, rt.Biography)

	switch id {
	case 0:
		http.Error(w, "Couldn't persist user in database", http.StatusBadRequest)
		return
	case -1:
		http.Error(w, "Couldn't create a user with given values (invalid email address)", http.StatusBadRequest)
		return
	case -2:
		http.Error(w, "Couldn't create a user with given values (name can't be empty)", http.StatusBadRequest)
		return
	case -3:
		http.Error(w, "Couldn't create a user with given values (username can't be empty)", http.StatusBadRequest)
		return
	case -4:
		http.Error(w, "Couldn't create a user with given values (username already exists)", http.StatusBadRequest)
		return
	case -5:
		http.Error(w, "Couldn't create a user with given values (password can't be empty)", http.StatusBadRequest)
		return
	case -6:
		http.Error(w, "Couldn't create a user with given values (invalid gender)", http.StatusBadRequest)
		return
	}

	renderJSON(w, model.ResponseId{Id: id})
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

	switch id {
	case 0:
		http.Error(w, "Couldn't persist user in database", http.StatusBadRequest)
		return
	case -1:
		http.Error(w, "Couldn't create a user with given values (invalid email address)", http.StatusBadRequest)
		return
	case -2:
		http.Error(w, "Couldn't create a user with given values (name can't be empty)", http.StatusBadRequest)
		return
	case -3:
		http.Error(w, "Couldn't create a user with given values (username can't be empty)", http.StatusBadRequest)
		return
	case -5:
		http.Error(w, "Couldn't create a user with given values (password can't be empty)", http.StatusBadRequest)
		return
	case -6:
		http.Error(w, "Couldn't create a user with given values (invalid gender)", http.StatusBadRequest)
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
