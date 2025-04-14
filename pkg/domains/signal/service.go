package signal

import (
	"context"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type Service interface {
	AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error)
	UpdateSignalIntervals(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error)
	DeleteSignalIntervals(ctx context.Context, id string) error
	GetSignalIntervalById(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error)
	GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error) {
	return s.repository.AddSignalIntervals(ctx, req)
}

func (s *service) UpdateSignalIntervals(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error) {
	return s.repository.UpdateSignalInterval(ctx, req)
}

func (s *service) DeleteSignalIntervals(ctx context.Context, id string) error {
	return s.repository.DeleteSignalInterval(ctx, id)
}

func (s *service) GetSignalIntervalById(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error) {
	return s.repository.GetSignalInterval(ctx, id)
}

func (s *service) GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error) {
	return s.repository.GetAllSignalIntervals(ctx)
}
