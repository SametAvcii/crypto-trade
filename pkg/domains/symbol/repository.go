package symbol

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	AddSymbol(req dtos.AddSymbolReq) (*dtos.AddSymbolRes, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSymbol(req dtos.AddSymbolReq) (*dtos.AddSymbolRes, error) {
	var symbol entities.Symbol
	symbol.FromDto(&req)
	err := r.db.Create(&symbol).Error
	if err != nil {
		return symbol.ToDto(), err

	}
	return symbol.ToDto(), nil
}
