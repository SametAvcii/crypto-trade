package exchange

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddExchange(req dtos.AddExchangeReq) (*dtos.AddExchangeRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
func (s *service) AddExchange(req dtos.AddExchangeReq) (*dtos.AddExchangeRes, error) {
	return s.repository.AddExchange(req)
}
