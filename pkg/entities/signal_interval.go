package entities

import (
	"strings"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

type SignalInterval struct {
	Base
	Symbol     string `json:"symbol"`   //BTCUSDT
	Interval   string `json:"interval"` // 1m, 5m, 15m, 1h, 4h, 1d
	ExchangeID string `json:"exchange_id"`
}

func (s *SignalInterval) FromDto(dto *dtos.AddSignalIntervalReq) error {
	//to lowercase symbol
	symbol := strings.ToLower(dto.Symbol)
	s.Symbol = symbol
	s.Interval = dto.Interval
	s.ExchangeID = dto.ExchangeId
	return nil
}
func (s *SignalInterval) ToDto() *dtos.AddSignalIntervalRes {
	return &dtos.AddSignalIntervalRes{
		ID:       s.ID.String(),
		Symbol:   s.Symbol,
		Interval: s.Interval,
	}
}
