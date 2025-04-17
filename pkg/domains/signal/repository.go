package signal

import (
	"context"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error)
	GetSignalInterval(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error)
	GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error)
	DeleteSignalInterval(ctx context.Context, id string) error
	UpdateSignalInterval(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error) {
	var signal entities.SignalInterval
	signal.FromDto(&req)
	err := r.db.WithContext(ctx).Create(&signal).Error
	if err != nil {
		return signal.ToDto(), err
	}
	return signal.ToDto(), nil
}

func (r *repository) GetSignalInterval(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error) {
	var signal entities.SignalInterval
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&signal).Error
	if err != nil {
		return signal.GetDto(), err
	}
	return signal.GetDto(), nil
}

func (r *repository) GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error) {
	var signals []entities.SignalInterval
	err := r.db.WithContext(ctx).Find(&signals).Error
	if err != nil {
		return nil, err
	}
	var result []dtos.GetSignalIntervalRes
	for _, signal := range signals {
		result = append(result, signal.GetDto())
	}
	return result, nil
}

func (r *repository) DeleteSignalInterval(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.SignalInterval{}).Error
}

func (r *repository) UpdateSignalInterval(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error) {
	var signal entities.SignalInterval
	err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&signal).Error
	if err != nil {
		return dtos.UpdateSignalIntervalRes{}, err
	}

	signal.UpdateFromDto(req)
	err = r.db.WithContext(ctx).Where("id = ?", req.ID).Updates(&signal).Error
	if err != nil {
		return dtos.UpdateSignalIntervalRes{}, err
	}

	return signal.ToDtoUpdate(), nil
}
