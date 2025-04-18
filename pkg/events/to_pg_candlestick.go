package events

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/shopspring/decimal"
)

type PgCandleStickHandler struct{}

func (d *PgCandleStickHandler) HandleMessage(msg *sarama.ConsumerMessage) {

	var mongoData dtos.MongoData
	if err := json.Unmarshal([]byte(msg.Value), &mongoData); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling message for candlestick signal",
			Message: "Error unmarshalling outer message structure: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(msg.Value),
		})
		log.Printf("Error unmarshalling outer message structure: %v", err)
		return
	}

	var payload dtos.CandlestickWs
	if err := json.Unmarshal([]byte(mongoData.Value), &payload); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling candlestick data",
			Message: "Error unmarshalling value field into CandlestickWs: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    mongoData.Value,
		})
		log.Printf("Error unmarshalling value field into CandlestickWs: %v", err)
		return
	}

	if !payload.Kline.IsKlineClosed {
		log.Println("Candlestick is not closed for write pg, skipping...")
		return
	}

	// Check if the candlestick is already in the database
	if checkCandlestickExists(payload.Kline.Symbol, payload.Kline.StartTime) {
		log.Println("Candlestick already exists, skipping...")
		return
	}

	/*var symbol entities.Symbol
	err := database.PgClient().Where("symbol = ?", payload.Kline.Symbol).First(&symbol).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching symbol from Postgres",
			Message: "Error fetching symbol from Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(mongoData.Value),
		})
		return
	}
	var exchange entities.Exchange
	err = database.PgClient().Where("id = ?", symbol.ExchangeID).First(&exchange).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching exchange from Postgres",
			Message: "Error fetching exchange from Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(mongoData.Value),
		})
		return
	}*/

	//payload.ExchangeId = exchange.ID.String()

	if err := insertCandlestick(payload); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting candlestick into Postgres",
			Message: "Error inserting candlestick into Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(mongoData.Value),
		})
		return
	}
	log.Println("Candlestick inserted into Postgres successfully")
}

func checkCandlestickExists(symbol string, startTime int64) bool {
	db := database.PgClient()
	var candlestick entities.Candlestick
	err := db.Model(&entities.Candlestick{}).Where("symbol = ? AND open_time = ?", symbol, startTime).First(&candlestick).Error
	if err != nil {
		return false
	}
	return true
}

func insertCandlestick(payload dtos.CandlestickWs) error {
	db := database.PgClient()
	openPrice, _ := decimal.NewFromString(payload.Kline.OpenPrice)
	closePrice, _ := decimal.NewFromString(payload.Kline.ClosePrice)
	highPrice, _ := decimal.NewFromString(payload.Kline.HighPrice)
	lowPrice, _ := decimal.NewFromString(payload.Kline.LowPrice)
	volume, _ := decimal.NewFromString(payload.Kline.BaseAssetVolume)
	quoteVolume, _ := decimal.NewFromString(payload.Kline.QuoteAssetVolume)
	takerBuyBaseVolume, _ := decimal.NewFromString(payload.Kline.TakerBuyBaseVolume)
	takerBuyQuoteVolume, _ := decimal.NewFromString(payload.Kline.TakerBuyQuoteVolume)
	ignore, _ := decimal.NewFromString(payload.Kline.Ignore)

	candlestick := entities.Candlestick{
		Symbol:              payload.Kline.Symbol,
		ExchangeId:          payload.ExchangeId,
		OpenTime:            payload.Kline.StartTime,
		CloseTime:           payload.Kline.CloseTime,
		Interval:            payload.Kline.Interval,
		Open:                openPrice,
		Close:               closePrice,
		High:                highPrice,
		Low:                 lowPrice,
		Volume:              volume,
		QuoteVolume:         quoteVolume,
		NumberOfTrades:      payload.Kline.NumberOfTrades,
		TakerBuyBaseVolume:  takerBuyBaseVolume,
		TakerBuyQuoteVolume: takerBuyQuoteVolume,
		Ignore:              ignore,
	}

	return db.Create(&candlestick).Error
}
