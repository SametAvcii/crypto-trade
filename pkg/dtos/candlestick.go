package dtos

import "github.com/shopspring/decimal"

/*
1744171200000,

	"75073.34000000",
	"77979.30000000",
	"75000.00000000",
	"77694.12000000",
	"9140.27184000",
	1744185599999,
	"702615734.33333770",
	1111545,
	"4869.28598000",
	"374294299.83193930",
	"0"
*/
type CandlestickRest struct {
	Symbol              string          `json:"symbol"`                       // Symbol
	ExchangeId          string          `json:"exchange_id"`                  // Exchange
	Interval            string          `json:"interval"`                     // Interval
	OpenTime            int64           `json:"open_time"`                    // Open time
	Open                decimal.Decimal `json:"open"`                         // Open
	High                decimal.Decimal `json:"high"`                         // High
	Low                 decimal.Decimal `json:"low"`                          // Low
	Close               decimal.Decimal `json:"close"`                        // Close
	Volume              decimal.Decimal `json:"volume"`                       // Volume
	CloseTime           int64           `json:"close_time"`                   // Close time
	QuoteVolume         decimal.Decimal `json:"quote_asset_volume"`           // Quote asset volume
	NumberOfTrades      int64           `json:"number_of_trades"`             // Number of trades
	TakerBuyBaseVolume  decimal.Decimal `json:"taker_buy_base_asset_volume"`  // Taker buy base asset volume
	TakerBuyQuoteVolume decimal.Decimal `json:"taker_buy_quote_asset_volume"` // Taker buy quote asset volume
	Ignore              decimal.Decimal `json:"ignore"`                       // Ignore

}

/*
	{
	  "e": "kline",     // Event type
	  "E": 1638747660000,   // Event time
	  "s": "BTCUSDT",    // Symbol
	  "k": {
	    "t": 1638747660000, // Kline start time
	    "T": 1638747719999, // Kline close time
	    "s": "BTCUSDT",  // Symbol
	    "i": "1m",      // Interval
	    "f": 100,       // First trade ID
	    "L": 200,       // Last trade ID
	    "o": "0.0010",  // Open price
	    "c": "0.0020",  // Close price
	    "h": "0.0025",  // High price
	    "l": "0.0015",  // Low price
	    "v": "1000",    // Base asset volume
	    "n": 100,       // Number of trades
	    "x": false,     // Is this kline closed?
	    "q": "1.0000",  // Quote asset volume
	    "V": "500",     // Taker buy base asset volume
	    "Q": "0.500",   // Taker buy quote asset volume
	    "B": "123456"   // Ignore
	  }
	}
*/
type CandlestickWs struct {
	ExchangeId string `json:"exchange_id"` // Exchange
	EventType  string `json:"e"`           // Event type
	EventTime  int64  `json:"E"`           // Event time
	Symbol     string `json:"s"`           // Symbol
	Kline      Kline  `json:"k"`           // Kline data
}
type Kline struct {
	StartTime           int64  `json:"t"` // Kline start time
	CloseTime           int64  `json:"T"` // Kline close time
	Symbol              string `json:"s"` // Symbol
	Interval            string `json:"i"` // Interval
	FirstTradeID        int64  `json:"f"` // First trade ID
	LastTradeID         int64  `json:"L"` // Last trade ID
	OpenPrice           string `json:"o"` // Open price
	ClosePrice          string `json:"c"` // Close price
	HighPrice           string `json:"h"` // High price
	LowPrice            string `json:"l"` // Low price
	BaseAssetVolume     string `json:"v"` // Base asset volume
	NumberOfTrades      int64  `json:"n"` // Number of trades
	IsKlineClosed       bool   `json:"x"` // Is this kline closed?
	QuoteAssetVolume    string `json:"q"` // Quote asset volume
	TakerBuyBaseVolume  string `json:"V"` // Taker buy base asset volume
	TakerBuyQuoteVolume string `json:"Q"` // Taker buy quote asset volume
	Ignore              string `json:"B"` // Ignore
}
