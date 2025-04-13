package entities

import "github.com/SametAvcii/crypto-trade/pkg/dtos"

const (
	// Exchange
	ExchangeActive  = 1
	ExchangePassive = 2
)

type Exchange struct {
	Base
	Name     string `json:"name"`      //Binance, Kucoin, etc
	WsUrl    string `json:"ws_url"`    //wss://ws-api.binance.com:443/ws-api/v3
	RestUrl  string `json:"rest_url"`  //https://api.binance.com/api/v3
	IsActive uint   `json:"is_active"` // 1 active, 2 passive
}

func (e *Exchange) FromDto(req *dtos.AddExchangeReq) {
	e.Name = req.Name
	e.WsUrl = req.WsUrl
	e.IsActive = ExchangeActive
}

func (e *Exchange) ToDto() *dtos.AddExchangeRes {
	return &dtos.AddExchangeRes{
		ID:    e.ID.String(),
		Name:  e.Name,
		WsUrl: e.WsUrl,
	}
}
