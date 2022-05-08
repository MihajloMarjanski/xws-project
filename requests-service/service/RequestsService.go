package service

import (
	"requests-service/model"
	"requests-service/repo"
)

type RequestService struct {
	reqRepo *repo.RequestRepository
}

func New() (*RequestService, error) {
	reqRepo, err := repo.New()
	if err != nil {
		return nil, err
	}

	return &RequestService{
		reqRepo: reqRepo,
	}, nil
}

func (s *RequestService) GetAllByRecieverId(rid uint) []model.Request {
	return s.reqRepo.GetAllByRecieverId(rid)
}

func (s *RequestService) AcceptRequest(sid, rid uint) {
	s.reqRepo.AcceptRequest(sid, rid)
}

func (s *RequestService) DeclineRequest(sid, rid uint) {
	s.reqRepo.DeclineRequest(sid, rid)
}
