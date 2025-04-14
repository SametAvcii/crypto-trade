package symbol

import (
	"context"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error)
	GetSymbol(ctx context.Context, id string) (dtos.GetSymbolRes, error)
	GetAllSymbols(ctx context.Context) ([]dtos.GetSymbolRes, error)
	DeleteSymbol(ctx context.Context, id string) error
	UpdateSymbol(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error) {
	return s.repository.AddSymbol(ctx, req)
}

func (s *service) GetSymbol(ctx context.Context, id string) (dtos.GetSymbolRes, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *service) GetAllSymbols(ctx context.Context) ([]dtos.GetSymbolRes, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) DeleteSymbol(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) UpdateSymbol(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error) {
	return s.repository.Update(ctx, req)
}
