package entities

import "github.com/SametAvcii/crypto-trade/pkg/dtos"

type Signal struct {
	Base
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"` // 1m, 5m, 15m, 1h, 4h, 1d
	Signal    string `json:"signal"`    // buy, sell, hold
	Indicator string `json:"indicator" gorm:"type:jsonb"`
	LastTrade string `json:"last_trade" gorm:"type:jsonb"` // JSON string of last trade data
}

func (s *Signal) FromDto(dto dtos.Signal) {
	s.Symbol = dto.Symbol
	s.Timeframe = dto.Timeframe
	s.Signal = dto.Signal
	s.Indicator = dto.IndicatorData
	s.LastTrade = dto.LastTrade
}
