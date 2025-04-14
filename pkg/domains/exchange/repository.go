package exchange

import (
	"context"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error)
	UpdateExchange(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error)
	DeleteExchange(ctx context.Context, id string) error
	GetExchangeById(ctx context.Context, id string) (dtos.GetExchangeRes, error)
	GetAllExchanges(ctx context.Context) ([]dtos.GetExchangeRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error) {
	var exchange entities.Exchange
	exchange.FromDto(req)
	err := r.db.WithContext(ctx).Create(&exchange).Error
	if err != nil {
		return exchange.ToDto(), err
	}
	return exchange.ToDto(), nil
}

func (r *repository) UpdateExchange(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error) {
	var exchange entities.Exchange
	if err := r.db.WithContext(ctx).First(&exchange, req.ID).Error; err != nil {
		return dtos.UpdateExchangeRes{}, err
	}

	exchange.FromDtoUpdate(req)
	if err := r.db.WithContext(ctx).Where("id = ?", req.ID).Updates(&exchange).Error; err != nil {
		return dtos.UpdateExchangeRes{}, err
	}

	return exchange.ToDtoUpdate(), nil
}

func (r *repository) DeleteExchange(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Exchange{}, id).Error
}

func (r *repository) GetExchangeById(ctx context.Context, id string) (dtos.GetExchangeRes, error) {
	var exchange entities.Exchange
	if err := r.db.WithContext(ctx).First(&exchange, id).Error; err != nil {
		return dtos.GetExchangeRes{}, err
	}
	return exchange.ToDtoGet(), nil
}

func (r *repository) GetAllExchanges(ctx context.Context) ([]dtos.GetExchangeRes, error) {
	var exchanges []entities.Exchange
	if err := r.db.WithContext(ctx).Find(&exchanges).Error; err != nil {
		return nil, err
	}

	var result []dtos.GetExchangeRes
	for _, exchange := range exchanges {
		result = append(result, exchange.ToDtoGet())
	}
	return result, nil
}
