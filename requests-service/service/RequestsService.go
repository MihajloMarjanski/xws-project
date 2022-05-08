package service

import (
	"requests-service/model"
	"requests-service/repo"
)

type RequestsService struct {
	reqRepo *repo.RequestsRepository
}

func New() (*RequestsService, error) {
	reqRepo, err := repo.New()
	if err != nil {
		return nil, err
	}

	return &RequestsService{
		reqRepo: reqRepo,
	}, nil
}

func (s *RequestsService) CloseDB() error {
	return s.reqRepo.Close()
}

func (s *RequestsService) GetAllByRecieverId(rid uint) []model.Request {
	return s.reqRepo.GetAllByRecieverId(rid)
}

func (s *RequestsService) AcceptRequest(sid, rid uint) {
	s.reqRepo.AcceptRequest(sid, rid)
}

func (s *RequestsService) DeclineRequest(sid, rid uint) {
	s.reqRepo.DeclineRequest(sid, rid)
}

func (s *RequestsService) SendRequest(sid, rid uint) {
	s.reqRepo.SendRequest(sid, rid)
}
