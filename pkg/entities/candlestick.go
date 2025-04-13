package entities

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/shopspring/decimal"
)

type Candlestick struct {
	Base
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

func (c *Candlestick) FromDto(req *dtos.CandlestickRest) {
	c.Symbol = req.Symbol
	c.ExchangeId = req.ExchangeId
	c.Interval = req.Interval
	c.OpenTime = req.OpenTime
	c.Open = req.Open
	c.High = req.High
	c.Low = req.Low
	c.Close = req.Close
	c.Volume = req.Volume
	c.CloseTime = req.CloseTime
	c.QuoteVolume = req.QuoteVolume
	c.NumberOfTrades = req.NumberOfTrades
	c.TakerBuyBaseVolume = req.TakerBuyBaseVolume
	c.TakerBuyQuoteVolume = req.TakerBuyQuoteVolume
	c.Ignore = req.Ignore
}

func (c *Candlestick) FromDtoWs(req *dtos.CandlestickWs) {
	c.Symbol = req.Symbol
	c.ExchangeId = req.ExchangeId
	c.Interval = req.Kline.Interval
	c.OpenTime = req.Kline.StartTime
	c.Open, _ = decimal.NewFromString(req.Kline.OpenPrice)
	c.High, _ = decimal.NewFromString(req.Kline.HighPrice)
	c.Low, _ = decimal.NewFromString(req.Kline.LowPrice)
	c.Close, _ = decimal.NewFromString(req.Kline.ClosePrice)
	c.Volume, _ = decimal.NewFromString(req.Kline.BaseAssetVolume)
	c.CloseTime = req.Kline.CloseTime
	c.QuoteVolume, _ = decimal.NewFromString(req.Kline.QuoteAssetVolume)
	c.NumberOfTrades = req.Kline.NumberOfTrades
	c.TakerBuyBaseVolume, _ = decimal.NewFromString(req.Kline.TakerBuyBaseVolume)
	c.TakerBuyQuoteVolume, _ = decimal.NewFromString(req.Kline.TakerBuyQuoteVolume)
	c.Ignore, _ = decimal.NewFromString(req.Kline.Ignore)
}
