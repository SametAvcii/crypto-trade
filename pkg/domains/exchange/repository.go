package exchange

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddExchange(req dtos.AddExchangeReq) (*dtos.AddExchangeRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddExchange(req dtos.AddExchangeReq) (*dtos.AddExchangeRes, error) {
	var exchange entities.Exchange
	exchange.FromDto(&req)
	err := r.db.Create(&exchange).Error
	if err != nil {
		return exchange.ToDto(), err

	}
	return exchange.ToDto(), nil
}
