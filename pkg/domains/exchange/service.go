package exchange

import (
	"context"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error)
	Update(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error)
	GetById(ctx context.Context, id string) (dtos.GetExchangeRes, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]dtos.GetExchangeRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error) {
	return s.repository.AddExchange(ctx, req)
}

func (s *service) Update(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error) {
	return s.repository.UpdateExchange(ctx, req)
}

func (s *service) GetById(ctx context.Context, id string) (dtos.GetExchangeRes, error) {
	return s.repository.GetExchangeById(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.DeleteExchange(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]dtos.GetExchangeRes, error) {
	return s.repository.GetAllExchanges(ctx)
}
