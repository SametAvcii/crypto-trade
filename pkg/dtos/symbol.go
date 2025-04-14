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

type UpdateSymbolReq struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	ExchangeID string `json:"exchange_id"`
}

type UpdateSymbolRes struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	ExchangeID string `json:"exchange_id"`
}

type GetSymbolReq struct {
	ID string `json:"id"`
}

type GetSymbolRes struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	ExchangeID string `json:"exchange_id"`
	IsActive   uint   `json:"is_active"` //1 active, 2 passive
}
