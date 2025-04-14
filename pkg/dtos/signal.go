package dtos

type Signal struct {
	Symbol        string `json:"symbol"`         // BTCUSDT
	Timeframe     string `json:"timeframe"`      // 1m, 5m, 15m, 1h, 4h, 1d
	Signal        string `json:"signal"`         // buy, sell, hold
	IndicatorData string `json:"indicator_data"` // JSON string of indicator data
	LastTrade     string `json:"last_trade"`     // JSON string of last trade data
	TradeType     string `json:"trade_type"`     // 1- Ma cross, 2- RSI, 3- MACD etc.
}

type MaData struct {
	MA50   string `json:"ma50"`   // 100.0
	MA200  string `json:"ma200"`  // 200.0
	Signal string `json:"signal"` // buy, sell, hold
}

type AddSignalIntervalReq struct {
	Symbol     string `json:"symbol"`      // BTCUSDT
	Interval   string `json:"interval"`    // 1m, 5m, 15m, 1h, 4h, 1d
	ExchangeId string `json:"exchange_id"` // 1
}
type AddSignalIntervalRes struct {
	ID       string `json:"id"`       // 1
	Symbol   string `json:"symbol"`   // BTCUSDT
	Interval string `json:"interval"` // 1m, 5m, 15m, 1h, 4h, 1d
}
type UpdateSignalIntervalReq struct {
	ID         string `json:"id"`          // 1
	Symbol     string `json:"symbol"`      // BTCUSDT
	Interval   string `json:"interval"`    // 1m, 5m, 15m, 1h, 4h, 1d
	ExchangeId string `json:"exchange_id"` // 1
}

type UpdateSignalIntervalRes struct {
	ID       string `json:"id"`       // 1
	Symbol   string `json:"symbol"`   // BTCUSDT
	Interval string `json:"interval"` // 1m, 5m, 15m, 1h, 4h, 1d
}

type GetSignalIntervalReq struct {
	ID string `json:"id"` // 1
}

type GetSignalIntervalRes struct {
	ID         string `json:"id"`          // 1
	Symbol     string `json:"symbol"`      // BTCUSDT
	Interval   string `json:"interval"`    // 1m, 5m, 15m, 1h, 4h, 1d
	ExchangeId string `json:"exchange_id"` // 1
	IsActive   int    `json:"is_active"`   // 1: active, 2: inactive
}
