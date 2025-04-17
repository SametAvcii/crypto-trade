package symbol

import (
	"context"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error)
	GetByID(ctx context.Context, id string) (dtos.GetSymbolRes, error)
	GetAll(ctx context.Context) ([]dtos.GetSymbolRes, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error) {
	var symbol entities.Symbol
	symbol.FromDto(&req)
	err := r.db.WithContext(ctx).Create(&symbol).Error
	if err != nil {
		return symbol.ToDto(), err
	}
	return symbol.ToDto(), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (dtos.GetSymbolRes, error) {
	var symbol entities.Symbol
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&symbol).Error
	if err != nil {
		return dtos.GetSymbolRes{}, err
	}
	return symbol.ToGetDto(), nil
}

func (r *repository) GetAll(ctx context.Context) ([]dtos.GetSymbolRes, error) {
	var symbols []entities.Symbol
	err := r.db.WithContext(ctx).Find(&symbols).Error
	if err != nil {
		return nil, err
	}
	var response []dtos.GetSymbolRes
	for _, symbol := range symbols {
		response = append(response, symbol.ToGetDto())
	}
	return response, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Symbol{}).Error
}

func (r *repository) Update(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error) {
	var symbol entities.Symbol
	err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&symbol).Error
	if err != nil {
		return dtos.UpdateSymbolRes{}, err
	}
	symbol.UpdateFromDto(req)
	err = r.db.WithContext(ctx).Where("id = ?", req.ID).Updates(&symbol).Error
	if err != nil {
		return dtos.UpdateSymbolRes{}, err
	}
	return symbol.ToDtoUpdate(), nil
}
