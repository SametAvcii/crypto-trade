package signal

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddSignalIntervals(req dtos.AddSignalIntervalReq) (*dtos.AddSignalIntervalRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddSignalIntervals(req dtos.AddSignalIntervalReq) (*dtos.AddSignalIntervalRes, error) {
	return s.repository.AddSignalIntervals(req)
}
