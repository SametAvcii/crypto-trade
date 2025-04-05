package symbol

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddSymbol(req dtos.AddSymbolReq) (*dtos.AddSymbolRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
func (s *service) AddSymbol(req dtos.AddSymbolReq) (*dtos.AddSymbolRes, error) {
	return s.repository.AddSymbol(req)
}
