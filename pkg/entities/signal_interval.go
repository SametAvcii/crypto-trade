package entities

import (
	"strings"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/google/uuid"
)

type SignalInterval struct {
	Base
	Symbol     string    `json:"symbol"`   //BTCUSDT
	Interval   string    `json:"interval"` // 1m, 5m, 15m, 1h, 4h, 1d
	ExchangeID uuid.UUID `json:"exchange_id"`
	IsActive   int       `json:"is_active"` // 1: active, 2: inactive
}

func (s *SignalInterval) FromDto(dto *dtos.AddSignalIntervalReq) error {
	//to lowercase symbol
	symbol := strings.ToLower(dto.Symbol)
	s.Symbol = symbol
	s.Interval = dto.Interval
	s.ExchangeID = uuid.MustParse(dto.ExchangeId)
	return nil
}
func (s *SignalInterval) ToDto() dtos.AddSignalIntervalRes {
	return dtos.AddSignalIntervalRes{
		ID:       s.ID.String(),
		Symbol:   s.Symbol,
		Interval: s.Interval,
	}
}
func (s *SignalInterval) UpdateFromDto(dto dtos.UpdateSignalIntervalReq) error {
	if dto.Symbol != "" {
		s.Symbol = strings.ToLower(dto.Symbol)
	}
	if dto.Interval != "" {
		s.Interval = dto.Interval
	}
	if dto.ExchangeId != "" {
		s.ExchangeID = uuid.MustParse(dto.ExchangeId)
	}
	return nil
}

func (s *SignalInterval) GetDto() dtos.GetSignalIntervalRes {
	return dtos.GetSignalIntervalRes{
		ID:         s.ID.String(),
		Symbol:     s.Symbol,
		Interval:   s.Interval,
		ExchangeId: s.ExchangeID.String(),
		IsActive:   s.IsActive,
	}
}

func (s *SignalInterval) ToDtoUpdate() dtos.UpdateSignalIntervalRes {
	return dtos.UpdateSignalIntervalRes{
		ID:       s.ID.String(),
		Symbol:   s.Symbol,
		Interval: s.Interval,
	}
}
