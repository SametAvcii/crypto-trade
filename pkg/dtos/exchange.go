package dtos

type AddExchangeReq struct {
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}

type AddExchangeRes struct {
	ID    string `json:"id"`     // 1
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}

type UpdateExchangeReq struct {
	ID    string `json:"id"`     //
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}

type UpdateExchangeRes struct {
	ID    string `json:"id"`     //
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}

type GetExchangeRes struct {
	ID    string `json:"id"`     //
	Name  string `json:"name"`   // Binance
	WsUrl string `json:"ws_url"` // wss://ws-api.binance.com:443/ws-api/v3
}
