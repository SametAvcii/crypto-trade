package entities

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/google/uuid"
)

const (
	SymbolActive  = 1
	SymbolPassive = 2
)

type Symbol struct {
	Base
	Symbol     string    `json:"symbol"` //BTCUSDT
	ExchangeID uuid.UUID `json:"exchange_id"`
	IsActive   uint      `json:"is_active"` //1 active, 2 passive
}

func (s *Symbol) FromDto(dto *dtos.AddSymbolReq) error {
	s.Symbol = dto.Symbol
	ExchangeID, err := uuid.Parse(dto.ExchangeID)
	if err != nil {
		return err
	}
	s.ExchangeID = ExchangeID
	return nil
}

func (s *Symbol) ToDto() *dtos.AddSymbolRes {
	return &dtos.AddSymbolRes{
		ID:         s.ID.String(),
		Symbol:     s.Symbol,
		ExchangeID: s.ExchangeID.String(),
	}
}
