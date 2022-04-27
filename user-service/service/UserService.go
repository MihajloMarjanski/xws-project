package service

import (
	"user-service/repo"
	//"user-service/model"
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

func (s *UserService) CloseDB() error {
	return s.userRepo.Close()
}

func (s *UserService) CreateUser(name string, email string, password string) int {

	return s.userRepo.CreateUser(name, email, password)
}
