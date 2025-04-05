package dtos

type AddSymbolReq struct {
	Symbol     string `json:"symbol"` //BTCUSDT
	ExchangeID string `json:"exchange_id"`
}

type AddSymbolRes struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"` //BTCUSDT
	ExchangeID string `json:"exchange_id"`
}
