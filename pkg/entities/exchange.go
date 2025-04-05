package entities

const (
	// Exchange
	ExchangeActive  = 1
	ExchangePassive = 2
)

type Exchange struct {
	Base
	Name     string `json:"name"`      //Binance, Kucoin, etc
	WsUrl    string `json:"ws_url"`    //wss://ws-api.binance.com:443/ws-api/v3
	IsActive uint   `json:"is_active"` // 1 active, 2 passive
}
