package service

import (
	"time"
	"user-service/model"
	"user-service/repo"
)

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

func (s *UserService) SearchUsers(username string) []model.UserDTO {
	return s.userRepo.SearchUsers(username)
}

func (s *UserService) GetByID(id int) model.UserDTO {
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

	return s.userRepo.CreateUser(name, email, password, username, gender, phonenumber, dateofbirth, biography)
}

func (s *UserService) AddInterest(interest string, userId uint) int {

	return s.userRepo.AddInterest(interest, userId)
}

func (s *UserService) AddExperience(company string, position string, from time.Time, until time.Time, userId uint) int {

	return s.userRepo.AddExperience(company, position, from, until, userId)
}

func (s *UserService) UpdateUser(id uint, name string, email string, password string, username string, gender model.Gender, phonenumber string, dateofbirth time.Time, biography string) int {

	return s.userRepo.UpdateUser(id, name, email, password, username, gender, phonenumber, dateofbirth, biography)
}
