package entities

import "github.com/SametAvcii/crypto-trade/pkg/dtos"

/*
type AggTrade struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Symbol       string `json:"s"`
	TradeID      int64  `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	TradeTime    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
}

*/

type SymbolPrice struct {
	Base
	Symbol       string `json:"symbol"`         //BTCUSDT
	Price        string `json:"price"`          // 100.0
	Quantity     string `json:"quantity"`       // 0.1
	TradeID      int64  `json:"trade_id"`       // 123456
	TradeTime    int64  `json:"trade_time"`     // 1234567890
	IsBuyerMaker bool   `json:"is_buyer_maker"` // true
	EventTime    int64  `json:"event_time"`     // 1234567890
	EventType    string `json:"event_type"`     // "aggTrade"
	MongoID      string `json:"mongo_id"`
}

func (s *SymbolPrice) FromDto(req *dtos.AggTrade) {
	s.Symbol = req.Symbol
	s.Price = req.Price
	s.Quantity = req.Quantity
	s.TradeID = req.TradeID
	s.TradeTime = req.TradeTime
	s.IsBuyerMaker = req.IsBuyerMaker
	s.EventTime = req.EventTime
	s.EventType = req.EventType
	s.Symbol = req.Symbol
}
