package dtos

type AddExchangeReq struct {
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}
