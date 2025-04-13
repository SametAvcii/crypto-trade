package signal

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddSignalIntervals(req dtos.AddSignalIntervalReq) (*dtos.AddSignalIntervalRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSignalIntervals(req dtos.AddSignalIntervalReq) (*dtos.AddSignalIntervalRes, error) {
	var signal entities.SignalInterval
	signal.FromDto(&req)
	err := r.db.Create(&signal).Error
	if err != nil {
		return signal.ToDto(), err
	}
	return signal.ToDto(), nil
}
